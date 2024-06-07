package likes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type LikeHandler struct {
	DB *database.Queries
}

func (apiCfg LikeHandler) LikePost(w http.ResponseWriter, r *http.Request, username string) {
	postId, err := utils.GetPostId(r)

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("error %v :", err))
		return
	}

	errRes := apiCfg.DB.LikePost(r.Context(), database.LikePostParams{
		Username: username,
		PostID:   postId,
	})
	fmt.Println("done")
	if errRes != nil {
		utils.ErrResponse(501, w, fmt.Sprintf("error %v:", errRes))
		return
	}

	data, _ := apiCfg.DB.UpdateLikeIncrease(r.Context(), postId)

	post, err := apiCfg.DB.GetPost(r.Context(), postId)

	if err == nil {
		apiCfg.DB.InsertNotification(r.Context(), database.InsertNotificationParams{
			ID:                  uuid.New(),
			SenderUsername:      username,
			ReceiverUsername:    post.Username,
			Content:             fmt.Sprintf("@%v likes your post: %s", username, utils.TrimToMaxChars(post.Content, 100)),
			NotificationsSource: "likes",
			Reference:           postId.String(),
			CreatedAt:           time.Now(),
		})
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "post successfully liked",
		Payload: data,
	})

}

func (apiCfg LikeHandler) UnlikePost(w http.ResponseWriter, r *http.Request, username string) {
	postId, err := utils.GetPostId(r)

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("error %v :", err))
		return
	}

	errRes := apiCfg.DB.UnlikePost(r.Context(), database.UnlikePostParams{
		Username: username,
		PostID:   postId,
	})

	if errRes != nil {
		utils.ErrResponse(501, w, fmt.Sprintf("error %v:", errRes))
		return
	}

	data, _ := apiCfg.DB.UpdateLikeDecrease(r.Context(), postId)

	apiCfg.DB.RemoveNotification(r.Context(), database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           postId.String(),
		NotificationsSource: "likes",
	})

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "post successfully unliked",
		Payload: data,
	})

}
