// package models

// import "time"

// type Comment struct {
// 	ID        uint      `json:"id" gorm:"primarykey"`
// 	Content   string    `json:"content" gorm:"type:text;not null"`
// 	UserID    uint      `json:"user_id"` // Foreign key to User
// 	BlogID    uint      `json:"blog_id"` // Foreign key to Blog
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`

//		User User `json:"user"`
//	}
package models

import (
	"time"

	"gorm.io/gorm" // Make sure gorm.io/gorm is imported
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"` // Belongs to User
	BlogID    uint      `json:"blog_id"`
	Blog      Blog      `json:"blog"` // Belongs to Blog
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// NEW: Field for approval status
	IsApproved bool `json:"is_approved" gorm:"default:false"`
}

// Ensure Comment also has CreatedAt and UpdatedAt, and soft delete if desired
func (comment *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return
}

func (comment *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	comment.UpdatedAt = time.Now()
	return
}
