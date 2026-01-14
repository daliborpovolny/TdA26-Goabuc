package feeds

import (
	"context"
	"fmt"
	"sync"
	"time"

	db "tourbackend/internal/database/gen"

	"github.com/google/uuid"
)

type Service struct {
	q          *db.Queries
	staticPath string

	// SSE: Mutex for thread-safe access to clients map
	// Map key is courseID, value is a list of channels for connected clients
	clientsMux sync.RWMutex
	clients    map[string][]chan FeedPostResponse
}

func NewService(queries *db.Queries, staticPath string) *Service {
	return &Service{
		q:          queries,
		staticPath: staticPath,
		clients:    make(map[string][]chan FeedPostResponse),
	}
}

// --- SSE Helpers ---

// subscribe adds a client channel to a specific course
func (s *Service) subscribe(courseID string) chan FeedPostResponse {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	ch := make(chan FeedPostResponse, 10) // Buffer to prevent blocking
	s.clients[courseID] = append(s.clients[courseID], ch)
	return ch
}

// unsubscribe removes a client channel
func (s *Service) unsubscribe(courseID string, ch chan FeedPostResponse) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	channels := s.clients[courseID]
	for i, c := range channels {
		if c == ch {
			// Remove from slice
			s.clients[courseID] = append(channels[:i], channels[i+1:]...)
			close(c)
			break
		}
	}
}

// broadcast sends a new post to all clients listening to that course
func (s *Service) broadcast(courseID string, post FeedPostResponse) {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	fmt.Println("broadcasting!")

	for _, ch := range s.clients[courseID] {

		select {
		case ch <- post:
			fmt.Println("chose this one?")
		default:
			// Skip if channel is full to prevent blocking the whole server
		}
	}
}

// --- DB Logic ---

func (s *Service) GetFeed(ctx context.Context, courseID string) ([]FeedPostResponse, error) {
	posts, err := s.q.GetPostsByCourse(ctx, courseID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	response := make([]FeedPostResponse, 0, 5)
	for _, p := range posts {
		response = append(response, dbFeedPostToFeedPost(p))
	}
	return response, nil
}

func (s *Service) CreateManualPost(ctx context.Context, courseID, message string) (FeedPostResponse, error) {
	now := time.Now().Unix()
	newPost, err := s.q.CreatePost(ctx, db.CreatePostParams{
		Uuid:       uuid.New().String(),
		CourseUuid: courseID,
		Type:       "manual",
		Message:    message,
		IsEdited:   false, // Assuming generated bool
		CreatedAt:  now,   // Adjust based on your SQLC time configuration
		UpdatedAt:  now,
	})
	if err != nil {
		fmt.Println(err)
		return FeedPostResponse{}, err
	}

	resp := dbFeedPostToFeedPost(newPost)

	// TRIGGER SSE: Notify all listeners
	go s.broadcast(courseID, resp)

	return resp, nil
}

func (s *Service) UpdatePost(ctx context.Context, courseID, postID, message string) (FeedPostResponse, error) {
	updatedPost, err := s.q.UpdatePost(ctx, db.UpdatePostParams{
		Message:   message,
		Uuid:      postID,
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		fmt.Println(err)
		return FeedPostResponse{}, err
	}

	resp := dbFeedPostToFeedPost(updatedPost)

	// TRIGGER SSE: Notify that a post changed (Client handles logic to update UI)
	go s.broadcast(courseID, resp)

	return resp, nil
}

func (s *Service) DeletePost(ctx context.Context, courseID, postID string) error {
	return s.q.DeletePost(ctx, postID)

	// Note: You might want to broadcast a "delete" event type here if your frontend supports it,
	// but the spec only showed Post objects in SSE.
}

// StreamFeed handles the connection lifecycle for SSE
func (s *Service) StreamFeed(ctx context.Context, courseID string) chan FeedPostResponse {
	return s.subscribe(courseID)
}

func (s *Service) EndStream(courseID string, ch chan FeedPostResponse) {
	s.unsubscribe(courseID, ch)
}
