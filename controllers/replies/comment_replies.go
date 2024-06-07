package replies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type ReplyStruct struct {
	Content   string          `json:"content"`
	Media     json.RawMessage `json:"media"`
	CommentId string          `json:"comment_id"`
}

type ReplyHandler struct {
	DB *database.Queries
}

func (apiCfg ReplyHandler) ReplyComment(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := ReplyStruct{}

	decorder.Decode(&params)

	commentId, err := uuid.Parse(params.CommentId)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
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
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", replyErr))
		return
	}

	comment, err := apiCfg.DB.GetComment(r.Context(), commentId)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	apiCfg.DB.InsertNotification(r.Context(), database.InsertNotificationParams{
		ID:                  uuid.New(),
		SenderUsername:      username,
		ReceiverUsername:    comment.Username,
		Content:             fmt.Sprintf("@%v replied to your comment %s", username, utils.TrimToMaxChars(comment.CommentText, 100)),
		NotificationsSource: "comment-reply",
		Reference:           params.CommentId,
		CreatedAt:           time.Now(),
	})

	count, _ := apiCfg.DB.UpdateReplyIncrease(r.Context(), commentId)

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "reply successful",
		Payload: count,
	})
}

func (apiCfg ReplyHandler) DeleteCommetReply(w http.ResponseWriter, r *http.Request, username string) {
	idstring := chi.URLParam(r, "replyID")

	replyId, err := uuid.Parse(idstring)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	reply, err := apiCfg.DB.GetReply(r.Context(), replyId)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	replyErr := apiCfg.DB.DeleteReply(r.Context(), replyId)

	if replyErr != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", replyErr))
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

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "reply successful",
		Payload: count,
	})
}

func (apiCfg ReplyHandler) GetReplies(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "commentID")

	commentId, err := uuid.Parse(idString)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	replies, err := apiCfg.DB.GetReplies(r.Context(), commentId)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleRepliesToReplies(replies))
}

func (apiCfg ReplyHandler) LikeReply(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "replyID")

	replyID, err := uuid.Parse(stringID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.LikeReply(r.Context(), database.LikeReplyParams{
		Username: username,
		ReplyID:  replyID,
	})

	if likesErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.UpdateReplyLikesCountIncrease(r.Context(), replyID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	apiCfg.DB.InsertNotification(r.Context(), database.InsertNotificationParams{
		ID:                  uuid.New(),
		SenderUsername:      username,
		ReceiverUsername:    response.Username,
		Content:             fmt.Sprintf("@%v likes your reply %s", username, utils.TrimToMaxChars(response.ReplyText, 100)),
		NotificationsSource: "reply-likes",
		Reference:           stringID,
		CreatedAt:           time.Now(),
	})

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "likes added successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg ReplyHandler) UnLikeReply(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	replyID, err := uuid.Parse(stringID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.RemoveReplyLike(r.Context(), database.RemoveReplyLikeParams{
		Username: username,
		ReplyID:  replyID,
	})

	if likesErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.UpdateReplyLikesCountDecrease(r.Context(), replyID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	apiCfg.DB.RemoveNotification(r.Context(), database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           stringID,
		NotificationsSource: "reply-likes",
	})

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "likes removed successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg ReplyHandler) EditReply(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := ReplyStruct{}

	decorder.Decode(&params)

	stringId := chi.URLParam(r, "replyID")

	replyId, err := uuid.Parse(stringId)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	replyErr := apiCfg.DB.EditReply(r.Context(), database.EditReplyParams{
		ReplyText: params.Content,
		Media:     params.Media,
		ID:        replyId,
	})

	if replyErr != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error %v", replyErr))
		return
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "reply successful",
	})
}
