package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "github.com/feproldo/effective-mobile/docs"
	db "github.com/feproldo/effective-mobile/internal/db/generated"
	subscriptionHandler "github.com/feproldo/effective-mobile/internal/handlers/subscriptions"
	"github.com/feproldo/effective-mobile/internal/middlewares"
	subscriptionService "github.com/feproldo/effective-mobile/internal/services/subscriptions"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title		Subscriptions service
// @version	1.0
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

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.Route("/subscriptions", func(r chi.Router) {
		r.Get("/", subsHandler.List)
		r.Get("/{id}", subsHandler.Get)

		r.Get("/user/{user_id}", subsHandler.GetByUserId)

		r.Post("/", subsHandler.Create)

		r.Delete("/{id}", subsHandler.Delete)

		r.Put("/{id}", subsHandler.Update)

		r.Get("/sum", subsHandler.Sum)
	})

	port := os.Getenv("PORT")

	log.Info().Msg("Service started on 0.0.0.0:" + port)
	log.Info().Msg("Swagger started on 0.0.0.0:" + port + "/swagger/index.html")

	http.ListenAndServe("0.0.0.0:"+port, router)
}
