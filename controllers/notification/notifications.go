package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type NotificationHandler struct {
	DB *database.Queries
}

func (apiCfg NotificationHandler) CreateNotification(params database.InsertNotificationParams) {
	ctx := context.Background()
	apiCfg.DB.InsertNotification(ctx, params)
}

func (apiCfg NotificationHandler) RemoveNotification(params database.RemoveNotificationParams) {
	ctx := context.Background()
	apiCfg.DB.RemoveNotification(ctx, params)
}

func (apiCfg NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request, username string) {
	notifications, err := apiCfg.DB.GetUserNotifications(r.Context(), username)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v ", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleNotificationsToNotifications(notifications))
}

type NotificationParams struct {
	Username           string `json:"username"`
	Reference          string `json:"reference"`
	NotificationSource string `json:"notifications_source"`
}

func (apiCfg NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request, username string) {
	decoder := json.NewDecoder(r.Body)

	params := []NotificationParams{}

	decoder.Decode(&params)

	for _, item := range params {
		apiCfg.DB.MarkNotificationAsRead(r.Context(), database.MarkNotificationAsReadParams{
			SenderUsername:      item.Username,
			Reference:           item.Reference,
			NotificationsSource: item.NotificationSource,
		})
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "notification mark as read",
	})
}
