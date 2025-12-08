package class

import (
	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
)

type Service interface {
	Create(teacherID int64, req *request.CreateClassRequest) (*response.ClassInfo, error)
	DeleteByID(teacherID int64, classID int64) error
	UpdateByID(teacherID int64, req *request.UpdateClassRequest) error
	ListForTeacher(teacherID int64) ([]*response.ClassInfo, error)

	GetByCode(code string) (*response.ClassInfo, error)
	GetByID(operatorID int64, classID int64) (*response.ClassInfo, error)

	Quit(studentID int64, classID int64) error
	JoinByCode(studentID int64, code string) error
	ListForStudent(studentID int64) ([]*response.ClassInfo, error)
	ListMembers(teacherID int64, classID int64) ([]*model.User, error)

	// 资源相关
	AddResource(teacherID int64, req *request.AddResourceToClassRequest) error
	RemoveResource(teacherID int64, req *request.RemoveResourceFromClassRequest) error
	ListResources(operatorID int64, req *request.ListClassResourcesRequest) (*response.PageResult, error)
}
