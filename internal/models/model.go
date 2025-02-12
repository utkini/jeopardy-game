package models

type Player struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Score int    `json:"score" db:"score"`
}

type Question struct {
	ID       int    `json:"id" db:"id"`
	Category string `json:"category" db:"category"`
	Question string `json:"question" db:"question"`
	Answer   string `json:"answer" db:"answer"`
	Points   int    `json:"points" db:"points"`
}
