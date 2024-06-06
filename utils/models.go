package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

type SignInPayload struct {
	User        database.User
	AccessToken string
}

func HandleUserToUser(dbUser database.User) (user User) {
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

func HandleUsersToUsers(dbUsers []database.User) (users []User) {
	convertedUsers := []User{}

	for _, element := range dbUsers {
		convertedUsers = append(convertedUsers, HandleUserToUser(element))
	}

	return convertedUsers
}

func HandleLoginToLogin(signInPayload SignInPayload) (signInResponse SignInJson) {
	return SignInJson{
		User:        HandleUserToUser(signInPayload.User),
		AccessToken: signInPayload.AccessToken,
	}
}

func HandleFollowToFollow(followList database.GetFollowingRow) (user User) {
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

func HandleFollowsToFollows(followList []database.GetFollowingRow) (users []User) {
	usersList := []User{}

	for _, userList := range followList {
		usersList = append(usersList, HandleFollowToFollow(userList))
	}

	return usersList
}

func HandleFollowerToFollower(followList database.GetFollowersRow) (user User) {
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

func HandleFollowersToFollowers(followList []database.GetFollowersRow) (users []User) {
	usersList := []User{}

	for _, userList := range followList {
		usersList = append(usersList, HandleFollowerToFollower(userList))
	}

	return usersList
}

type Post struct {
	ID                      uuid.UUID       `json:"id"`
	UserID                  uuid.UUID       `json:"user_id"`
	Content                 string          `json:"content"`
	Media                   json.RawMessage `json:"media"`
	Username                string          `json:"username"`
	CreatedAt               string          `json:"created_at"`
	UpdatedAt               string          `json:"updated_at"`
	LikesCount              int             `json:"likes_count"`
	CommentCount            int             `json:"comments_count"`
	BookmarksCount          int             `json:"bookmarks_count"`
	ShareCount              int             `json:"shared_count"`
	Name                    string          `json:"name"`
	ProfilePicture          string          `json:"profile_picture"`
	IsVerified              bool            `json:"is_verified"`
	LikedUsersUsername      interface{}     `json:"liked_users"`
	BookmarkedUsersUsername interface{}     `json:"bookmarked_users"`
	IsSharedPost            bool            `json:"is_shared"`
	SharedPostID            uuid.UUID       `json:"shared_post_id"`
}

func convertUsernamesToString(usernames interface{}) []string {
	byteArray, ok := usernames.([]uint8)

	var result string

	if !ok {
		fmt.Println("failed to convert error")
	} else {
		for _, item := range byteArray {
			result += string(item)
		}
	}
	result = strings.Trim(result, "{}")

	if result == "null" {
		return []string{}
	}

	users := strings.Split(result, ",")

	return users
}

func HandlePostToPost(dbPost database.GetPostRow) (post Post) {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	return Post{
		ID:                      dbPost.ID,
		UserID:                  dbPost.UserID,
		Content:                 dbPost.Content,
		Media:                   dbPost.Media,
		Username:                dbPost.Username,
		CreatedAt:               dbPost.CreatedAt.Format(JavascriptISOString),
		UpdatedAt:               dbPost.UpdatedAt.Format(JavascriptISOString),
		LikesCount:              int(dbPost.LikesCount),
		CommentCount:            int(dbPost.CommentCount),
		BookmarksCount:          int(dbPost.BookmarksCount),
		ShareCount:              int(dbPost.ShareCount),
		Name:                    dbPost.Name,
		ProfilePicture:          dbPost.ProfilePicture.String,
		IsVerified:              dbPost.IsVerified.Bool,
		LikedUsersUsername:      convertUsernamesToString(dbPost.LikedUsersUsername),
		BookmarkedUsersUsername: convertUsernamesToString(dbPost.BookmarkedUsersUsername),
		IsSharedPost:            dbPost.IsSharedPost,
		SharedPostID:            dbPost.SharedPostID.UUID,
	}
}

func HandlePostsToPosts(dbPosts []database.GetUserPostsRow) (posts []Post) {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	initPosts := []Post{}

	for _, post := range dbPosts {
		initPosts = append(initPosts, Post{
			ID:                      post.ID,
			UserID:                  post.UserID,
			Content:                 post.Content,
			Media:                   post.Media,
			Username:                post.Username,
			CreatedAt:               post.CreatedAt.Format(JavascriptISOString),
			UpdatedAt:               post.UpdatedAt.Format(JavascriptISOString),
			LikesCount:              int(post.LikesCount),
			CommentCount:            int(post.CommentCount),
			BookmarksCount:          int(post.BookmarksCount),
			ShareCount:              int(post.ShareCount),
			Name:                    post.Name,
			ProfilePicture:          post.ProfilePicture.String,
			IsVerified:              post.IsVerified.Bool,
			LikedUsersUsername:      convertUsernamesToString(post.LikedUsersUsername),
			BookmarkedUsersUsername: convertUsernamesToString(post.BookmarkedUsersUsername),
			IsSharedPost:            post.IsSharedPost,
			SharedPostID:            post.SharedPostID.UUID,
		})
	}

	return initPosts
}

func HandleFeedsToFeeds(dbPosts []database.GetFeedPostsRow) (posts []Post) {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	initPosts := []Post{}

	for _, post := range dbPosts {
		initPosts = append(initPosts, Post{
			ID:                      post.ID,
			UserID:                  post.UserID,
			Content:                 post.Content,
			Media:                   post.Media,
			Username:                post.Username,
			CreatedAt:               post.CreatedAt.Format(JavascriptISOString),
			UpdatedAt:               post.UpdatedAt.Format(JavascriptISOString),
			LikesCount:              int(post.LikesCount),
			CommentCount:            int(post.CommentCount),
			BookmarksCount:          int(post.BookmarksCount),
			ShareCount:              int(post.ShareCount),
			Name:                    post.Name,
			ProfilePicture:          post.ProfilePicture.String,
			IsVerified:              post.IsVerified.Bool,
			LikedUsersUsername:      convertUsernamesToString(post.LikedUsersUsername),
			BookmarkedUsersUsername: convertUsernamesToString(post.BookmarkedUsersUsername),
			IsSharedPost:            post.IsSharedPost,
			SharedPostID:            post.SharedPostID.UUID,
		})
	}

	return initPosts
}

type Comment struct {
	ID                 uuid.UUID       `json:"id"`
	CommentText        string          `json:"comment_text"`
	Media              json.RawMessage `json:"media"`
	Username           string          `json:"username"`
	PostID             uuid.UUID       `json:"post_id"`
	LikesCount         int32           `json:"likes_count"`
	ReplyCount         int32           `json:"reply_count"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	Name               string          `json:"name"`
	ProfilePicture     string          `json:"profile_picture"`
	IsVerified         bool            `json:"is_verified"`
	LikedUsersUsername interface{}     `json:"liked_users"`
}

func HandleCommentToComment(dbComment database.GetCommentsRow) (comment Comment) {
	return Comment{
		ID:                 dbComment.ID,
		CommentText:        dbComment.CommentText,
		Media:              dbComment.Media,
		Username:           dbComment.Username,
		PostID:             dbComment.PostID,
		LikesCount:         dbComment.LikesCount,
		ReplyCount:         dbComment.ReplyCount,
		CreatedAt:          dbComment.CreatedAt,
		UpdatedAt:          dbComment.UpdatedAt,
		Name:               dbComment.Name,
		ProfilePicture:     dbComment.ProfilePicture.String,
		IsVerified:         dbComment.IsVerified.Bool,
		LikedUsersUsername: convertUsernamesToString(dbComment.LikedUsersUsername),
	}
}

func HandleCommentsToComments(dbComments []database.GetCommentsRow) (comments []Comment) {
	commentConverts := []Comment{}

	for _, dbComment := range dbComments {
		commentConverts = append(commentConverts, HandleCommentToComment(dbComment))
	}

	return commentConverts
}

func HandleBookmarksToBookmarks(dbPosts []database.GetBookmarkedPostsRow) []Post {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	posts := []Post{}

	for _, dbPost := range dbPosts {
		posts = append(posts, Post{
			ID:                      dbPost.ID,
			UserID:                  dbPost.UserID,
			Content:                 dbPost.Content,
			Media:                   dbPost.Media,
			Username:                dbPost.Username,
			CreatedAt:               dbPost.CreatedAt.Format(JavascriptISOString),
			UpdatedAt:               dbPost.UpdatedAt.Format(JavascriptISOString),
			LikesCount:              int(dbPost.LikesCount),
			CommentCount:            int(dbPost.CommentCount),
			BookmarksCount:          int(dbPost.BookmarksCount),
			ShareCount:              int(dbPost.ShareCount),
			Name:                    dbPost.Name,
			ProfilePicture:          dbPost.ProfilePicture.String,
			IsVerified:              dbPost.IsVerified.Bool,
			LikedUsersUsername:      convertUsernamesToString(dbPost.LikedUsersUsername),
			BookmarkedUsersUsername: convertUsernamesToString(dbPost.BookmarkedUsersUsername),
			IsSharedPost:            dbPost.IsSharedPost,
			SharedPostID:            dbPost.SharedPostID.UUID,
		})
	}

	return posts
}

type Notification struct {
	ID                  uuid.UUID `json:"id"`
	SenderUsername      string    `json:"username"`
	Content             string    `json:"content"`
	CreatedAt           string    `json:"created_at"`
	IsViewed            bool      `json:"is_viewed"`
	NotificationsSource string    `json:"notification_source"`
	Reference           string    `json:"reference"`
	Name                string    `json:"name"`
	ProfilePicture      string    `json:"profile_picture"`
	IsVerified          bool      `json:"is_verified"`
}

func HandleNotificationsToNotifications(dbNotifications []database.GetUserNotificationsRow) []Notification {
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	notifications := []Notification{}

	for _, dbNotification := range dbNotifications {
		notifications = append(notifications, Notification{
			ID:                  dbNotification.ID,
			SenderUsername:      dbNotification.SenderUsername,
			Content:             dbNotification.Content,
			CreatedAt:           dbNotification.CreatedAt.Format(JavascriptISOString),
			IsViewed:            dbNotification.IsViewed,
			NotificationsSource: dbNotification.NotificationsSource,
			Reference:           dbNotification.Reference,
			Name:                dbNotification.Name,
			ProfilePicture:      dbNotification.ProfilePicture.String,
			IsVerified:          dbNotification.IsVerified.Bool,
		})
	}

	return notifications
}

func HandleRepliesToReplies(dbReplies []database.GetRepliesRow) (reply []Comment) {
	replies := []Comment{}

	for _, reply := range dbReplies {
		replies = append(replies, Comment{
			ID:                 reply.ID,
			CommentText:        reply.ReplyText,
			Media:              reply.Media,
			Username:           reply.Username,
			PostID:             reply.CommentID,
			LikesCount:         reply.LikesCount,
			CreatedAt:          reply.CreatedAt,
			UpdatedAt:          reply.UpdatedAt,
			Name:               reply.Name,
			ProfilePicture:     reply.ProfilePicture.String,
			IsVerified:         reply.IsVerified.Bool,
			LikedUsersUsername: convertUsernamesToString(reply.LikedUsersUsername),
		})
	}

	return replies
}
