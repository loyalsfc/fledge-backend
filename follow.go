package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type FollowParams struct {
	Following uuid.UUID `json:"following"`
	Follower  uuid.UUID `json:"follower"`
}

func (apiCfg apiCfg) follow(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := FollowParams{}

	decorder.Decode(&params)

	follow, err := apiCfg.DB.NewFollower(r.Context(), database.NewFollowerParams{
		FollowerID:  params.Follower,
		FollowingID: params.Following,
	})

	if err != nil {
		errResponse(405, w, fmt.Sprintf("Error occured %v", err))
		return
	}

	jsonResponse(200, w, follow)
}

func getIdFromParams(r *http.Request) (uuid.UUID, error) {
	id := r.URL.Query().Get("id")

	if id == "" {
		return uuid.New(), errors.New("no id found")
	}

	userId, err := uuid.Parse(id)

	if err != nil {
		return uuid.New(), errors.New("invalid userId")
	}

	return userId, nil
}

func (apiCfg apiCfg) getFollower(w http.ResponseWriter, r *http.Request, username string) {

	userId, err := getIdFromParams(r)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	list, err := apiCfg.DB.GetFollowers(r.Context(), userId)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	jsonResponse(200, w, handleFollowersToFollowers(list))
}

func (apiCfg apiCfg) getFollowing(w http.ResponseWriter, r *http.Request, username string) {

	userId, err := getIdFromParams(r)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	list, err := apiCfg.DB.GetFollowing(r.Context(), userId)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("Error occured %v ", err))
		return
	}

	jsonResponse(200, w, handleFollowsToFollows(list))
}
