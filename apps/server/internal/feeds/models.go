package feeds

import (
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"
)

type FeedPostResponse struct {
	UUID      string `json:"uuid"`
	Type      string `json:"type"` // "manual" or "auto"
	Message   string `json:"message"`
	Edited    bool   `json:"edited"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreatePostRequest struct {
	Message string `json:"message"`
}

type UpdatePostRequest struct {
	Message string `json:"message"`
	Edited  bool   `json:"edited"` // Often ignored in logic, but present in spec
}

func dbFeedPostToFeedPost(dbFeedPost db.FeedPost) FeedPostResponse {
	return FeedPostResponse{
		UUID:      dbFeedPost.Uuid,
		Type:      dbFeedPost.Type,
		Message:   dbFeedPost.Message,
		Edited:    dbFeedPost.IsEdited,
		CreatedAt: utils.UnixToIso(dbFeedPost.CreatedAt),
		UpdatedAt: utils.UnixToIso(dbFeedPost.UpdatedAt),
	}
}
