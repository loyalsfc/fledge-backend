package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/loyalsfc/fledge-backend/internal/database"

	_ "github.com/lib/pq"
)

type apiCfg struct {
	DB *database.Queries
}

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

	apiCfg := apiCfg{
		DB: db,
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	v1Router := chi.NewRouter()

	v1Router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connected"))
	})

	v1Router.Get("/users", apiCfg.getUsers)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.getUser))
	v1Router.Post("/user", apiCfg.createUser)
	v1Router.Post("/sign", apiCfg.userSignin)

	v1Router.Put("/update-profile-image", apiCfg.middlewareAuth(apiCfg.changeUserProfile))
	v1Router.Put("/update-cover-image", apiCfg.middlewareAuth(apiCfg.changeUserCoverImage))
	v1Router.Put("/update-user-profile", apiCfg.middlewareAuth(apiCfg.updateUserProfile))

	v1Router.Post("/follow", apiCfg.middlewareAuth(apiCfg.follow))
	v1Router.Post("/unfollow", apiCfg.middlewareAuth(apiCfg.unfollow))
	v1Router.Get("/get-followers", apiCfg.middlewareAuth(apiCfg.getFollower))
	v1Router.Get("/get-following", apiCfg.middlewareAuth(apiCfg.getFollowing))

	v1Router.Post("/new-post", apiCfg.middlewareAuth(apiCfg.makePost))

	router.Mount("/v1", v1Router)

	http.ListenAndServe(":3333", router)
}
