package repository

import (
	"fmt"
	"jeopardy-game/internal/models"

	"github.com/jmoiron/sqlx"
)

type QuestionRepository struct {
	db *sqlx.DB
}

func NewQuestionRepo(db *sqlx.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) GetAll() ([]models.Question, error) {
	const op = "repository.question.GetQuestions"
	var questions []models.Question
	err := r.db.Select(&questions, "SELECT * FROM questions")
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return questions, nil
}

func (r *QuestionRepository) GetByCategory(category string) ([]models.Question, error) {
	const op = "repository.question.GetByCategory"
	var questions []models.Question
	err := r.db.Select(&questions, "SELECT * FROM questions where category = ?", category)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return questions, nil
}

func (r *QuestionRepository) GetAnswer(questionId int) (string, error) {
	const op = "repository.question.GetAnswer"
	var answer string
	err := r.db.Get(&answer, "SELECT answer FROM questions where id = ?", questionId)
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}
	return answer, nil
}

func (r *QuestionRepository) Create(category string, question string, answer string, points int) (int64, error) {
	const op = "repository.question.Create"
	res, err := r.db.Exec(
		"INSERT INTO questions (category, question, answer, points) VALUES (?, ?, ?, ?)",
		category, question, answer, points,
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
