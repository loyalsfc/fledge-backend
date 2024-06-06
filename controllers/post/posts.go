package post

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

type PostParams struct {
	UserID       uuid.UUID       `json:"user_id"`
	Content      string          `json:"content"`
	Media        json.RawMessage `json:"media"`
	IsShared     bool            `json:"is_shared"`
	SharedPostId uuid.UUID       `json:"shared_post_id"`
}

type PostHandler struct {
	DB *database.Queries
}

func (apiCfg PostHandler) MakePost(w http.ResponseWriter, r *http.Request, username string) {
	decoder := json.NewDecoder(r.Body)
	params := PostParams{}

	decoder.Decode(&params)

	postId, err := apiCfg.DB.NewPost(r.Context(), database.NewPostParams{
		ID:           uuid.New(),
		UserID:       params.UserID,
		Username:     username,
		Content:      params.Content,
		Media:        params.Media,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		IsSharedPost: params.IsShared,
		SharedPostID: uuid.NullUUID{Valid: params.IsShared, UUID: params.SharedPostId},
	})

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("Error %v ", err))
		return
	}

	post, err := apiCfg.DB.GetPost(r.Context(), postId)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("Error: %v", err))
		return
	}

	if params.IsShared {
		apiCfg.DB.IncreaseShareCount(r.Context(), params.SharedPostId)
	}

	utils.JsonResponse(200, w, utils.HandlePostToPost(post))
}

func (apiCfg PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request, username string) {
	profileUsername := chi.URLParam(r, "username")

	if profileUsername == "" {
		utils.ErrResponse(200, w, "Invalid username")
		return
	}

	posts, err := apiCfg.DB.GetUserPosts(r.Context(), profileUsername)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("Error: %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandlePostsToPosts(posts))

}

func (apiCfg PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postID")

	if postId == "" {
		utils.ErrResponse(404, w, "post id not found")
		return
	}

	id, err := uuid.Parse(postId)

	if err != nil {
		utils.ErrResponse(403, w, "invalid post id")
		return
	}

	post, err := apiCfg.DB.GetPost(r.Context(), id)

	if err != nil {
		utils.ErrResponse(400, w, fmt.Sprintf("error %v ", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandlePostToPost(post))
}

func (apiCfg PostHandler) GetUserFeeds(w http.ResponseWriter, r *http.Request, username string) {
	var idString string = r.URL.Query().Get("id")

	if idString == "" {
		utils.ErrResponse(401, w, "no id found")
		return
	}

	userID, err := uuid.Parse(idString)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error: %v", err))
		return
	}

	feeds, err := apiCfg.DB.GetFeedPosts(r.Context(), userID)

	if err != nil {
		utils.ErrResponse(401, w, fmt.Sprintf("error: %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleFeedsToFeeds(feeds))
}

func (apiCfg PostHandler) GetBookmarkedPosts(w http.ResponseWriter, r *http.Request, username string) {
	posts, err := apiCfg.DB.GetBookmarkedPosts(r.Context(), username)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error: %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.HandleBookmarksToBookmarks(posts))
}

func (apiCfg PostHandler) DeletePost(w http.ResponseWriter, r *http.Request, username string) {
	idString := chi.URLParam(r, "postID")

	postId, err := uuid.Parse(idString)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	deleteErr := apiCfg.DB.DeletePost(r.Context(), postId)

	if deleteErr != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "post delete successfully",
	})
}

func (apiCfg PostHandler) EditPost(w http.ResponseWriter, r *http.Request, username string) {
	idString := chi.URLParam(r, "postID")

	postId, err := uuid.Parse(idString)

	if err != nil {
		utils.ErrResponse(403, w, fmt.Sprintf("error %v", err))
		return
	}

	decorder := json.NewDecoder(r.Body)

	params := PostParams{}

	decorder.Decode(&params)

	postErr := apiCfg.DB.EditPost(r.Context(), database.EditPostParams{
		Content: params.Content,
		Media:   params.Media,
		ID:      postId,
	})

	if postErr != nil {
		utils.ErrResponse(404, w, fmt.Sprintf("error %v", postErr))
		return
	}

	utils.JsonResponse(200, w, utils.Response{
		Status:  "success",
		Message: "post edit successful",
	})
}
