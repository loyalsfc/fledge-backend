package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		fmt.Println("Password doesn't match")
		return false
	}

	return true
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type LikesParam struct {
	PostId string `json:"post_id"`
}

func GetPostId(r *http.Request) (uuid.UUID, error) {
	decorder := json.NewDecoder(r.Body)

	params := LikesParam{}

	decorder.Decode(&params)

	postId, err := uuid.Parse(params.PostId)

	if err != nil {
		return uuid.New(), errors.New("invalid post id")
	}

	return postId, nil
}

func TrimToMaxChars(s string, maxChars int) string {
	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars]
}
