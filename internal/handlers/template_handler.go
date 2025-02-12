package handlers

import (
	"html/template"
	"jeopardy-game/internal/models"
	"jeopardy-game/internal/repository"
	"net/http"
	"path/filepath"
)

type TemplateHandler struct {
	repo      *repository.MainRepository
	templates *template.Template
}

func NewTemplateHandler(
	repo *repository.MainRepository,
) (*TemplateHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.gohtml"))
	if err != nil {
		return nil, err
	}

	return &TemplateHandler{
		repo:      repo,
		templates: tmpl,
	}, nil
}

func (h *TemplateHandler) Index(w http.ResponseWriter, r *http.Request) {
	players, err := h.repo.Player.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Players": players,
	}

	h.templates.ExecuteTemplate(w, "layout", data)
}

func (h *TemplateHandler) CreatePlayerHtmx(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player := models.Player{
		Name: r.FormValue("name"),
	}

	_, err := h.repo.Player.Create(player.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Players": []models.Player{player},
	}

	h.templates.ExecuteTemplate(w, "players-list", data)
}

func (h *TemplateHandler) GameHtmx(w http.ResponseWriter, r *http.Request) {
	questions, err := h.repo.Question.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Questions": questions,
	}

	h.templates.ExecuteTemplate(w, "questions", data)
}
