package repository

import (
	"fmt"
	"jeopardy-game/internal/models"

	"github.com/jmoiron/sqlx"
)

type PlayerRepository struct {
	db *sqlx.DB
}

func NewPlayerRepo(db *sqlx.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (r *PlayerRepository) GetAll() ([]models.Player, error) {
	const op = "repository.player.GetPlayers"
	var players []models.Player
	err := r.db.Select(&players, "SELECT * FROM players ORDER BY score DESC")
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return players, nil
}

func (r *PlayerRepository) GetByID(id int) (*models.Player, error) {
	var player models.Player
	err := r.db.Get(&player, "SELECT * FROM players WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) Create(name string) (int64, error) {
	const op = "repository.player.Create"
	result, err := r.db.Exec("INSERT INTO players (name, score) VALUES (?, 0)", name)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	return result.LastInsertId()
}

func (r *PlayerRepository) UpdateScore(id int, score int) error {
	_, err := r.db.Exec("UPDATE players SET score = ? WHERE id = ?", score, id)
	return err
}
