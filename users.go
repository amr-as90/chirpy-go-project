package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/amr-as90/chirpy-go-project/internal/auth"
	"github.com/amr-as90/chirpy-go-project/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Validate e-mail isn't empty
	if user.Email == "" {
		http.Error(w, "E-mail is invalid.", http.StatusBadRequest)
		return
	}
	//Validate e-mail contains an @
	if !strings.Contains(user.Email, "@") {
		http.Error(w, "E-mail is invalid.", http.StatusBadRequest)
		return
	}
	//Validate Password isn't empty
	if user.Password == "" {
		http.Error(w, "Password cannot be empty.", http.StatusBadRequest)
		return
	}

	//If pwd is valid, hash it
	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Something went wrong. Password cannot be hashed.", http.StatusInternalServerError)
		return
	}

	userParams := database.CreateUserParams{
		Email:          user.Email,
		HashedPassword: user.Password,
	}

	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), userParams)
	if err != nil {
		http.Error(w, "Something went wrong creating the user", http.StatusInternalServerError)
		return
	}

	response := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Something went wrong encoding JSON", http.StatusInternalServerError)
		return
	}

}

func (cfg *apiConfig) handlerAuthenticate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, 500, "Something went wrong decoding JSON", err)
	}
	// Get the user from the database and store as returnedUser
	returnedUser, err := cfg.dbQueries.GetUser(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(returnedUser.HashedPassword, user.Password)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	response := User{
		ID:        returnedUser.ID,
		CreatedAt: returnedUser.CreatedAt,
		UpdatedAt: returnedUser.UpdatedAt,
		Email:     returnedUser.Email,
	}

	respondWithJSON(w, 200, response)

}
