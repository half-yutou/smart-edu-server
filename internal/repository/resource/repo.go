package resource

import (
	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
)

type Repository interface {
	Create(resource *model.Resource) error
	Update(resource *model.Resource) error
	Delete(id int64) error
	GetByID(id int64) (*model.Resource, error)

	// List 支持多条件筛选和分页
	List(req *request.ListResourcesRequest) ([]*model.Resource, int64, error)

	// ListByUploader 查询某个用户上传的资源
	ListByUploader(uploaderID int64) ([]*model.Resource, error)
}
