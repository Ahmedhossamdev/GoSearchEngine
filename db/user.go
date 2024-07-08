package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
	IsAdmin  bool   `gorm:"default:false" json:"isAdmin"`

	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
	user := User{
		Email:    "admin@gmail.com",
		Password: "admin",
		IsAdmin:  true,
	}

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}

	// Store hashed password in user object
	user.Password = string(hashedPassword)

	// Create user in database
	if err := DBConn.Create(&user).Error; err != nil {
		return errors.New("error creating user")
	}
	return nil
}



func (u *User) LoginAsAdmin(email string, password string) (*User, error) {
	var user User
	if err := DBConn.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
