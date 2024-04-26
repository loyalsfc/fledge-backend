package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type parameters struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInPayload struct {
	User        database.User
	AccessToken string
}

var secretKey = []byte("secret-key")

func (apiCfg *apiCfg) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsers(r.Context())

	if err != nil {
		errResponse(403, w, fmt.Sprintf("Error getting users %v", err))
	}

	jsonResponse(200, w, handleUsersToUsers(users))

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
		Password: passwordEncryption(params.Password),
	})

	if err != nil {
		errResponse(400, w, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	jsonResponse(200, w, handleUserToUser(user))
}

func (apiCfg apiCfg) userSignin(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)

	params := signin{}

	err := decorder.Decode(&params)

	if err != nil {
		errResponse(400, w, fmt.Sprintf("Error paring json %v", err))
	}

	user, err := apiCfg.DB.SignIn(r.Context(), params.Email)

	if err != nil {
		errResponse(401, w, "invalid email or password")
		return
	}

	isPasswordMatched := comparePassword(user.Password, params.Password)

	if !isPasswordMatched {
		errResponse(401, w, "invalid email or password")
		return
	}

	jwtString, err := createToken(user.Username)

	if err != nil {
		errResponse(500, w, "Internal error occured")
		return
	}

	payload := SignInPayload{
		User:        user,
		AccessToken: jwtString,
	}

	jsonResponse(200, w, handleLoginToLogin(payload))
}

func generateUsername(name string) string {
	// Split the name into first and last name
	parts := strings.Fields(name)
	if len(parts) < 2 {
		// Handle cases where there's no space between the first and last name
		// You can implement your own logic here
		return name
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error generating hashped password", err)
	}

	return string(hashedPassword)
}

func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		fmt.Println("Password doesn't match")
		return false
	}

	return true
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (apiCfg apiCfg) getUser(w http.ResponseWriter, r *http.Request, username string) {
	profileUsername := r.URL.Query().Get("username")
	user, err := apiCfg.DB.GetUser(r.Context(), profileUsername)

	if err != nil {
		errResponse(404, w, fmt.Sprintf("Error occured: %v", err))
		return
	}

	jsonResponse(200, w, handleUserToUser(user))
}
