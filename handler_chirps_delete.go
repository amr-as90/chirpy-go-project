package main

import (
	"net/http"

	"github.com/amr-as90/chirpy-go-project/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//Get the token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	//Validate the token
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	//Extract the chirp ID from the URL and store it
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	//Get the chirp from the database and compare userID in JWT to userID in the Chirp
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	if userID != dbChirp.UserID {
		respondWithError(w, 403, "Unable to delete chirp, unauthorized.", err)
		return
	}

	//Make database call to delete the chirp, using the chirp ID and ?token?
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	w.WriteHeader(204)
}
