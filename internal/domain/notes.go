package domain

import "time"

type Note struct {
	ID        string
	Title     string
	Content   string
	UserID    string
	CreatedAt time.Time
}
