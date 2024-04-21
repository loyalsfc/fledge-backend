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
	v1Router.Post("/user", apiCfg.createUser)

	router.Mount("/v1", v1Router)

	http.ListenAndServe(":3333", router)
}
