package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type BoookmarksParam struct {
	PostID string `json:"post_id"`
}

func (apiCfg apiCfg) addBookmarks(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := BoookmarksParam{}

	decorder.Decode(&params)

	postID, err := uuid.Parse(params.PostID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	bookmarkErr := apiCfg.DB.AddBookmarks(r.Context(), database.AddBookmarksParams{
		Username: username,
		PostID:   postID,
	})

	if bookmarkErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v", bookmarkErr))
		return
	}

	bookmarksCount, err := apiCfg.DB.UpdateBookmarksIncrease(r.Context(), postID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "successful",
		Message: "bookmark added",
		Payload: bookmarksCount,
	})
}

func (apiCfg apiCfg) removeBookmarks(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := BoookmarksParam{}

	decorder.Decode(&params)

	postID, err := uuid.Parse(params.PostID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	bookmarkErr := apiCfg.DB.RemoveBookmarks(r.Context(), database.RemoveBookmarksParams{
		Username: username,
		PostID:   postID,
	})

	if bookmarkErr != nil {
		errResponse(403, w, fmt.Sprintf("error %v", bookmarkErr))
		return
	}

	bookmarksCount, err := apiCfg.DB.UpdateBookmarksDecrease(r.Context(), postID)

	if err != nil {
		errResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	jsonResponse(200, w, Response{
		Status:  "successful",
		Message: "bookmark added",
		Payload: bookmarksCount,
	})
}
