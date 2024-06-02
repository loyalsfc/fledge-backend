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

	post, err := apiCfg.DB.GetPost(r.Context(), params.PostID)

	if err == nil {
		apiCfg.createNotification(database.InsertNotificationParams{
			ID:                  uuid.New(),
			SenderUsername:      username,
			ReceiverUsername:    post.Username,
			Content:             fmt.Sprintf("@%v commented on your post: %s", username, trimToMaxChars(post.Content, 100)),
			NotificationsSource: "comments",
			Reference:           params.PostID.String(),
			CreatedAt:           time.Now(),
		})
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

	postID, commentErr := apiCfg.DB.DeleteComment(r.Context(), commentID)

	if commentErr != nil {
		errResponse(404, w, fmt.Sprintf("error %v", commentErr))
		return
	}

	commentCount, err := apiCfg.DB.UpdateCommentDecrease(r.Context(), postID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	apiCfg.removeNotification(database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           commentID.String(),
		NotificationsSource: "comments",
	})

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

func (apiCfg apiCfg) likeComment(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	commentID, err := uuid.Parse(stringID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.LikeComment(r.Context(), database.LikeCommentParams{
		Username:  username,
		CommentID: commentID,
	})

	if likesErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.IncreaseCommentLikeCount(r.Context(), commentID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	apiCfg.DB.InsertNotification(r.Context(), database.InsertNotificationParams{
		ID:                  uuid.New(),
		SenderUsername:      username,
		ReceiverUsername:    response.Username,
		Content:             fmt.Sprintf("@%v likes your comment %s", username, trimToMaxChars(response.CommentText, 100)),
		NotificationsSource: "comment-likes",
		Reference:           commentID.String(),
		CreatedAt:           time.Now(),
	})

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "likes added successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg apiCfg) unLikeComment(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	commentID, err := uuid.Parse(stringID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.RemoveCommentLike(r.Context(), database.RemoveCommentLikeParams{
		Username:  username,
		CommentID: commentID,
	})

	if likesErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.DecreaseCommentLikeCount(r.Context(), commentID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	apiCfg.DB.RemoveNotification(r.Context(), database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           commentID.String(),
		NotificationsSource: "comment-likes",
	})

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "likes removed successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg apiCfg) editComment(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := CommentParams{}

	decorder.Decode(&params)

	stringId := chi.URLParam(r, "commentID")

	commentId, err := uuid.Parse(stringId)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	editErr := apiCfg.DB.EditComment(r.Context(), database.EditCommentParams{
		CommentText: params.Content,
		Media:       params.Media,
		ID:          commentId,
	})

	if editErr != nil {
		errResponse(401, w, fmt.Sprintf("error : %v", editErr))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "comment deleted successfully",
	})
}
