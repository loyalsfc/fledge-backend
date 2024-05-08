package main

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/loyalsfc/fledge-backend/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Bio            string    `json:"bio"`
	Profession     string    `json:"profession"`
	IsVerified     bool      `json:"is_verified"`
	IsActive       bool      `json:"is_active"`
	ProfilePicture string    `json:"profile_picture"`
	CoverPicture   string    `json:"cover_picture"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

type SignInJson struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func handleUserToUser(dbUser database.User) (user User) {
	return User{
		ID:             dbUser.ID,
		Name:           dbUser.Name,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		Bio:            dbUser.Bio.String,
		Profession:     dbUser.Profession.String,
		IsVerified:     dbUser.IsVerified.Bool,
		IsActive:       dbUser.IsActive.Bool,
		ProfilePicture: dbUser.ProfilePicture.String,
		CoverPicture:   dbUser.CoverPicture.String,
		CreatedAt:      dbUser.CreatedAt.Time.String(),
		UpdatedAt:      dbUser.UpdatedAt.Time.String(),
	}
}

func handleUsersToUsers(dbUsers []database.User) (users []User) {
	convertedUsers := []User{}

	for _, element := range dbUsers {
		convertedUsers = append(convertedUsers, handleUserToUser(element))
	}

	return convertedUsers
}

func handleLoginToLogin(signInPayload SignInPayload) (signInResponse SignInJson) {
	return SignInJson{
		User:        handleUserToUser(signInPayload.User),
		AccessToken: signInPayload.AccessToken,
	}
}

func handleFollowToFollow(followList database.GetFollowingRow) (user User) {
	return User{
		ID:             followList.ID,
		Name:           followList.Name,
		Username:       followList.Username,
		Email:          followList.Email,
		Bio:            followList.Bio.String,
		Profession:     followList.Profession.String,
		IsVerified:     followList.IsVerified.Bool,
		IsActive:       followList.IsActive.Bool,
		ProfilePicture: followList.ProfilePicture.String,
		CoverPicture:   followList.CoverPicture.String,
		CreatedAt:      followList.CreatedAt.Time.String(),
		UpdatedAt:      followList.UpdatedAt.Time.String(),
	}
}

func handleFollowsToFollows(followList []database.GetFollowingRow) (users []User) {
	usersList := []User{}

	for _, userList := range followList {
		usersList = append(usersList, handleFollowToFollow(userList))
	}

	return usersList
}

func handleFollowerToFollower(followList database.GetFollowersRow) (user User) {
	return User{
		ID:             followList.ID,
		Name:           followList.Name,
		Username:       followList.Username,
		Email:          followList.Email,
		Bio:            followList.Bio.String,
		Profession:     followList.Profession.String,
		IsVerified:     followList.IsVerified.Bool,
		IsActive:       followList.IsActive.Bool,
		ProfilePicture: followList.ProfilePicture.String,
		CoverPicture:   followList.CoverPicture.String,
		CreatedAt:      followList.CreatedAt.Time.String(),
		UpdatedAt:      followList.UpdatedAt.Time.String(),
	}
}

func handleFollowersToFollowers(followList []database.GetFollowersRow) (users []User) {
	usersList := []User{}

	for _, userList := range followList {
		usersList = append(usersList, handleFollowerToFollower(userList))
	}

	return usersList
}

type Post struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Content   string          `json:"content"`
	Media     json.RawMessage `json:"media"`
	Username  string          `json:"username"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func handlePostToPost(dbPost database.Post) (post Post) {
	return Post{
		ID:        dbPost.ID,
		UserID:    dbPost.UserID,
		Content:   dbPost.Content,
		Media:     dbPost.Media,
		Username:  dbPost.Username,
		CreatedAt: dbPost.CreatedAt.GoString(),
		UpdatedAt: dbPost.UpdatedAt.GoString(),
	}
}
