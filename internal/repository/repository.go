package repository

import (
	"jeopardy-game/internal/models"

	"github.com/jmoiron/sqlx"
)

type MainRepository struct {
	Player   IPlayerRepository
	Question IQuestionRepository
	Category ICategoryRepository
}

func NewRepository(db *sqlx.DB) *MainRepository {
	playerRepo := NewPlayerRepo(db)
	questionRepo := NewQuestionRepo(db)
	categoryRepo := NewCategoryRepo(db)
	return &MainRepository{
		Player:   playerRepo,
		Question: questionRepo,
		Category: categoryRepo,
	}
}

type IPlayerRepository interface {
	GetAll() ([]models.Player, error)
	GetByID(id int) (*models.Player, error)
	Create(name string) (int64, error)
	UpdateScore(id int, score int) error
}

type ICategoryRepository interface {
	GetAll() ([]models.Category, error)
	Create(name string) (int64, error)
	GetByName(name string) (*models.Category, error)
}

type IQuestionRepository interface {
	GetAll() ([]models.Question, error)
	GetByID(id int) (*models.Question, error)
	Create(q *models.Question) (int64, error)
	MarkAsAnswered(id int, wasCorrect bool) error
	GetByCategoryID(categoryID int) ([]models.Question, error)
}
