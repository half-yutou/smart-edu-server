package user

import (
	"smarteduhub/internal/model/dto/request"
)

type Service interface {
	Register(req *request.RegisterRequest) error
	Login(req *request.LoginRequest) (*request.LoginResponse, error)
	Logout(token string) error
	UpdateProfile(userID int64, req *request.UpdateProfileRequest) error
}
