package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type parameters struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (apiCfg *apiCfg) getUsers(w http.ResponseWriter, r *http.Request) {
	user, err := apiCfg.DB.GetUser(r.Context())

	if err != nil {
		errResponse(403, w, fmt.Sprintf("Error getting users %v", err))
	}

	jsonResponse(200, w, user)

}

func (apiCfg *apiCfg) createUser(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decorder.Decode(&params)

	if err != nil {
		errResponse(500, w, "invalid parameters")
		return
	}

	username := generateUsername(params.Name)

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Name:     params.Name,
		Email:    params.Email,
		Username: username,
		Password: params.Password,
	})

	if err != nil {
		errResponse(400, w, fmt.Sprintf("Error parsing json %v", err))
	}

	jsonResponse(200, w, user)
}

func generateUsername(name string) string {
	// Split the name into first and last name
	parts := strings.Fields(name)
	if len(parts) < 2 {
		// Handle cases where there's no space between the first and last name
		// You can implement your own logic here
		return ""
	}
	firstName := strings.ToLower(parts[0])
	lastName := strings.ToLower(parts[1])

	// Generate the username by concatenating the first letter of the first name with the last name
	username := string(firstName[0]) + lastName

	// You might want to ensure the username is unique in your system
	// You can append a unique identifier if necessary
	// For example:
	// username += "_" + generateUniqueIdentifier()

	return username
}

func passwordEncryption(password string) string {
	// salt, _ := bcrypt.salt(10)
	return ""
}
