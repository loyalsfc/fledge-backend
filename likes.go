package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func trimToMaxChars(s string, maxChars int) string {
	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars]
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

	post, err := apiCfg.DB.GetPost(r.Context(), postId)

	if err == nil {
		apiCfg.createNotification(database.InsertNotificationParams{
			ID:                  uuid.New(),
			SenderUsername:      username,
			ReceiverUsername:    post.Username,
			Content:             fmt.Sprintf("@%v likes your post: %s", username, trimToMaxChars(post.Content, 100)),
			NotificationsSource: "likes",
			Reference:           postId.String(),
			CreatedAt:           time.Now(),
		})
	}

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

	apiCfg.removeNotification(database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           postId.String(),
		NotificationsSource: "likes",
	})

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "post successfully unliked",
		Payload: data,
	})

}

// func (apiCfg apiCfg) postLikes(w http.ResponseWriter, r *http.Request, username string) {

// }
