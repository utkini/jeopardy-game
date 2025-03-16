package handlers

import (
	"html/template"
	"jeopardy-game/internal/models"
	"jeopardy-game/internal/repository"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type TemplateHandler struct {
	repo        *repository.MainRepository
	templates   *template.Template
	currentTurn int   // Index of the current player's turn
	playerOrder []int // Array to maintain the order of player IDs
}

func NewTemplateHandler(
	repo *repository.MainRepository,
) (*TemplateHandler, error) {
	tmpl := template.New("")

	// First parse the layout template as it contains the base structure
	layoutPath := filepath.Join("templates", "layout.html")
	tmpl, err := tmpl.ParseFiles(layoutPath)
	if err != nil {
		return nil, err
	}

	// Then parse all other templates
	patterns := []string{
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "players-list.html"),
		filepath.Join("templates", "questions.html"),
		filepath.Join("templates", "question-modal.html"),
	}

	for _, pattern := range patterns {
		_, err := tmpl.ParseFiles(pattern)
		if err != nil {
			return nil, err
		}
	}

	return &TemplateHandler{
		repo:        repo,
		templates:   tmpl,
		currentTurn: 0,
		playerOrder: make([]int, 0),
	}, nil
}

func (h *TemplateHandler) Index(w http.ResponseWriter, r *http.Request) {
	players, err := h.repo.Player.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Players":    players,
		"CurrentTab": "players",
	}

	h.templates.ExecuteTemplate(w, "layout", data)
}

func (h *TemplateHandler) CreatePlayerHtmx(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Create new player
	playerID, err := h.repo.Player.Create(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add player to the order
	h.playerOrder = append(h.playerOrder, int(playerID))
	log.Printf("Added player %s (ID: %d) to the game. Current order: %v", name, playerID, h.playerOrder)

	players, err := h.repo.Player.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Players":    players,
		"CurrentTab": "players",
	}

	h.templates.ExecuteTemplate(w, "players-list", data)
}

func (h *TemplateHandler) QuestionsHtmx(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.Category.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all questions
	allQuestions, err := h.repo.Question.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a map of questions by category ID
	questionsByCategory := make(map[int][]models.Question)
	for _, q := range allQuestions {
		questionsByCategory[q.CategoryID] = append(questionsByCategory[q.CategoryID], q)
	}

	// Sort questions by points within each category
	for _, questions := range questionsByCategory {
		sort.Slice(questions, func(i, j int) bool {
			return questions[i].Points < questions[j].Points
		})
	}

	// Create category rows with their questions
	type CategoryRow struct {
		Category  models.Category
		Questions []models.Question
	}
	var categoryRows []CategoryRow

	for _, cat := range categories {
		row := CategoryRow{
			Category:  cat,
			Questions: questionsByCategory[cat.ID],
		}
		categoryRows = append(categoryRows, row)
	}

	data := map[string]interface{}{
		"CategoryRows": categoryRows,
		"CurrentTab":   "questions",
	}

	h.templates.ExecuteTemplate(w, "questions", data)
}

func (h *TemplateHandler) GetQuestionHtmx(w http.ResponseWriter, r *http.Request) {
	// Extract question ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Printf("Invalid URL path: %s", r.URL.Path)
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]
	questionID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error parsing question ID: %v", err)
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Get question from repository
	question, err := h.repo.Question.GetByID(questionID)
	if err != nil {
		log.Printf("Error getting question %d: %v", questionID, err)
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	// Get all players
	players, err := h.repo.Player.GetAll()
	if err != nil {
		log.Printf("Error getting players: %v", err)
		http.Error(w, "Error getting players", http.StatusInternalServerError)
		return
	}

	// Get current player based on turn
	var currentPlayer models.Player
	if len(h.playerOrder) > 0 {
		currentPlayerID := h.playerOrder[h.currentTurn%len(h.playerOrder)]
		for _, p := range players {
			if p.ID == currentPlayerID {
				currentPlayer = p
				break
			}
		}
	}

	// Log the question being displayed
	log.Printf("Displaying question %d for player %s (Turn %d/%d)",
		questionID, currentPlayer.Name, h.currentTurn%len(h.playerOrder)+1, len(h.playerOrder))

	// Execute template with question and current player
	data := map[string]interface{}{
		"Question":      question,
		"CurrentPlayer": currentPlayer,
	}
	if err := h.templates.ExecuteTemplate(w, "question-modal", data); err != nil {
		log.Printf("Error executing question modal template: %v", err)
		http.Error(w, "Error displaying question", http.StatusInternalServerError)
		return
	}
}

