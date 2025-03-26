package main

import (
	"encoding/json"
	"net/http"

	"github.com/amr-as90/chirpy-go-project/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type Data struct {
		UserID string `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	userApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't find API key", err)
		return
	}

	if userApiKey != cfg.PolkaAPIKey {
		respondWithError(w, 401, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, 204, "Invalid event", nil)
		return
	}

	userUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "Invalid user ID format", err)
		return
	}

	cfg.db.UpgradeUser(r.Context(), userUUID)

	respondWithJSON(w, 204, nil)
}
