package homework

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

func (r *repositoryImpl) Create(homework *model.Homework) error {
	return database.DB.Create(homework).Error
}

func (r *repositoryImpl) GetByID(id int64) (*model.Homework, error) {
	var hw model.Homework
	err := database.DB.Preload("Questions").First(&hw, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hw, nil
}

func (r *repositoryImpl) Delete(id int64) error {
	// Select("Questions") 会触发级联删除，前提是 questions 表没有设外键约束或者外键允许级联
	// 为了保险，我们也可以显式删除
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("homework_id = ?", id).Delete(&model.Question{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Homework{}, id).Error
	})
}

func (r *repositoryImpl) Update(homework *model.Homework) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新基础信息
		if err := tx.Model(homework).Updates(map[string]interface{}{
			"title":    homework.Title,
			"deadline": homework.Deadline,
		}).Error; err != nil {
			return err
		}

		// 2. 替换题目
		// Replace 会找出差异：删除多余的，添加新增的，更新已有的
		return tx.Model(homework).Association("Questions").Replace(homework.Questions)
	})
}

func (r *repositoryImpl) ListByCreator(creatorID int64) ([]*model.Homework, error) {
	var list []*model.Homework
	err := database.DB.
		Preload("Class"). // 关联班级信息
		Where("creator_id = ?", creatorID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *repositoryImpl) ListByClass(classID int64) ([]*model.Homework, error) {
	var list []*model.Homework
	err := database.DB.
		Preload("Questions").
		Where("class_id = ?", classID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *repositoryImpl) GetSubmission(homeworkID, studentID int64) (*model.Submission, error) {
	var sub model.Submission
	err := database.DB.Preload("Details").
		Where("homework_id = ? AND student_id = ?", homeworkID, studentID).
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *repositoryImpl) GetSubmissionByID(id int64) (*model.Submission, error) {
	var sub model.Submission
	err := database.DB.Preload("Details").First(&sub, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *repositoryImpl) ListSubmissions(homeworkID int64) ([]*model.Submission, error) {
	var list []*model.Submission
	err := database.DB.
		Preload("Student").
		Preload("Details").
		Where("homework_id = ?", homeworkID).
		Order("submitted_at desc").
		Find(&list).Error
	return list, err
}

func (r *repositoryImpl) SaveSubmission(sub *model.Submission) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if sub.ID > 0 {
			// Update: 更新主表 + 手动替换详情
			// 注意：Updates 默认只更新非零值字段。如果想更新所有字段，用 Select("*") 或者 map
			if err := tx.Model(sub).Updates(map[string]interface{}{
				"status":       sub.Status,
				"total_score":  sub.TotalScore,
				"submitted_at": sub.SubmittedAt,
				"feedback":     sub.Feedback,
				"ai_feedback":  sub.AIFeedback,
			}).Error; err != nil {
				return err
			}

			// 手动删除旧详情
			if err := tx.Where("submission_id = ?", sub.ID).Delete(&model.SubmissionDetail{}).Error; err != nil {
				return err
			}

			// 手动插入新详情
			if len(sub.Details) > 0 {
				for i := range sub.Details {
					sub.Details[i].SubmissionID = sub.ID
				}
				if err := tx.Create(&sub.Details).Error; err != nil {
					return err
				}
			}
			return nil
		}
		// Create
		return tx.Create(sub).Error
	})
}
