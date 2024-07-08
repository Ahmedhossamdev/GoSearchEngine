package services

import (
	"Ahmedhossamdev/search-engine/db"
	"Ahmedhossamdev/search-engine/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin() error {
	user := models.User{
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
	if err := db.GetDB().Create(&user).Error; err != nil {
		return errors.New("error creating user")
	}
	return nil
}

func LoginAsAdmin(email string, password string) (*models.User, error) {
	var user models.User
	if err := db.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
