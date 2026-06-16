package config

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"rest-api-disbursement-system/internal/models"
)

type seedUser struct {
	Username string
	Password string
	Role     string
}

func SeedUsers(db *gorm.DB) error {
	users := []seedUser{
		{
			Username: "admin",
			Password: "admin123",
			Role:     "admin",
		},
		{
			Username: "operator",
			Password: "operator123",
			Role:     "operator",
		},
	}

	for _, seed := range users {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		var user models.User
		err = db.Where("username = ?", seed.Username).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.User{
				Username:     seed.Username,
				Role:         seed.Role,
				PasswordHash: string(passwordHash),
			}

			if err := db.Create(&user).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}

		user.Username = seed.Username
		user.Role = seed.Role
		user.PasswordHash = string(passwordHash)

		if err := db.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
