package models

import (
    "time"

    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Id               uint      `json:"id,omitempty"`
    FirstName        string    `json:"first_name,omitempty"`
    LastName         string    `json:"last_name,omitempty"`
    Email            string    `json:"email,omitempty"`
    Password         []byte    `json:"-"` // Changed to "-" to ensure password hash is never serialized
    Phone            string    `json:"phone,omitempty"`
    Role             string    `json:"role" gorm:"type:varchar(50);default:'user'"`
    
    // NEW PROFILE FIELDS
    Bio              string    `json:"bio,omitempty"`               // User's short biography
    ProfilePictureURL string   `json:"profile_picture_url,omitempty"` // URL to user's profile picture
    Location         string    `json:"location,omitempty"`          // User's location
    Website          string    `json:"website,omitempty"`           // User's personal website or blog

    CreatedAt        time.Time `json:"created_at"` // Added for consistency
    UpdatedAt        time.Time `json:"updated_at"` // Added for consistency

    // GORM associations (if you have them)
    Blogs            []Blog    `gorm:"foreignKey:UserID" json:"-"`    // User has many blogs (add if you have Blog model)
    Comments         []Comment `gorm:"foreignKey:UserID" json:"-"` // User has many comments (add if you have Comment model)
}

func (user *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = hashedPassword
    return nil
}

func (user *User) ComparePassword(password string) error {
    return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}