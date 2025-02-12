package repository

import (
	"github.com/jmoiron/sqlx"
)

type MainRepository struct {
	Player   *PlayerRepository
	Question *QuestionRepository
}

func NewRepository(db *sqlx.DB) *MainRepository {
	playerRepo := NewPlayerRepo(db)
	questionRepo := NewQuestionRepo(db)
	return &MainRepository{
		Player:   playerRepo,
		Question: questionRepo,
	}
}
