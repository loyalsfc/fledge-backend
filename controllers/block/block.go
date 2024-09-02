package block

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/utils"
)

type BlockParams struct {
	BlockerId uuid.UUID `json:"blocker_id"`
	BlockedId uuid.UUID `json:"blocked_id"`
}

type BlockHandler struct {
	DB *database.Queries
}

func (b BlockHandler) Block(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := &BlockParams{}

	decorder.Decode(&params)

	err := b.DB.Block(r.Context(), database.BlockParams{
		BlockerID: params.BlockerId,
		BlockedID: params.BlockedId,
	})

	if err != nil {
		utils.ErrResponse(401, w, err.Error())
		return
	}

	utils.JsonResponse(200, w, nil)

}

func (b BlockHandler) UnBlock(w http.ResponseWriter, r *http.Request, username string) {
	decorder := json.NewDecoder(r.Body)

	params := &BlockParams{}

	decorder.Decode(&params)

	err := b.DB.Unblock(r.Context(), database.UnblockParams{
		BlockerID: params.BlockerId,
		BlockedID: params.BlockedId,
	})

	if err != nil {
		utils.ErrResponse(401, w, err.Error())
		return
	}

	utils.JsonResponse(200, w, nil)

}

func (b BlockHandler) GetBlocks(w http.ResponseWriter, r *http.Request, username string) {
	param := chi.URLParam(r, "userID")

	userId, err := uuid.Parse(param)

	if err != nil {
		utils.ErrResponse(401, w, err.Error())
		return
	}

	blockedUsers, err := b.DB.GetBlocked(r.Context(), userId)

	if err != nil {
		utils.ErrResponse(401, w, err.Error())
		return
	}

	utils.JsonResponse(200, w, blockedUsers)
}
