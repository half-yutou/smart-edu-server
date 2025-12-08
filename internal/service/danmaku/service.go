package danmaku

import (
	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	danmakuRepo "smarteduhub/internal/repository/danmaku"
)

type Service interface {
	Send(senderID int64, req *request.SendDanmakuRequest) error
	List(resourceID int64) ([]*model.Danmaku, error)
}

type serviceImpl struct {
	repo danmakuRepo.Repository
}

func NewService() Service {
	return &serviceImpl{
		repo: danmakuRepo.NewRepository(),
	}
}

func (s *serviceImpl) Send(senderID int64, req *request.SendDanmakuRequest) error {
	d := &model.Danmaku{
		ResourceID: req.ResourceID,
		UserID:     senderID,
		Content:    req.Content,
		TimePoint:  req.Time,
		Color:      req.Color,
	}
	if d.Color == "" {
		d.Color = "#ffffff"
	}
	return s.repo.Create(d)
}

func (s *serviceImpl) List(resourceID int64) ([]*model.Danmaku, error) {
	return s.repo.GetByResourceID(resourceID)
}
