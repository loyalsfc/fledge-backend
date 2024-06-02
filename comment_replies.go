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

type ReplyStruct struct {
	Content   string          `json:"content"`
	Media     json.RawMessage `json:"media"`
	CommentId string          `json:"comment_id"`
}

func (apiCfg apiCfg) replyComment(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := ReplyStruct{}

	decorder.Decode(&params)

	commentId, err := uuid.Parse(params.CommentId)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	replyErr := apiCfg.DB.NewReply(r.Context(), database.NewReplyParams{
		ID:        uuid.New(),
		ReplyText: params.Content,
		Media:     params.Media,
		Username:  username,
		CommentID: commentId,
	})

	if replyErr != nil {
		errResponse(401, w, fmt.Sprintf("error %v", replyErr))
		return
	}

	comment, err := apiCfg.DB.GetComment(r.Context(), commentId)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	apiCfg.DB.InsertNotification(r.Context(), database.InsertNotificationParams{
		ID:                  uuid.New(),
		SenderUsername:      username,
		ReceiverUsername:    comment.Username,
		Content:             fmt.Sprintf("@%v replied to your comment %s", username, trimToMaxChars(comment.CommentText, 100)),
		NotificationsSource: "comment-reply",
		Reference:           params.CommentId,
		CreatedAt:           time.Now(),
	})

	count, _ := apiCfg.DB.UpdateReplyIncrease(r.Context(), commentId)

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "reply successful",
		Payload: count,
	})
}

func (apiCfg apiCfg) deleteCommetReply(w http.ResponseWriter, r *http.Request, username string) {
	idstring := chi.URLParam(r, "replyID")

	replyId, err := uuid.Parse(idstring)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	reply, err := apiCfg.DB.GetReply(r.Context(), replyId)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	replyErr := apiCfg.DB.DeleteReply(r.Context(), replyId)

	if replyErr != nil {
		errResponse(401, w, fmt.Sprintf("error %v", replyErr))
		return
	}

	apiCfg.DB.RemoveNotification(r.Context(), database.RemoveNotificationParams{
		SenderUsername:      username,
		NotificationsSource: "comment-reply",
		Reference:           reply.CommentID.String(),
	})

	count, err := apiCfg.DB.UpdateReplyDecrease(r.Context(), reply.CommentID)

	if err != nil {
		fmt.Println("err", err)
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "reply successful",
		Payload: count,
	})
}

func (apiCfg apiCfg) getReplies(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "commentID")

	commentId, err := uuid.Parse(idString)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	replies, err := apiCfg.DB.GetReplies(r.Context(), commentId)

	if err != nil {
		errResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	jsonResponse(200, w, handleRepliesToReplies(replies))
}

func (apiCfg apiCfg) likeReply(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "replyID")

	replyID, err := uuid.Parse(stringID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.LikeReply(r.Context(), database.LikeReplyParams{
		Username: username,
		ReplyID:  replyID,
	})

	if likesErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.UpdateReplyLikesCountIncrease(r.Context(), replyID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	apiCfg.DB.InsertNotification(r.Context(), database.InsertNotificationParams{
		ID:                  uuid.New(),
		SenderUsername:      username,
		ReceiverUsername:    response.Username,
		Content:             fmt.Sprintf("@%v likes your reply %s", username, trimToMaxChars(response.ReplyText, 100)),
		NotificationsSource: "reply-likes",
		Reference:           stringID,
		CreatedAt:           time.Now(),
	})

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "likes added successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg apiCfg) unLikeReply(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	replyID, err := uuid.Parse(stringID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.RemoveReplyLike(r.Context(), database.RemoveReplyLikeParams{
		Username: username,
		ReplyID:  replyID,
	})

	if likesErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.UpdateReplyLikesCountDecrease(r.Context(), replyID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	apiCfg.DB.RemoveNotification(r.Context(), database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           stringID,
		NotificationsSource: "reply-likes",
	})

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "likes removed successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg apiCfg) editReply(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := ReplyStruct{}

	decorder.Decode(&params)

	stringId := chi.URLParam(r, "replyID")

	replyId, err := uuid.Parse(stringId)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	replyErr := apiCfg.DB.EditReply(r.Context(), database.EditReplyParams{
		ReplyText: params.Content,
		Media:     params.Media,
		ID:        replyId,
	})

	if replyErr != nil {
		errResponse(401, w, fmt.Sprintf("error %v", replyErr))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "reply successful",
	})
}
