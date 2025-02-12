package handlers

import (
	"encoding/json"
	"jeopardy-game/internal/models"
	"jeopardy-game/internal/repository"
	"net/http"
)

type PlayerHandler struct {
	repo *repository.PlayerRepository
}

func NewPlayerHandler(repo *repository.PlayerRepository) *PlayerHandler {
	return &PlayerHandler{repo: repo}
}

func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.repo.Create(player.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}
