package bookmarks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type BoookmarksParam struct {
	PostID string `json:"post_id"`
}

type BookMarkHandler struct {
	DB *database.Queries
}

func (apiCfg BookMarkHandler) AddBookmarks(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := BoookmarksParam{}

	decorder.Decode(&params)

	postID, err := uuid.Parse(params.PostID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	bookmarkErr := apiCfg.DB.AddBookmarks(r.Context(), database.AddBookmarksParams{
		Username: username,
		PostID:   postID,
	})

	if bookmarkErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", bookmarkErr))
		return
	}

	bookmarksCount, err := apiCfg.DB.UpdateBookmarksIncrease(r.Context(), postID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "successful",
		Message: "bookmark added",
		Payload: bookmarksCount,
	})
}

func (apiCfg BookMarkHandler) RemoveBookmarks(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := BoookmarksParam{}

	decorder.Decode(&params)

	postID, err := uuid.Parse(params.PostID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	bookmarkErr := apiCfg.DB.RemoveBookmarks(r.Context(), database.RemoveBookmarksParams{
		Username: username,
		PostID:   postID,
	})

	if bookmarkErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", bookmarkErr))
		return
	}

	bookmarksCount, err := apiCfg.DB.UpdateBookmarksDecrease(r.Context(), postID)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "successful",
		Message: "bookmark added",
		Payload: bookmarksCount,
	})
}
