package main

import (
	"net/http"

	"github.com/amr-as90/chirpy-go-project/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token, err := auth.GetRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't retrieve token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Was unable to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
