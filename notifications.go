package main

import (
	"context"
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
