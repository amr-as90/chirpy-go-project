package main

import (
	"net/http"
	"time"

	"github.com/amr-as90/chirpy-go-project/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token, err := auth.GetRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't retrieve token", err)
		return
	}

	u, err := cfg.db.GetValidRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Couldn't retrieve user from token", err)
		return
	}

	if u.RevokedAt.Valid {
		// If RevokedAt has a valid value (is not NULL), the token is revoked
		respondWithError(w, 401, "Token has been revoked", nil)
		return
	}

	// Create JWT token from returned user

	const expirationTime = 3600 * time.Second

	accessToken, err := auth.MakeJWT(
		u.UserID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	resp := response{
		Token: accessToken,
	}

	respondWithJSON(w, 200, resp)

}
