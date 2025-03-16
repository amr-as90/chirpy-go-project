package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/amr-as90/chirpy-go-project/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	createParams := struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}{
		Body:   ProfanityCheck(params.Body),
		UserID: params.UserID,
	}

	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams(createParams))
	if err != nil {
		http.Error(w, "Something went wrong creating chirp", http.StatusInternalServerError)
		return
	}

	response := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var allChirps []database.Chirp

	allChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, "Couldn't get all items in database", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, allChirps)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirpID := r.PathValue("chirpID")
	chirpUuid, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 500, "Couldn't convert chirpID to UUID", err)
		return
	}

	returnChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, 404, "Not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnChirp)

}
