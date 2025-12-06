package resource

import (
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
)

type Service interface {
	Create(uploaderID int64, req *request.CreateResourceRequest) (*response.ResourceInfo, error)
	Update(operatorID int64, req *request.UpdateResourceRequest) error
	Delete(operatorID int64, resourceID int64) error
	GetByID(resourceID int64) (*response.ResourceInfo, error)

	// List 公共资源列表 (分页)
	List(req *request.ListResourcesRequest) (*response.PageResult, error)

	// ListMyResources 查看我的资源
	ListMyResources(uploaderID int64) ([]*response.ResourceInfo, error)
}
