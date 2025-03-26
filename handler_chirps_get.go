package main

import (
	"net/http"
	"sort"

	"github.com/amr-as90/chirpy-go-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get query parameters
	authorIDParam := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	// Default sort order is ascending
	isDescending := sortOrder == "desc"

	var dbChirps []database.Chirp
	var err error

	// Check if filtering by author
	if authorIDParam != "" {
		authorID, err := uuid.Parse(authorIDParam)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		dbChirps, err = cfg.db.GetUserChirps(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	} else {
		// Get all chirps if no author filter
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	}

	// Sort the chirps based on sort parameter
	if isDescending {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
		})
	} else {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt)
		})
	}

	// Convert database chirps to API chirps
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	// Everything went ok
	respondWithJSON(w, http.StatusOK, chirps)
}
