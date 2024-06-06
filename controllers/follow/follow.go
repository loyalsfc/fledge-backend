package follow

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/controllers/notification"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type FollowParams struct {
	Following uuid.UUID `json:"following"`
	Follower  uuid.UUID `json:"follower"`
}

type FollowHandler struct {
	DB *database.Queries
}

func (apiCfg FollowHandler) Follow(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := FollowParams{}

	decorder.Decode(&params)

	follow, err := apiCfg.DB.NewFollower(r.Context(), database.NewFollowerParams{
		FollowerID:  params.Follower,
		FollowingID: params.Following,
	})

	if err != nil {
		utils.ErrResponse(405, w, fmt.Sprintf("Error occured %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserById(r.Context(), params.Following)

	if err == nil {
		notification.NotificationHandler.CreateNotification(notification.NotificationHandler{}, database.InsertNotificationParams{
			ID:                  uuid.New(),
			SenderUsername:      username,
			ReceiverUsername:    user.Username,
			Content:             fmt.Sprintf("@%v is now following you", username),
			NotificationsSource: "follow",
			Reference:           params.Following.String(),
			CreatedAt:           time.Now(),
		})
	}

	utils.JsonResponse(200, w, follow)
}

func getIdFromParams(r *http.Request) (uuid.UUID, error) {
	id := chi.URLParam(r, "userID")

	if id == "" {
		return uuid.New(), errors.New("no id found")
	}

	userId, err := uuid.Parse(id)

	if err != nil {
		return uuid.New(), errors.New("invalid userId")
	}

	return userId, nil
}

func (apiCfg FollowHandler) GetFollower(w http.ResponseWriter, r *http.Request) {

	userId, err := getIdFromParams(r)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	list, err := apiCfg.DB.GetFollowers(r.Context(), userId)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleFollowersToFollowers(list))
}

func (apiCfg FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {

	userId, err := getIdFromParams(r)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	list, err := apiCfg.DB.GetFollowing(r.Context(), userId)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleFollowsToFollows(list))
}

func (apiCfg FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := FollowParams{}

	decorder.Decode(&params)

	err := apiCfg.DB.UnfollowOne(r.Context(), database.UnfollowOneParams{
		FollowerID:  params.Follower,
		FollowingID: params.Following,
	})

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("Error %v", err))
		return
	}

	notification.NotificationHandler.RemoveNotification(notification.NotificationHandler{}, database.RemoveNotificationParams{
		SenderUsername:      username,
		Reference:           params.Following.String(),
		NotificationsSource: "follow",
	})

	utils.JsonResponse(200, w, params)
}
