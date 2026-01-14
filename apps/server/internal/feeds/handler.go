package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"
)

type Handler struct {
	*handlers.Handler // Your existing base handler
	service           *Service
	pathToStatic      string
}

func NewHandler(pathToStatic string, service *Service, queries *db.Queries, isDeployed bool) *Handler {
	return &Handler{
		handlers.NewHandler(queries, isDeployed),
		service,
		pathToStatic,
	}
}

// GET /courses/{courseId}/feed
func (h *Handler) GetCourseFeed(c echo.Context) error {
	courseID := c.Param("courseId")

	feed, err := h.service.GetFeed(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, feed)
}

// POST /courses/{courseId}/feed
func (h *Handler) CreateFeedPost(c echo.Context) error {
	courseID := c.Param("courseId")
	var req CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	fmt.Println("creating post: ", req.Message)

	post, err := h.service.CreateManualPost(c.Request().Context(), courseID, req.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, post)
}

// PUT /courses/{courseId}/feed/{postId}
func (h *Handler) UpdateFeedPost(c echo.Context) error {
	courseID := c.Param("courseId")
	postID := c.Param("postId")

	var req UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	post, err := h.service.UpdatePost(c.Request().Context(), courseID, postID, req.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, post)
}

// DELETE /courses/{courseId}/feed/{postId}
func (h *Handler) DeleteFeedPost(c echo.Context) error {
	courseID := c.Param("courseId")
	postID := c.Param("postId")

	err := h.service.DeletePost(c.Request().Context(), courseID, postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET /courses/{courseId}/feed/stream
func (h *Handler) StreamFeed(c echo.Context) error {

	fmt.Println("in streamFeed")

	courseID := c.Param("courseId")

	// 1. Set SSE Headers
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	// 2. Subscribe to the service
	msgChan := h.service.StreamFeed(c.Request().Context(), courseID)

	// Ensure cleanup when client disconnects
	defer h.service.EndStream(courseID, msgChan)

	// 3. Listen for events or context cancellation
	// Flush keeps the connection open
	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "Streaming not supported")
	}

	// Send an initial ping or comment to establish connection (optional)
	fmt.Fprintf(c.Response(), ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case post := <-msgChan:
			fmt.Println("sending msg")
			// Format:
			// event: new_post
			// data: {...json...}
			//

			data, _ := json.Marshal(post)
			fmt.Fprintf(c.Response(), "event: new_post\n")
			fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			flusher.Flush()

		case <-c.Request().Context().Done():
			// Client disconnected
			return nil
		}
	}
}
