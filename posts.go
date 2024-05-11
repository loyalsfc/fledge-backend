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
	UserID  uuid.UUID       `json:"user_id"`
	Content string          `json:"content"`
	Media   json.RawMessage `json:"media"`
}

func (apiCfg apiCfg) makePost(w http.ResponseWriter, r *http.Request, username string) {
	decoder := json.NewDecoder(r.Body)
	params := PostParams{}

	decoder.Decode(&params)

	postId, err := apiCfg.DB.NewPost(r.Context(), database.NewPostParams{
		ID:        uuid.New(),
		UserID:    params.UserID,
		Username:  username,
		Content:   params.Content,
		Media:     params.Media,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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
