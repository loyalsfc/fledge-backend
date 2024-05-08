package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	post, err := apiCfg.DB.NewPost(r.Context(), database.NewPostParams{
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

	jsonResponse(200, w, handlePostToPost(post))
}
