package user

import (
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
)

type Service interface {
	Register(req *request.RegisterRequest) error
	Login(req *request.LoginRequest) (*response.LoginResponse, error)
	Logout(token string) error
	UpdateProfile(userID int64, req *request.UpdateProfileRequest) error
}
