package models

type Blog struct {
    Id          uint   `json:"id,omitempty"`
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Image       string `json:"image,omitempty"`
    UserID      uint   `json:"user_id,omitempty"` 
    User        User   `json:"user" gorm:"foreignKey:UserID"`
}
