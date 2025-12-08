package class

import (
	"crypto/rand"
	"errors"
	"math/big"

	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	classRepo "smarteduhub/internal/repository/class"
	classMemberRepo "smarteduhub/internal/repository/class_member"
)

const (
	// 定义邀请码字符集
	charset    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length     = 10
	maxRetries = 5
)

type serviceImpl struct {
	classRepo       classRepo.Repository
	classMemberRepo classMemberRepo.Repository
}

var _ Service = (*serviceImpl)(nil)

func NewService() Service {
	return &serviceImpl{
		classRepo:       classRepo.NewRepository(),
		classMemberRepo: classMemberRepo.NewRepository(),
	}
}

// toClassInfo 将 model.Class 转换为 response.ClassInfo
func (s *serviceImpl) toClassInfo(c *model.Class) *response.ClassInfo {
	// 统计班级成员数量
	memberCount, _ := s.classMemberRepo.CountMembers(c.ID)
	info := &response.ClassInfo{
		ID:          c.ID,
		Name:        c.Name,
		Code:        c.Code,
		SubjectID:   c.SubjectID,
		SubjectName: "",
		GradeID:     c.GradeID,
		GradeName:   "",
		TeacherID:   c.TeacherID,
		TeacherName: "",
		MemberCount: memberCount,
		CreatedAt:   c.CreatedAt,
	}

	if c.Subject != nil {
		info.SubjectName = c.Subject.Name
	}
	if c.Grade != nil {
		info.GradeName = c.Grade.Name
	}
	if c.Teacher != nil {
		info.TeacherName = c.Teacher.Nickname
	}
	return info
}

func (s *serviceImpl) Create(teacherID int64, req *request.CreateClassRequest) (*response.ClassInfo, error) {
	var lastErr error
	// 乐观锁重试机制：尝试 maxRetries 次
	for i := 0; i < maxRetries; i++ {
		// 1. 生成随机码 (纯内存操作，无需查库)
		code, err := s.generateRandomCode()
		if err != nil {
			return nil, err
		}

		class := &model.Class{
			Name:      req.Name,
			Code:      code,
			TeacherID: teacherID,
			SubjectID: req.SubjectID,
			GradeID:   req.GradeID,
		}

		// 2. 尝试插入数据库
		// 如果 code 重复，数据库会报 duplicate key value violates unique constraint
		// 利用数据库的原子性来保证唯一性，无需应用层显式加锁
		if err = s.classRepo.Create(class); err == nil {
			// 插入成功后，class.ID 已经被 GORM 填充
			// 我们使用这个 ID 重新查询数据库，利用 Preload 加载关联信息 (Subject, Grade, Teacher)
			createdClass, err := s.classRepo.GetByID(class.ID)
			if err != nil {
				return nil, err
			}
			return s.toClassInfo(createdClass), nil
		} else {
			lastErr = err
		}
	}
	// 超过最大重试次数，返回最后一次的错误
	return nil, lastErr
}

func (s *serviceImpl) UpdateByID(teacherID int64, req *request.UpdateClassRequest) error {
	class, err := s.classRepo.GetByID(req.ClassID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}

	if class.TeacherID != teacherID {
		return errors.New("permission denied: you are not the owner of this class")
	}

	if req.Name != nil {
		class.Name = *req.Name
	}
	if req.GradeID != nil {
		class.GradeID = *req.GradeID
	}
	if req.SubjectID != nil {
		class.SubjectID = *req.SubjectID
	}

	return s.classRepo.Save(class)
}

func (s *serviceImpl) ListForTeacher(teacherID int64) ([]*response.ClassInfo, error) {
	classes, err := s.classRepo.ListForTeacher(teacherID)
	if err != nil {
		return nil, err
	}

	var infos []*response.ClassInfo
	for _, c := range classes {
		infos = append(infos, s.toClassInfo(c))
	}
	return infos, nil
}

