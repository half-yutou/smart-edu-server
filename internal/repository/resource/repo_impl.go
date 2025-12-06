package resource

import (
	"errors"

	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/pkg/database"

	"gorm.io/gorm"
)

type repositoryImpl struct{}

var _ Repository = (*repositoryImpl)(nil)

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (r *repositoryImpl) Create(resource *model.Resource) error {
	return database.DB.Create(resource).Error
}

func (r *repositoryImpl) Update(resource *model.Resource) error {
	return database.DB.Save(resource).Error
}

func (r *repositoryImpl) Delete(id int64) error {
	return database.DB.Delete(&model.Resource{}, id).Error
}

func (r *repositoryImpl) GetByID(id int64) (*model.Resource, error) {
	var res model.Resource
	err := database.DB.
		Preload("Subject").
		Preload("Grade").
		Preload("Uploader").
		First(&res, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *repositoryImpl) List(req *request.ListResourcesRequest) ([]*model.Resource, int64, error) {
	var resources []*model.Resource
	var count int64

	db := database.DB.Model(&model.Resource{})

	// 动态构建查询条件
	if req.SubjectID != 0 {
		db = db.Where("subject_id = ?", req.SubjectID)
	}
	if req.GradeID != 0 {
		db = db.Where("grade_id = ?", req.GradeID)
	}
	if req.ResType != "" {
		db = db.Where("res_type = ?", req.ResType)
	}
	if req.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 先查询总数
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (req.Page - 1) * req.PageSize
	err := db.
		Preload("Subject").
		Preload("Grade").
		Preload("Uploader").
		Order("created_at desc").
		Offset(offset).
		Limit(req.PageSize).
		Find(&resources).Error

	return resources, count, err
}

func (r *repositoryImpl) ListByUploader(uploaderID int64) ([]*model.Resource, error) {
	var resources []*model.Resource
	err := database.DB.
		Preload("Subject").
		Preload("Grade").
		Preload("Uploader").
		Where("uploader_id = ?", uploaderID).
		Order("created_at desc").
		Find(&resources).Error
	return resources, err
}
