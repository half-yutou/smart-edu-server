package user

import (
	"errors"

	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/pkg/utils"
	userRepo "smarteduhub/internal/repository/user"

	"github.com/click33/sa-token-go/stputil"
)

type serviceImpl struct {
	repo userRepo.Repository
}

var _ Service = (*serviceImpl)(nil)

func NewService() Service {
	return &serviceImpl{
		repo: userRepo.NewRepository(),
	}
}

func (s *serviceImpl) Register(req *request.RegisterRequest) error {
	existUser, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return err
	}
	if existUser != nil {
		return errors.New("username already exists")
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPwd,
		Role:     req.Role,
		Nickname: req.Nickname,
	}

	return s.repo.CreateUser(user)
}

func (s *serviceImpl) Login(req *request.LoginRequest) (*request.LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	// 用户不存在
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// 密码错误
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Sa-Token Login
	token, err := stputil.Login(user.ID)
	if err != nil {
		return nil, err
	}

	// Store role in session
	session, err := stputil.GetSession(user.ID)
	if err == nil {
		_ = session.Set("role", user.Role)
	}

	return &request.LoginResponse{
		ID:    user.ID,
		Token: token,
		Role:  user.Role,
		Name:  user.Nickname,
	}, nil
}

func (s *serviceImpl) Logout(token string) error {
	return stputil.LogoutByToken(token)
}

func (s *serviceImpl) UpdateProfile(userID int64, req *request.UpdateProfileRequest) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	return s.repo.UpdateUser(user)
}
