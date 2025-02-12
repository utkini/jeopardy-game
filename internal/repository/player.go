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
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) GetAll() ([]models.Player, error) {
	const op = "repository.player.GetPlayers"
	var players []models.Player
	err := r.db.Select(&players, "SELECT * FROM players")
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return players, nil
}

func (r *PlayerRepository) Create(name string) (int64, error) {
	const op = "repository.player.Create"
	res, err := r.db.Exec(
		"INSERT INTO players(name) VALUES(?)",
		name,
	)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	return id, nil
}
