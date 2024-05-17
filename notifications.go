package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loyalsfc/fledge-backend/internal/database"
)

func (apiCfg apiCfg) createNotification(params database.InsertNotificationParams) {
	ctx := context.Background()
	apiCfg.DB.InsertNotification(ctx, params)
}

func (apiCfg apiCfg) removeNotification(params database.RemoveNotificationParams) {
	ctx := context.Background()
	apiCfg.DB.RemoveNotification(ctx, params)
}

func (apiCfg apiCfg) getUserNotifications(w http.ResponseWriter, r *http.Request, username string) {
	notifications, err := apiCfg.DB.GetUserNotifications(r.Context(), username)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v ", err))
		return
	}

	jsonResponse(200, w, handleNotificationsToNotifications(notifications))
}

type NotificationParams struct {
	Username           string `json:"username"`
	Reference          string `json:"reference"`
	NotificationSource string `json:"notifications_source"`
}

func (apiCfg apiCfg) markNotificationAsRead(w http.ResponseWriter, r *http.Request, username string) {
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

	jsonResponse(200, w, Response{
		Status:  "success",
		Message: "notification mark as read",
	})
}
