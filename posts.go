package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type PostParams struct {
	UserID       uuid.UUID       `json:"user_id"`
	Content      string          `json:"content"`
	Media        json.RawMessage `json:"media"`
	IsShared     bool            `json:"is_shared"`
	SharedPostId uuid.UUID       `json:"shared_post_id"`
}

func (apiCfg apiCfg) makePost(w http.ResponseWriter, r *http.Request, username string) {
	decoder := json.NewDecoder(r.Body)
	params := PostParams{}

	decoder.Decode(&params)

	postId, err := apiCfg.DB.NewPost(r.Context(), database.NewPostParams{
		ID:           uuid.New(),
		UserID:       params.UserID,
		Username:     username,
		Content:      params.Content,
		Media:        params.Media,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		IsSharedPost: params.IsShared,
		SharedPostID: uuid.NullUUID{Valid: params.IsShared, UUID: params.SharedPostId},
	})

	if err != nil {
		errResponse(400, w, fmt.Sprintf("Error %v ", err))
		return
	}

	post, err := apiCfg.DB.GetPost(r.Context(), postId)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("Error: %v", err))
		return
	}

	if params.IsShared {
		apiCfg.DB.IncreaseShareCount(r.Context(), params.SharedPostId)
	}

	jsonResponse(200, w, handlePostToPost(post))
}

func (apiCfg apiCfg) getUserPosts(w http.ResponseWriter, r *http.Request, username string) {
	profileUsername := chi.URLParam(r, "username")

	if profileUsername == "" {
		errResponse(200, w, "Invalid username")
		return
	}

	posts, err := apiCfg.DB.GetUserPosts(r.Context(), profileUsername)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("Error: %v", err))
		return
	}

	jsonResponse(200, w, handlePostsToPosts(posts))

}

func (apiCfg apiCfg) getPost(w http.ResponseWriter, r *http.Request, username string) {
	postId := chi.URLParam(r, "postID")

	if postId == "" {
		errResponse(404, w, "post id not found")
		return
	}

	id, err := uuid.Parse(postId)

	if err != nil {
		errResponse(403, w, "invalid post id")
		return
	}

	post, err := apiCfg.DB.GetPost(r.Context(), id)

	if err != nil {
		errResponse(400, w, fmt.Sprintf("error %v ", err))
		return
	}

	jsonResponse(200, w, handlePostToPost(post))
}

func (apiCfg apiCfg) getUserFeeds(w http.ResponseWriter, r *http.Request, username string) {
	var idString string = r.URL.Query().Get("id")

	if idString == "" {
		errResponse(401, w, "no id found")
		return
	}

	userID, err := uuid.Parse(idString)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error: %v", err))
		return
	}

	feeds, err := apiCfg.DB.GetFeedPosts(r.Context(), userID)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error: %v", err))
		return
	}

	jsonResponse(200, w, handleFeedsToFeeds(feeds))
}

func (apiCfg apiCfg) getBookmarkedPosts(w http.ResponseWriter, r *http.Request, username string) {
	posts, err := apiCfg.DB.GetBookmarkedPosts(r.Context(), username)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error: %v", err))
		return
	}

	jsonResponse(200, w, handleBookmarksToBookmarks(posts))
}

func (apiCfg apiCfg) deletePost(w http.ResponseWriter, r *http.Request, username string) {
	idString := chi.URLParam(r, "postID")

	postId, err := uuid.Parse(idString)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	deleteErr := apiCfg.DB.DeletePost(r.Context(), postId)

	if deleteErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "post delete successfully",
	})
}

func (apiCfg apiCfg) editPost(w http.ResponseWriter, r *http.Request, username string) {
	idString := chi.URLParam(r, "postID")

	postId, err := uuid.Parse(idString)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	decorder := json.NewDecoder(r.Body)

	params := PostParams{}

	decorder.Decode(&params)

	postErr := apiCfg.DB.EditPost(r.Context(), database.EditPostParams{
		Content: params.Content,
		Media:   params.Media,
		ID:      postId,
	})

	if postErr != nil {
		errResponse(404, w, fmt.Sprintf("error %v", postErr))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "post edit successful",
	})
}
