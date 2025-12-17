package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"
	"tourbackend/internal/utils"

	"github.com/labstack/echo/v4"
)

var COOKIE_LIFETIME = time.Hour * 24 * 7

type AuthHandler struct {
	*handlers.Handler
}

func NewAuthHandler(queries *db.Queries, isDeployed bool) *AuthHandler {
	return &AuthHandler{handlers.NewHandler(queries, isDeployed)}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	r := h.NewReqCtx(c)

	var req LoginRequest
	if err := r.Echo.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request body")
	}

	user, err := r.Queries.GetUserByEmail(r.Ctx, req.Email)
	if err != nil {
		if utils.IsNoRowsError(err) {
			return r.Error(http.StatusBadRequest, "Unknown email")
		}
		return err
	}

	isCorrect := utils.CheckPasswordHash(req.Password, user.Hash)
	if !isCorrect {
		return r.Error(http.StatusUnauthorized, "invalid password")
	}

	newToken, err := utils.NewSessionToken()
	if err != nil {
		c.Logger().Errorf("failed to generate a session token: %v", err)
		return r.Error(http.StatusInternalServerError, "internal server error")
	}

	cookie := h.createHttpCookie(newToken)

	_, err = r.Queries.CreateSession(r.Ctx, db.CreateSessionParams{
		UserID:    user.ID,
		Token:     newToken,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: cookie.Expires.Unix(),
	})
	if err != nil {
		c.Logger().Errorf("failed to create the session in the database: %v", err)
		return r.Error(http.StatusInternalServerError, "internal server error")
	}

	c.SetCookie(cookie)

	c.Logger().Infof("logged in a user: %v", user.Email)
	return r.JSONMsg(http.StatusCreated, "logged in user")
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	r := h.NewReqCtx(c)

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return r.Error(http.StatusBadRequest, "unhashable password")
	}

	user, err := r.Queries.CreateUser(r.Ctx, db.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Hash:      hash,
		Email:     req.Email,
	})
	if err != nil {
		if utils.IsUniqueConstraintError(err) {
			return r.Error(http.StatusBadRequest, "email already in use")
		}

		c.Logger().Errorf("failed to create a user: %v", err)
		return r.Error(http.StatusInternalServerError, "internal server error")
	}

	newToken, err := utils.NewSessionToken()
	if err != nil {
		c.Logger().Errorf("failed to generate a session token: %v", err)
		return r.Error(http.StatusInternalServerError, "internal server error")
	}

	cookie := h.createHttpCookie(newToken)

	_, err = r.Queries.CreateSession(r.Ctx, db.CreateSessionParams{
		UserID:    user.ID,
		Token:     newToken,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: cookie.Expires.Unix(),
	})
	if err != nil {
		c.Logger().Errorf("failed to create the session in the database: %v", err)
		return r.Error(http.StatusInternalServerError, "internal server error")
	}

	c.SetCookie(cookie)

	c.Logger().Infof("registered a user: %v", user.Email)
	return r.JSONMsg(http.StatusCreated, "registered user")
}

type ProfileResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (h *AuthHandler) Profile(c echo.Context) error {
	r := h.NewReqCtx(c)

	if r.User == nil {
		return r.Error(http.StatusUnauthorized, "authentication required")
	}

	c.Logger().Infof("returned profile of user: %v", r.User.Email)

	return c.JSON(http.StatusOK, ProfileResponse{
		FirstName: r.User.FirstName,
		LastName:  r.User.LastName,
		Email:     r.User.Email,
	})
}

func (h *AuthHandler) createHttpCookie(tokenValue string) *http.Cookie {
	return &http.Cookie{
		Name:     "auth_token",
		Value:    tokenValue,
		Expires:  time.Now().Add(COOKIE_LIFETIME),
		HttpOnly: true,
		Secure:   h.IsDeployed,
	}
}

func validateToken(token string, queries *db.Queries, ctx context.Context) (*db.User, error) {

	authInfo, err := queries.GetUserBySessionToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if authInfo.ExpiresAt <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	return &db.User{
		ID:        authInfo.UserID,
		FirstName: authInfo.FirstName,
		LastName:  authInfo.LastName,
		Hash:      authInfo.Hash,
		Email:     authInfo.Email,
	}, nil
}
