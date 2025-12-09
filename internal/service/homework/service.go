package homework

import (
	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
)

type Service interface {
	Create(teacherID int64, req *request.CreateHomeworkRequest) error
	Delete(operatorID int64, req *request.DeleteHomeworkRequest) error
	Update(operatorID int64, req *request.UpdateHomeworkRequest) error
	GetByID(id int64) (*model.Homework, error)
	ListByCreator(creatorID int64) ([]*model.Homework, error)
	ListByClass(studentID, classID int64) ([]*response.HomeworkListItem, error)
	Submit(studentID int64, req *request.SubmitHomeworkRequest) error
	GetSubmission(homeworkID, studentID int64) (*model.Submission, error)

	ListSubmissions(teacherID, homeworkID int64) ([]*model.Submission, error)
	GradeSubmission(teacherID int64, req *request.ManualGradeRequest) error
}
