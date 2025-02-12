package main

import (
	"jeopardy-game/internal/config"
	"jeopardy-game/internal/db"
	"jeopardy-game/internal/handlers"
	"jeopardy-game/internal/repository"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// configMgr := config.NewGameConfigManager("configs/questions.yaml")
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("staring jeopardy-game")
	log.Info("config", cfg)

	database, err := db.NewDB(cfg.StoragePath)
	if err != nil {
		log.Error(
			"failed to init storage", slog.Attr{
				Key:   "error_message",
				Value: slog.StringValue(err.Error()),
			},
		)
		os.Exit(1)
	}
	repo := repository.NewRepository(database)

	templateHandler, err := handlers.NewTemplateHandler(repo)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/", templateHandler.Index)
	router.Post("/player", templateHandler.CreatePlayerHtmx)
	router.Get("/game", templateHandler.GameHtmx)

	// router.Route("/api", func(r chi.Router) {
	// 	r.Get("/players", playerHandler.GetPlayers)
	// 	r.Post("/player", playerHandler.CreatePlayer)
	// })
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.Attr{
			Key:   "error_message",
			Value: slog.StringValue(err.Error()),
		})
	}
	log.Info("server stopped")
}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return log
}
