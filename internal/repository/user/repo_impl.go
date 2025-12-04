package user

import (
	"errors"

	"smarteduhub/internal/model"
	"smarteduhub/internal/pkg/database"

	"gorm.io/gorm"
)

type repositoryImpl struct{}

// Ensure implementation
var _ Repository = (*repositoryImpl)(nil)

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (r *repositoryImpl) CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func (r *repositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *repositoryImpl) UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func (r *repositoryImpl) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