func (s *serviceImpl) DeleteByID(teacherID int64, classID int64) error {
	return s.classRepo.Delete(teacherID, classID)
}

func (s *serviceImpl) GetByCode(code string) (*response.ClassInfo, error) {
	class, err := s.classRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, nil
	}
	return s.toClassInfo(class), nil
}

func (s *serviceImpl) GetByID(operatorID int64, classID int64) (*response.ClassInfo, error) {
	// 1. 获取班级信息 (包含 Preload 的关联数据)
	class, err := s.classRepo.GetByID(classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}

	// 2. 校验权限：必须是该班级的老师或者成员
	// 2.1 检查是否是老师
	if class.TeacherID == operatorID {
		return s.toClassInfo(class), nil
	}

	// 2.2 检查是否是成员
	isMember, err := s.classRepo.IsMember(classID, operatorID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return s.toClassInfo(class), nil
	}

	return nil, errors.New("permission denied: you are not a member of this class")
}

// ListForStudent 获取学生加入的班级
func (s *serviceImpl) ListForStudent(studentID int64) ([]*response.ClassInfo, error) {
	classes, err := s.classRepo.ListForStudent(studentID)
	if err != nil {
		return nil, err
	}

	var infos []*response.ClassInfo
	for _, c := range classes {
		infos = append(infos, s.toClassInfo(c))
	}
	return infos, nil
}

// Quit 学生退出班级
func (s *serviceImpl) Quit(studentID int64, classID int64) error {
	return s.classRepo.Quit(studentID, classID)
}

// JoinByCode 学生通过邀请码加入班级
func (s *serviceImpl) JoinByCode(studentID int64, code string) error {
	return s.classRepo.JoinByCode(studentID, code)
}

// AddResource 老师添加资源到班级
func (s *serviceImpl) AddResource(teacherID int64, req *request.AddResourceToClassRequest) error {
	class, err := s.classRepo.GetByID(req.ClassID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}
	if class.TeacherID != teacherID {
		return errors.New("permission denied: you are not the owner of this class")
	}

	return s.classRepo.AddResource(req.ClassID, req.ResourceID)
}

func (s *serviceImpl) RemoveResource(teacherID int64, req *request.RemoveResourceFromClassRequest) error {
	class, err := s.classRepo.GetByID(req.ClassID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}
	if class.TeacherID != teacherID {
		return errors.New("permission denied")
	}

	return s.classRepo.RemoveResource(req.ClassID, req.ResourceID)
}

func (s *serviceImpl) ListResources(operatorID int64, req *request.ListClassResourcesRequest) (*response.PageResult, error) {
	// 权限检查：必须是班级成员（老师或学生）
	class, err := s.classRepo.GetByID(req.ClassID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}

	// 1. 是老师？
	isTeacher := class.TeacherID == operatorID
	// 2. 是学生？
	isMember, _ := s.classRepo.IsMember(req.ClassID, operatorID)

	if !isTeacher && !isMember {
		return nil, errors.New("permission denied: you are not a member of this class")
	}

	// 查询
	list, total, err := s.classRepo.ListResources(req.ClassID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换
	var infos []*response.ResourceInfo
	for _, r := range list {
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
		infos = append(infos, info)
	}

	return &response.PageResult{
		List:  infos,
		Total: total,
	}, nil
}

func (s *serviceImpl) ListMembers(teacherID int64, classID int64) ([]*model.User, error) {
	// 1. 检查权限 (必须是该班级的老师)
	class, err := s.classRepo.GetByID(classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}
	if class.TeacherID != teacherID {
		return nil, errors.New("permission denied")
	}

	return s.classRepo.ListMembers(classID)
}

// generateRandomCode 生成随机字符串，不包含查重逻辑
func (s *serviceImpl) generateRandomCode() (string, error) {
	code := make([]byte, length)
	for j := 0; j < length; j++ {
		// 使用 crypto/rand 生成真随机数
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[j] = charset[num.Int64()]
	}
	return string(code), nil
}
