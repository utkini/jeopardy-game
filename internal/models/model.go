package models

type Player struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Score int    `json:"score" db:"score"`
}

type Question struct {
	ID         int    `json:"id" db:"id"`
	CategoryID int    `json:"category_id" db:"category_id"`
	Category   string `json:"category" db:"category"`
	Question   string `json:"question" db:"question"`
	Answer     string `json:"answer" db:"answer"`
	Points     int    `json:"points" db:"points"`
	MediaType  string `json:"media_type" db:"media_type"`   // can be "none", "image", "audio", "video"
	MediaURL   string `json:"media_url" db:"media_url"`     // path or URL to media file
	IsAnswered bool   `json:"is_answered" db:"is_answered"` // tracks if question was answered
	WasCorrect bool   `json:"was_correct" db:"was_correct"` // tracks if answer was correct
}

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
