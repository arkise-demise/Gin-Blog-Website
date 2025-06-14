package models

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id"` // Foreign key to User
	BlogID    uint      `json:"blog_id"` // Foreign key to Blog
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user"`
}
