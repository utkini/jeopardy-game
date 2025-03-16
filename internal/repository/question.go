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
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) GetAll() ([]models.Question, error) {
	var questions []models.Question
	query := `
		SELECT q.*, c.name as category 
		FROM questions q 
		JOIN categories c ON q.category_id = c.id 
		ORDER BY c.name, q.points`
	err := r.db.Select(&questions, query)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *QuestionRepository) GetByID(id int) (*models.Question, error) {
	var question models.Question
	query := `
		SELECT q.*, c.name as category 
		FROM questions q 
		JOIN categories c ON q.category_id = c.id 
		WHERE q.id = ?`
	err := r.db.Get(&question, query, id)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepository) Create(q *models.Question) (int64, error) {
	query := `
		INSERT INTO questions (
			category_id, question, answer, points, 
			media_type, media_url, is_answered, was_correct
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		q.CategoryID, q.Question, q.Answer, q.Points,
		q.MediaType, q.MediaURL, q.IsAnswered, q.WasCorrect)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *QuestionRepository) MarkAsAnswered(id int, wasCorrect bool) error {
	_, err := r.db.Exec(
		"UPDATE questions SET is_answered = 1, was_correct = ? WHERE id = ?",
		wasCorrect, id)
	return err
}

func (r *QuestionRepository) GetByCategoryID(categoryID int) ([]models.Question, error) {
	var questions []models.Question
	query := `
		SELECT q.*, c.name as category 
		FROM questions q 
		JOIN categories c ON q.category_id = c.id 
		WHERE q.category_id = ? 
		ORDER BY q.points`
	err := r.db.Select(&questions, query, categoryID)
	if err != nil {
		return nil, err
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

func (r *QuestionRepository) GetByCategory(category string) ([]models.Question, error) {
	const op = "repository.question.GetByCategory"
	var questions []models.Question
	err := r.db.Select(&questions, "SELECT * FROM questions where category = ?", category)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return questions, nil
}
