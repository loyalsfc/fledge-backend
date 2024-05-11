package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type LikesParam struct {
	PostId string `json:"post_id"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Payload int32  `json:"payload"`
}

func getPostId(r *http.Request) (uuid.UUID, error) {
	decorder := json.NewDecoder(r.Body)

	params := LikesParam{}

	decorder.Decode(&params)

	postId, err := uuid.Parse(params.PostId)

	if err != nil {
		return uuid.New(), errors.New("invalid post id")
	}

	return postId, nil
}

func (apiCfg apiCfg) likePost(w http.ResponseWriter, r *http.Request, username string) {
	postId, err := getPostId(r)

	if err != nil {
		errResponse(400, w, fmt.Sprintf("error %v :", err))
		return
	}

	errRes := apiCfg.DB.LikePost(r.Context(), database.LikePostParams{
		Username: username,
		PostID:   postId,
	})

	if errRes != nil {
		errResponse(501, w, fmt.Sprintf("error %v:", errRes))
		return
	}

	data, _ := apiCfg.DB.UpdateLikeIncrease(r.Context(), postId)

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "post successfully liked",
		Payload: data,
	})

}

func (apiCfg apiCfg) unlikePost(w http.ResponseWriter, r *http.Request, username string) {
	postId, err := getPostId(r)

	if err != nil {
		errResponse(400, w, fmt.Sprintf("error %v :", err))
		return
	}

	errRes := apiCfg.DB.UnlikePost(r.Context(), database.UnlikePostParams{
		Username: username,
		PostID:   postId,
	})

	if errRes != nil {
		errResponse(501, w, fmt.Sprintf("error %v:", errRes))
		return
	}

	data, _ := apiCfg.DB.UpdateLikeDecrease(r.Context(), postId)

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "post successfully unliked",
		Payload: data,
	})

}

// func (apiCfg apiCfg) postLikes(w http.ResponseWriter, r *http.Request, username string) {

// }
