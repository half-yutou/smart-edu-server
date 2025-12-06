package class

import (
	"smarteduhub/internal/model"
)

type Repository interface {
	Create(class *model.Class) error
	Save(class *model.Class) error
	Delete(teacherID int64, classID int64) error
	ListForTeacher(teacherID int64) ([]*model.Class, error)

	GetByCode(code string) (*model.Class, error)
	GetByID(id int64) (*model.Class, error)

	Quit(studentID int64, classID int64) error
	JoinByCode(studentID int64, code string) error
	ListForStudent(studentID int64) ([]*model.Class, error)

	IsMember(classID int64, userID int64) (bool, error)

	// 资源关联
	AddResource(classID, resourceID int64) error
	RemoveResource(classID, resourceID int64) error
	ListResources(classID int64, page, pageSize int) ([]*model.Resource, int64, error)
}
