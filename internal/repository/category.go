package repository

import (
	"jeopardy-game/internal/models"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Select(&categories, "SELECT * FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Create(name string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CategoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Get(&category, "SELECT * FROM categories WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
