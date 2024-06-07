package comment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/controllers/notification"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type CommentParams struct {
	Content string          `json:"content"`
	Media   json.RawMessage `json:"media"`
	PostID  uuid.UUID       `json:"post_id"`
}

type CommentHandler struct {
	DB *database.Queries
}

func (apiCfg CommentHandler) PostComment(w http.ResponseWriter, r *http.Request, username string) {
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
		utils.ErrResponse(401, w, fmt.Sprintf("error : %v", err))
		return
	}

	commentCount, err := apiCfg.DB.UpdateCommentIncrease(r.Context(), params.PostID)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error : %v", err))
		return
	}

	post, err := apiCfg.DB.GetPost(r.Context(), params.PostID)

	if err == nil {
		notification.NotificationHandler.CreateNotification(notification.NotificationHandler{}, database.InsertNotificationParams{
			ID:                  uuid.New(),
			SenderUsername:      username,
			ReceiverUsername:    post.Username,
			Content:             fmt.Sprintf("@%v commented on your post: %s", username, utils.TrimToMaxChars(post.Content, 100)),
			NotificationsSource: "comments",
			Reference:           params.PostID.String(),
			CreatedAt:           time.Now(),
		})
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "comment deleted successfully",
		Payload: commentCount,
	})
}

func (apiCfg CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request, username string) {
	IDString := chi.URLParam(r, "commentID")

	commentID, err := uuid.Parse(IDString)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	postID, commentErr := apiCfg.DB.DeleteComment(r.Context(), commentID)

	if commentErr != nil {
		utils.ErrResponse(404, w, fmt.Sprintf("error %v", commentErr))
		return
	}

	commentCount, err := apiCfg.DB.UpdateCommentDecrease(r.Context(), postID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	notification.NotificationHandler.RemoveNotification(notification.NotificationHandler{}, database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           commentID.String(),
		NotificationsSource: "comments",
	})

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "comment deleted successfully",
		Payload: commentCount,
	})

}

func (apiCfg CommentHandler) GetPostComments(w http.ResponseWriter, r *http.Request) {
	stringID := chi.URLParam(r, "postID")

	postId, err := uuid.Parse(stringID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	comments, err := apiCfg.DB.GetComments(r.Context(), postId)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleCommentsToComments(comments))

}

func (apiCfg CommentHandler) LikeComment(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	commentID, err := uuid.Parse(stringID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.LikeComment(r.Context(), database.LikeCommentParams{
		Username:  username,
		CommentID: commentID,
	})

	if likesErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.IncreaseCommentLikeCount(r.Context(), commentID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	notification.NotificationHandler.CreateNotification(notification.NotificationHandler{}, database.InsertNotificationParams{
		ID:                  uuid.New(),
		SenderUsername:      username,
		ReceiverUsername:    response.Username,
		Content:             fmt.Sprintf("@%v likes your comment %s", username, utils.TrimToMaxChars(response.CommentText, 100)),
		NotificationsSource: "comment-likes",
		Reference:           commentID.String(),
		CreatedAt:           time.Now(),
	})

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "likes added successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg CommentHandler) UnLikeComment(w http.ResponseWriter, r *http.Request, username string) {
	stringID := chi.URLParam(r, "postID")

	commentID, err := uuid.Parse(stringID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	likesErr := apiCfg.DB.RemoveCommentLike(r.Context(), database.RemoveCommentLikeParams{
		Username:  username,
		CommentID: commentID,
	})

	if likesErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", likesErr))
		return
	}

	response, err := apiCfg.DB.DecreaseCommentLikeCount(r.Context(), commentID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v:", err))
		return
	}

	notification.NotificationHandler.RemoveNotification(notification.NotificationHandler{}, database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           commentID.String(),
		NotificationsSource: "comment-likes",
	})

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "likes removed successfully",
		Payload: response.LikesCount,
	})
}

func (apiCfg CommentHandler) EditComment(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := CommentParams{}

	decorder.Decode(&params)

	stringId := chi.URLParam(r, "commentID")

	commentId, err := uuid.Parse(stringId)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	editErr := apiCfg.DB.EditComment(r.Context(), database.EditCommentParams{
		CommentText: params.Content,
		Media:       params.Media,
		ID:          commentId,
	})

	if editErr != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error : %v", editErr))
		return
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "comment deleted successfully",
	})
}
