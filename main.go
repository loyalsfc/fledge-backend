package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/loyalsfc/fledge-backend/controllers/block"
	"github.com/loyalsfc/fledge-backend/controllers/bookmarks"
	"github.com/loyalsfc/fledge-backend/controllers/comment"
	"github.com/loyalsfc/fledge-backend/controllers/follow"
	"github.com/loyalsfc/fledge-backend/controllers/likes"
	"github.com/loyalsfc/fledge-backend/controllers/notification"
	"github.com/loyalsfc/fledge-backend/controllers/post"
	"github.com/loyalsfc/fledge-backend/controllers/replies"
	"github.com/loyalsfc/fledge-backend/controllers/user"
	"github.com/loyalsfc/fledge-backend/internal/database"
	"github.com/loyalsfc/fledge-backend/middlewares"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	dbString := os.Getenv("DB_STRING")

	conn, err := sql.Open("postgres", dbString)

	if err != nil {
		log.Fatal("Database error", err)
	}

	db := database.New(conn)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	v1Router := chi.NewRouter()

	v1Router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connected"))
	})

	middlewareAuth := middlewares.MiddlewareHandler{
		DB: db,
	}

	userRoutes := user.ApiCfg{
		DB: db,
	}

	v1Router.Get("/users", userRoutes.GetUsers)
	v1Router.Get("/user/{username}", userRoutes.GetUser)
	v1Router.Post("/user", userRoutes.CreateUser)
	v1Router.Post("/sign", userRoutes.Signin)
	v1Router.Get("/suggested-users", middlewareAuth.MiddlewareAuth(userRoutes.SuggestedUsers))

	v1Router.Put("/update-profile-image", middlewareAuth.MiddlewareAuth(userRoutes.ChangeUserProfile))
	v1Router.Put("/update-cover-image", middlewareAuth.MiddlewareAuth(userRoutes.ChangeUserCoverImage))
	v1Router.Put("/update-user-profile", middlewareAuth.MiddlewareAuth(userRoutes.UpdateUserProfile))

	followHandler := follow.FollowHandler{
		DB: db,
	}
	v1Router.Post("/follow", middlewareAuth.MiddlewareAuth(followHandler.Follow))
	v1Router.Post("/unfollow", middlewareAuth.MiddlewareAuth(followHandler.Unfollow))
	v1Router.Get("/get-followers/{userID}", followHandler.GetFollower)
	v1Router.Get("/get-following/{userID}", followHandler.GetFollowing)

	postHandler := post.PostHandler{
		DB: db,
	}
	v1Router.Post("/new-post", middlewareAuth.MiddlewareAuth(postHandler.MakePost))
	v1Router.Get("/user-posts/{username}", middlewareAuth.MiddlewareAuth(postHandler.GetUserPosts))
	v1Router.Get("/post/{postID}", postHandler.GetPost)
	v1Router.Get("/feeds", middlewareAuth.MiddlewareAuth(postHandler.GetUserFeeds))
	v1Router.Delete("/post/{postID}", middlewareAuth.MiddlewareAuth(postHandler.DeletePost))
	v1Router.Put("/post/{postID}", middlewareAuth.MiddlewareAuth(postHandler.EditPost))

	likesHandler := likes.LikeHandler{
		DB: db,
	}
	v1Router.Post("/like-post", middlewareAuth.MiddlewareAuth(likesHandler.LikePost))
	v1Router.Post("/unlike-post", middlewareAuth.MiddlewareAuth(likesHandler.UnlikePost))

	commentHandler := comment.CommentHandler{
		DB: db,
	}
	v1Router.Post("/new-comment", middlewareAuth.MiddlewareAuth(commentHandler.PostComment))
	v1Router.Delete("/comment/{commentID}", middlewareAuth.MiddlewareAuth(commentHandler.DeleteComment))
	v1Router.Get("/post-comments/{postID}", commentHandler.GetPostComments)
	v1Router.Post("/like-comment/{postID}", middlewareAuth.MiddlewareAuth(commentHandler.LikeComment))
	v1Router.Post("/unlike-comment/{postID}", middlewareAuth.MiddlewareAuth(commentHandler.UnLikeComment))
	v1Router.Put("/comment/{commentID}", middlewareAuth.MiddlewareAuth(commentHandler.EditComment))

	replyHandler := replies.ReplyHandler{
		DB: db,
	}
	v1Router.Get("/comment-replies/{commentID}", replyHandler.GetReplies)
	v1Router.Post("/reply-comment", middlewareAuth.MiddlewareAuth(replyHandler.ReplyComment))
	v1Router.Delete("/reply/{replyID}", middlewareAuth.MiddlewareAuth(replyHandler.DeleteCommetReply))
	v1Router.Post("/like-reply/{replyID}", middlewareAuth.MiddlewareAuth(replyHandler.LikeReply))
	v1Router.Post("/unlike-reply/{replyID}", middlewareAuth.MiddlewareAuth(replyHandler.UnLikeReply))
	v1Router.Put("/reply/{replyID}", middlewareAuth.MiddlewareAuth(replyHandler.EditReply))

	bookmarkHandler := bookmarks.BookMarkHandler{
		DB: db,
	}
	v1Router.Post("/add-bookmarks", middlewareAuth.MiddlewareAuth(bookmarkHandler.AddBookmarks))
	v1Router.Post("/remove-bookmarks", middlewareAuth.MiddlewareAuth(bookmarkHandler.RemoveBookmarks))
	v1Router.Get("/bookmarks", middlewareAuth.MiddlewareAuth(postHandler.GetBookmarkedPosts))

	notificationHandler := notification.NotificationHandler{
		DB: db,
	}
	v1Router.Get("/notifications", middlewareAuth.MiddlewareAuth(notificationHandler.GetUserNotifications))
	v1Router.Put("/update-notification", middlewareAuth.MiddlewareAuth(notificationHandler.MarkNotificationAsRead))

	blockHandler := block.BlockHandler{
		DB: db,
	}

	v1Router.Post("/block", middlewareAuth.MiddlewareAuth(blockHandler.Block))
	v1Router.Post("/unblock", middlewareAuth.MiddlewareAuth(blockHandler.UnBlock))
	v1Router.Get("/block-list/{userID}", middlewareAuth.MiddlewareAuth(blockHandler.GetBlocks))

	router.Mount("/v1", v1Router)

	http.ListenAndServe(":3333", router)
}
