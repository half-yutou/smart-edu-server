package class_member

import (
	"smarteduhub/internal/model"
	"smarteduhub/internal/pkg/database"
)

type repositoryImpl struct {
}

// Ensure implementation
var _ Repository = (*repositoryImpl)(nil)

func NewRepository() Repository {
	return &repositoryImpl{}
}

// CountMembers 统计班级成员数量
func (r *repositoryImpl) CountMembers(classID int64) (int, error) {
	var count int64
	err := database.DB.Model(&model.ClassMember{}).
		Where("class_id = ?", classID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
