package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	db "github.com/feproldo/effective-mobile/internal/db/generated"
	subscriptionHandler "github.com/feproldo/effective-mobile/internal/handlers/subscriptions"
	"github.com/feproldo/effective-mobile/internal/middlewares"
	subscriptionService "github.com/feproldo/effective-mobile/internal/services/subscriptions"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	godotenv.Load()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Error().Err(err).Msg("Database connection error")
		return
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(time.Hour)

	queries := db.New(conn)

	subsService := subscriptionService.NewService(queries)
	subsHandler := subscriptionHandler.NewHandler(subsService)

	router := chi.NewRouter()

	router.Use(middlewares.ZeroLogLogger)

	router.Route("/subscriptions", func(r chi.Router) {
		r.Get("/", subsHandler.List)
		r.Get("/{id}", subsHandler.Get)

		r.Get("/user/{user_id}", subsHandler.GetByUserId)

		r.Post("/", subsHandler.Create)

		r.Delete("/{id}", subsHandler.Create)
	})

	port := os.Getenv("PORT")
	log.Info().Msg("Service started on 0.0.0.0:" + port)
	http.ListenAndServe("0.0.0.0:"+port, router)
}
