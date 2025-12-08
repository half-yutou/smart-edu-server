package resource

import (
	"errors"

	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	resourceRepo "smarteduhub/internal/repository/resource"
)

type serviceImpl struct {
	repo resourceRepo.Repository
}

var _ Service = (*serviceImpl)(nil)

func NewService() Service {
	return &serviceImpl{
		repo: resourceRepo.NewRepository(),
	}
}

func (s *serviceImpl) toResourceInfo(r *model.Resource) *response.ResourceInfo {
	info := &response.ResourceInfo{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		ResType:     r.ResType,
		FileURL:     r.FileURL,
		SubjectID:   r.SubjectID,
		GradeID:     r.GradeID,
		UploaderID:  r.UploaderID,
		Duration:    r.Duration,
		CreatedAt:   r.CreatedAt,
	}
	if r.Subject != nil {
		info.SubjectName = r.Subject.Name
	}
	if r.Grade != nil {
		info.GradeName = r.Grade.Name
	}
	if r.Uploader != nil {
		info.UploaderName = r.Uploader.Nickname
	}
	return info
}

func (s *serviceImpl) Create(uploaderID int64, req *request.CreateResourceRequest) (*response.ResourceInfo, error) {
	res := &model.Resource{
		Title:       req.Title,
		Description: req.Description,
		ResType:     req.ResType,
		FileURL:     req.FileURL,
		SubjectID:   req.SubjectID,
		GradeID:     req.GradeID,
		UploaderID:  uploaderID,
		Duration:    req.Duration,
	}

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return s.toResourceInfo(res), nil
}

func (s *serviceImpl) Update(operatorID int64, req *request.UpdateResourceRequest) error {
	res, err := s.repo.GetByID(req.ResourceID)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("resource not found")
	}

	// 权限检查：只能修改自己的资源
	if res.UploaderID != operatorID {
		return errors.New("permission denied")
	}

	// 部分更新
	if req.Title != nil {
		res.Title = *req.Title
	}
	if req.Description != nil {
		res.Description = *req.Description
	}
	if req.SubjectID != nil {
		res.SubjectID = *req.SubjectID
	}
	if req.GradeID != nil {
		res.GradeID = *req.GradeID
	}
	if req.ResType != nil {
		res.ResType = *req.ResType
	}
	if req.FileURL != nil {
		res.FileURL = *req.FileURL
	}

	return s.repo.Update(res)
}

func (s *serviceImpl) Delete(operatorID int64, resourceID int64) error {
	res, err := s.repo.GetByID(resourceID)
	if err != nil {
		return err
	}
	if res == nil {
		return nil // 资源不存在，视为删除成功
	}

	// 权限检查
	if res.UploaderID != operatorID {
		return errors.New("permission denied")
	}

	return s.repo.Delete(resourceID)
}

func (s *serviceImpl) GetByID(resourceID int64) (*response.ResourceInfo, error) {
	res, err := s.repo.GetByID(resourceID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("resource not found")
	}
	return s.toResourceInfo(res), nil
}

func (s *serviceImpl) List(req *request.ListResourcesRequest) (*response.PageResult, error) {
	list, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	var infos []*response.ResourceInfo
	for _, r := range list {
		infos = append(infos, s.toResourceInfo(r))
	}

	return &response.PageResult{
		List:  infos,
		Total: total,
	}, nil
}

func (s *serviceImpl) ListMyResources(uploaderID int64) ([]*response.ResourceInfo, error) {
	list, err := s.repo.ListByUploader(uploaderID)
	if err != nil {
		return nil, err
	}

	var infos []*response.ResourceInfo
	for _, r := range list {
		infos = append(infos, s.toResourceInfo(r))
	}
	return infos, nil
}