func (h *TemplateHandler) AnswerQuestionHtmx(w http.ResponseWriter, r *http.Request) {
	// Extract question ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Printf("Invalid URL path: %s", r.URL.Path)
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-2] // ID is second to last part in /question/{id}/answer
	questionID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error parsing question ID: %v", err)
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Get question from repository
	question, err := h.repo.Question.GetByID(questionID)
	if err != nil {
		log.Printf("Error getting question %d: %v", questionID, err)
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	// Get all players
	players, err := h.repo.Player.GetAll()
	if err != nil {
		log.Printf("Error getting players: %v", err)
		http.Error(w, "Error getting players", http.StatusInternalServerError)
		return
	}

	if len(h.playerOrder) == 0 {
		log.Printf("No players available to answer question")
		http.Error(w, "No players available", http.StatusBadRequest)
		return
	}

	// Get current player
	currentPlayerID := h.playerOrder[h.currentTurn%len(h.playerOrder)]
	var currentPlayer models.Player
	for _, p := range players {
		if p.ID == currentPlayerID {
			currentPlayer = p
			break
		}
	}

	// Parse form for correct/incorrect answer
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error processing answer", http.StatusBadRequest)
		return
	}

	wasCorrect := r.FormValue("correct") == "true"
	points := question.Points
	if !wasCorrect {
		points = -points // Deduct points for incorrect answer
	}

	// Update player's score
	log.Printf("Updating score for player %s (ID: %d): %d points (correct: %v)",
		currentPlayer.Name, currentPlayer.ID, points, wasCorrect)
	if err := h.repo.Player.UpdateScore(currentPlayer.ID, points); err != nil {
		log.Printf("Error updating player score: %v", err)
		http.Error(w, "Error updating score", http.StatusInternalServerError)
		return
	}

	// Mark question as answered
	if err := h.repo.Question.MarkAsAnswered(questionID, wasCorrect); err != nil {
		log.Printf("Error marking question as answered: %v", err)
		http.Error(w, "Error updating question status", http.StatusInternalServerError)
		return
	}

	// Move to next player's turn
	h.currentTurn = (h.currentTurn + 1) % len(players)
	log.Printf("Moving to next player's turn: %s", players[h.currentTurn].Name)

	// Get updated categories and questions
	categories, err := h.repo.Category.GetAll()
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		http.Error(w, "Error updating game board", http.StatusInternalServerError)
		return
	}

	// Get all questions
	allQuestions, err := h.repo.Question.GetAll()
	if err != nil {
		log.Printf("Error getting questions: %v", err)
		http.Error(w, "Error updating game board", http.StatusInternalServerError)
		return
	}

	// Create a map of questions by category ID
	questionsByCategory := make(map[int][]models.Question)
	for _, q := range allQuestions {
		questionsByCategory[q.CategoryID] = append(questionsByCategory[q.CategoryID], q)
	}

	// Sort questions by points within each category
	for _, questions := range questionsByCategory {
		sort.Slice(questions, func(i, j int) bool {
			return questions[i].Points < questions[j].Points
		})
	}

	// Create category rows with their questions
	type CategoryRow struct {
		Category  models.Category
		Questions []models.Question
	}
	var categoryRows []CategoryRow

	for _, cat := range categories {
		row := CategoryRow{
			Category:  cat,
			Questions: questionsByCategory[cat.ID],
		}
		categoryRows = append(categoryRows, row)
	}

	data := map[string]interface{}{
		"CategoryRows": categoryRows,
		"CurrentTab":   "questions",
	}

	// Execute the questions template
	if err := h.templates.ExecuteTemplate(w, "questions", data); err != nil {
		log.Printf("Error executing questions template: %v", err)
		http.Error(w, "Error updating game board", http.StatusInternalServerError)
		return
	}
}

func (h *TemplateHandler) ResetGameHtmx(w http.ResponseWriter, r *http.Request) {
	// Reset all player scores
	players, err := h.repo.Player.GetAll()
	if err != nil {
		log.Printf("Error getting players: %v", err)
		http.Error(w, "Error resetting game", http.StatusInternalServerError)
		return
	}

	for _, player := range players {
		if err := h.repo.Player.UpdateScore(player.ID, 0); err != nil {
			log.Printf("Error resetting score for player %s: %v", player.Name, err)
			http.Error(w, "Error resetting scores", http.StatusInternalServerError)
			return
		}
	}

	// Reset all questions to unanswered state
	questions, err := h.repo.Question.GetAll()
	if err != nil {
		log.Printf("Error getting questions: %v", err)
		http.Error(w, "Error resetting questions", http.StatusInternalServerError)
		return
	}

	for _, q := range questions {
		if q.IsAnswered {
			if err := h.repo.Question.MarkAsAnswered(q.ID, false); err != nil {
				log.Printf("Error resetting question %d: %v", q.ID, err)
				http.Error(w, "Error resetting questions", http.StatusInternalServerError)
				return
			}
		}
	}

	// Reset turn counter
	h.currentTurn = 0
	log.Printf("Game reset complete. Scores and questions reset, turn counter reset to 0")

	// Return updated players list
	players, err = h.repo.Player.GetAll()
	if err != nil {
		log.Printf("Error getting updated players: %v", err)
		http.Error(w, "Error getting players", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Players":    players,
		"CurrentTab": "players",
	}

	h.templates.ExecuteTemplate(w, "players-list", data)
}
