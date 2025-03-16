package main

import (
	"jeopardy-game/internal/config"
	"jeopardy-game/internal/db"
	"jeopardy-game/internal/handlers"
	"jeopardy-game/internal/models"
	"jeopardy-game/internal/repository"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("starting jeopardy-game")

	// Load game config
	gameConfig, err := config.NewGameConfigManager("configs/questions.yaml")
	if err != nil {
		log.Error("failed to load game config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize database
	database, err := db.NewDB(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	repo := repository.NewRepository(database)

	// Initialize questions from config
	if err := initializeQuestions(repo, gameConfig.Config); err != nil {
		log.Error("failed to initialize questions", slog.String("error", err.Error()))
		os.Exit(1)
	}

	templateHandler, err := handlers.NewTemplateHandler(repo)
	if err != nil {
		log.Error("failed to init template handler", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := mux.NewRouter()

	// Serve static files for media content
	fileServer := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Routes
	router.HandleFunc("/", templateHandler.Index).Methods("GET")
	router.HandleFunc("/player", templateHandler.CreatePlayerHtmx).Methods("POST")
	router.HandleFunc("/questions", templateHandler.QuestionsHtmx).Methods("GET")
	router.HandleFunc("/question/{id}", templateHandler.GetQuestionHtmx).Methods("GET")
	router.HandleFunc("/question/{id}/answer", templateHandler.AnswerQuestionHtmx).Methods("POST")
	router.HandleFunc("/game/reset", templateHandler.ResetGameHtmx).Methods("POST")

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func initializeQuestions(repo *repository.MainRepository, gameConfig config.GameConfig) error {
	// Create categories and questions
	for _, cat := range gameConfig.Categories {
		// Check if category exists
		existingCat, err := repo.Category.GetByName(cat.Name)
		var categoryID int64
		if err == nil {
			// Category exists, use its ID
			categoryID = int64(existingCat.ID)
		} else {
			// Create new category
			categoryID, err = repo.Category.Create(cat.Name)
			if err != nil {
				return err
			}
		}

		// Create questions for this category
		for _, q := range cat.Questions {
			question := &models.Question{
				CategoryID: int(categoryID),
				Question:   q.Question,
				Answer:     q.Answer,
				Points:     q.Points,
				MediaType:  q.MediaType,
				MediaURL:   q.MediaURL,
			}
			_, err := repo.Question.Create(question)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
