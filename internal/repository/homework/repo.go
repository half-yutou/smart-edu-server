package homework

import "smarteduhub/internal/model"

type Repository interface {
	Create(homework *model.Homework) error
	GetByID(id int64) (*model.Homework, error)
	Delete(id int64) error
	Update(homework *model.Homework) error
	ListByCreator(creatorID int64) ([]*model.Homework, error)
	// Submission 相关
	GetSubmission(homeworkID, studentID int64) (*model.Submission, error)
	GetSubmissionByID(id int64) (*model.Submission, error)
	SaveSubmission(sub *model.Submission) error
	ListSubmissions(homeworkID int64) ([]*model.Submission, error)
	ListByClass(classID int64) ([]*model.Homework, error)
}