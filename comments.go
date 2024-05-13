package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type CommentParams struct {
	Content string          `json:"content"`
	Media   json.RawMessage `json:"media"`
	PostID  uuid.UUID       `json:"post_id"`
}

func (apiCfg apiCfg) postComment(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := CommentParams{}

	decorder.Decode(&params)

	err := apiCfg.DB.NewComments(r.Context(), database.NewCommentsParams{
		ID:          uuid.New(),
		CommentText: params.Content,
		Media:       params.Media,
		Username:    username,
		PostID:      params.PostID,
	})

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error : %v", err))
		return
	}

	commentCount, err := apiCfg.DB.UpdateCommentIncrease(r.Context(), params.PostID)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error : %v", err))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "comment deleted successfully",
		Payload: commentCount,
	})
}

func (apiCfg apiCfg) deleteComment(w http.ResponseWriter, r *http.Request, username string) {
	IDString := chi.URLParam(r, "commentID")

	commentID, err := uuid.Parse(IDString)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	commentErr := apiCfg.DB.DeleteComment(r.Context(), commentID)

	if commentErr != nil {
		errResponse(404, w, fmt.Sprintf("error %v", commentErr))
		return
	}

	commentCount, err := apiCfg.DB.UpdateCommentDecrease(r.Context(), commentID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "comment deleted successfully",
		Payload: commentCount,
	})

}

func (apiCfg apiCfg) getPostComments(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	postId, err := uuid.Parse(stringID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	comments, err := apiCfg.DB.GetComments(r.Context(), postId)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	jsonResponse(200, w, handleCommentsToComments(comments))

}
