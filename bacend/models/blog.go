// package models

//	type Blog struct {
//	    Id          uint   `json:"id,omitempty"`
//	    Title       string `json:"title,omitempty"`
//	    Description string `json:"description,omitempty"`
//	    Image       string `json:"image,omitempty"`
//	    UserID      uint   `json:"user_id,omitempty"`
//	    User        User   `json:"user" gorm:"foreignKey:UserID"`
//	}
package models

import (
	"time"

	"gorm.io/gorm" // Make sure gorm.io/gorm is imported
)

type Blog struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user"` // Belongs To User
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// NEW: Field for approval status
	IsApproved bool `json:"is_approved" gorm:"default:false"`
}

// Ensure Blog also has CreatedAt and UpdatedAt, and soft delete if desired
func (blog *Blog) BeforeCreate(tx *gorm.DB) (err error) {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	return
}

func (blog *Blog) BeforeUpdate(tx *gorm.DB) (err error) {
	blog.UpdatedAt = time.Now()
	return
}
