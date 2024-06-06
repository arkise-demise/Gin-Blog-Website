package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id 			uint       `json:"id,omitempty"`
	FirstName 	string     `json:"first_name,omitempty"`
	LastName 	string     `json:"last_name,omitempty"`
	Email 		string     `json:"email,omitempty"`
	Password 	[]byte     `json:"password,omitempty"`
	Phone 		string     `json:"phone,omitempty"`

}

func (user *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = hashedPassword
    return nil
}

func (user *User)ComparePassword(password string)error{
	return bcrypt.CompareHashAndPassword(user.Password,[]byte(password))
}
