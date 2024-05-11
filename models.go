package main

import (
	"encoding/json"
	"fmt"

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
	ID                 uuid.UUID       `json:"id"`
	UserID             uuid.UUID       `json:"user_id"`
	Content            string          `json:"content"`
	Media              json.RawMessage `json:"media"`
	Username           string          `json:"username"`
	CreatedAt          string          `json:"created_at"`
	UpdatedAt          string          `json:"updated_at"`
	LikesCount         int             `json:"likes_count"`
	CommentCount       int             `json:"comments_count"`
	BookmarksCount     int             `json:"bookmarks_count"`
	ShareCount         int             `json:"shared_count"`
	Name               string          `json:"name"`
	ProfilePicture     string          `json:"profile_picture"`
	IsVerified         bool            `json:"is_verified"`
	LikedUsersUsername interface{}     `json:"liked_users"`
}

func handlePostToPost(dbPost database.GetPostRow) (post Post) {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	return Post{
		ID:             dbPost.ID,
		UserID:         dbPost.UserID,
		Content:        dbPost.Content,
		Media:          dbPost.Media,
		Username:       dbPost.Username,
		CreatedAt:      dbPost.CreatedAt.Format(JavascriptISOString),
		UpdatedAt:      dbPost.UpdatedAt.Format(JavascriptISOString),
		LikesCount:     int(dbPost.LikesCount),
		CommentCount:   int(dbPost.CommentCount),
		BookmarksCount: int(dbPost.BookmarksCount),
		ShareCount:     int(dbPost.ShareCount),
		Name:           dbPost.Name,
		ProfilePicture: dbPost.ProfilePicture.String,
		IsVerified:     dbPost.IsVerified.Bool,
	}
}

func handlePostsToPosts(dbPosts []database.GetUserPostsRow) (posts []Post) {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	initPosts := []Post{}

	for _, post := range dbPosts {

		usernames := post.LikedUsersUsername

		byteArray, ok := usernames.([]uint8)

		var result string

		if !ok {
			fmt.Println("failed to convert error")
		} else {
			for _, item := range byteArray {
				result += string(item)
			}
		}

		initPosts = append(initPosts, Post{
			ID:                 post.ID,
			UserID:             post.UserID,
			Content:            post.Content,
			Media:              post.Media,
			Username:           post.Username,
			CreatedAt:          post.CreatedAt.Format(JavascriptISOString),
			UpdatedAt:          post.UpdatedAt.Format(JavascriptISOString),
			LikesCount:         int(post.LikesCount),
			CommentCount:       int(post.CommentCount),
			BookmarksCount:     int(post.BookmarksCount),
			ShareCount:         int(post.ShareCount),
			Name:               post.Name,
			ProfilePicture:     post.ProfilePicture.String,
			IsVerified:         post.IsVerified.Bool,
			LikedUsersUsername: result,
		})
	}

	return initPosts
}

func handleFeedsToFeeds(dbPosts []database.GetFeedPostsRow) (posts []Post) {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	initPosts := []Post{}

	for _, post := range dbPosts {

		usernames := post.LikedUsersUsername

		byteArray, ok := usernames.([]uint8)

		var result string

		if !ok {
			fmt.Println("failed to convert error")
		} else {
			for _, item := range byteArray {
				result += string(item)
			}
		}

		initPosts = append(initPosts, Post{
			ID:                 post.ID,
			UserID:             post.UserID,
			Content:            post.Content,
			Media:              post.Media,
			Username:           post.Username,
			CreatedAt:          post.CreatedAt.Format(JavascriptISOString),
			UpdatedAt:          post.UpdatedAt.Format(JavascriptISOString),
			LikesCount:         int(post.LikesCount),
			CommentCount:       int(post.CommentCount),
			BookmarksCount:     int(post.BookmarksCount),
			ShareCount:         int(post.ShareCount),
			Name:               post.Name,
			ProfilePicture:     post.ProfilePicture.String,
			IsVerified:         post.IsVerified.Bool,
			LikedUsersUsername: result,
		})
	}

	return initPosts
}
