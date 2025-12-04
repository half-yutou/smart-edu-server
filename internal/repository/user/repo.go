package user

import "smarteduhub/internal/model"

type Repository interface {
	CreateUser(user *model.User) error
	GetUserByID(id int64) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
}
