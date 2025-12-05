package class

import (
	"errors"

	"smarteduhub/internal/model"
	"smarteduhub/internal/pkg/database"

	"gorm.io/gorm"
)

type repositoryImpl struct{}

var _ Repository = (*repositoryImpl)(nil)

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (r *repositoryImpl) Create(class *model.Class) error {
	return database.DB.Create(class).Error
}

func (r *repositoryImpl) Save(class *model.Class) error {
	return database.DB.Save(class).Error
}

func (r *repositoryImpl) GetByCode(code string) (*model.Class, error) {
	var class model.Class
	err := database.DB.
		Preload("Teacher").
		Preload("Subject").
		Preload("Grade").
		Where("code = ?", code).
		First(&class).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

func (r *repositoryImpl) GetByID(id int64) (*model.Class, error) {
	var class model.Class
	err := database.DB.
		Preload("Teacher").
		Preload("Subject").
		Preload("Grade").
		First(&class, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

func (r *repositoryImpl) ListForTeacher(teacherID int64) ([]*model.Class, error) {
	var classes []*model.Class
	err := database.DB.
		Preload("Teacher").
		Preload("Subject").
		Preload("Grade").
		Where("teacher_id = ?", teacherID).
		Order("created_at desc").
		Find(&classes).Error
	return classes, err
}

func (r *repositoryImpl) Delete(teacherID int64, classID int64) error {
	return database.DB.Where("id = ? AND teacher_id = ?", classID, teacherID).Delete(&model.Class{}).Error
}

func (r *repositoryImpl) ListForStudent(studentID int64) ([]*model.Class, error) {
	var classes []*model.Class
	err := database.DB.
		Preload("Teacher").
		Preload("Subject").
		Preload("Grade").
		Joins("JOIN class_members ON classes.id = class_members.class_id").
		Where("class_members.student_id = ?", studentID).
		Order("classes.created_at desc").
		Find(&classes).Error
	return classes, err
}

// Quit 学生退出班级
func (r *repositoryImpl) Quit(studentID int64, classID int64) error {
	class, err := r.GetByID(classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 检查学生是否已加入该班级
	var count int64
	database.DB.
		Model(&model.ClassMember{}).
		Where("class_id = ? AND student_id = ?", class.ID, studentID).
		Count(&count)
	if count == 0 {
		return errors.New("学生未加入该班级")
	}

	// 删除班级成员关系
	return database.DB.
		Where("class_id = ? AND student_id = ?", classID, studentID).
		Delete(&model.ClassMember{}).Error
}

func (r *repositoryImpl) JoinByCode(studentID int64, code string) error {
	// 1. 检查班级是否存在
	class, err := r.GetByCode(code)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 2. 检查学生是否已加入该班级
	var count int64
	database.DB.
		Model(&model.ClassMember{}).
		Where("class_id = ? AND student_id = ?", class.ID, studentID).
		Count(&count)
	if count > 0 {
		return errors.New("学生已加入该班级")
	}

	// 3. 加入班级
	return database.DB.Create(&model.ClassMember{
		ClassID:   class.ID,
		StudentID: studentID,
	}).Error
}

func (r *repositoryImpl) IsMember(classID int64, userID int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND student_id = ?", classID, userID).
		Count(&count).Error
	return count > 0, err
}
