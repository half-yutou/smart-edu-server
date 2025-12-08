package danmaku

import (
	"smarteduhub/internal/model"
	"smarteduhub/internal/pkg/database"
)

type Repository interface {
	Create(d *model.Danmaku) error
	GetByResourceID(resourceID int64) ([]*model.Danmaku, error)
}

type repositoryImpl struct{}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (r *repositoryImpl) Create(d *model.Danmaku) error {
	return database.DB.Create(d).Error
}

func (r *repositoryImpl) GetByResourceID(resourceID int64) ([]*model.Danmaku, error) {
	var list []*model.Danmaku
	err := database.DB.Where("resource_id = ?", resourceID).Order("time_point asc").Find(&list).Error
	return list, err
}
