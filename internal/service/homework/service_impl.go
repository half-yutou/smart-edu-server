package homework

import (
	"errors"
	"time"

	"smarteduhub/internal/config"
	"smarteduhub/internal/model"
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/pkg/ai"
	"smarteduhub/internal/pkg/ocr"
	classRepo "smarteduhub/internal/repository/class"
	homeworkRepo "smarteduhub/internal/repository/homework"
)

type serviceImpl struct {
	homeworkRepo homeworkRepo.Repository
	classRepo    classRepo.Repository
	ocrClient    ocr.Client
	aiClient     ai.Client
}

var _ Service = (*serviceImpl)(nil)

func NewService() Service {
	// 初始化 AI 客户端
	var aiClient ai.Client
	if key := config.Global.AI.DeepSeek.APIKey; key != "" {
		aiClient = ai.NewDeepSeekClient(key, config.Global.AI.DeepSeek.BaseURL)
	} else {
		aiClient = ai.NewMockClient()
	}

	// 初始化 OCR 客户端
	var ocrClient ocr.Client
	if url := config.Global.AI.BaiduOCR.APIURL; url != "" {
		ocrClient = ocr.NewBaiduClient(url, config.Global.AI.BaiduOCR.Token)
	} else {
		ocrClient = ocr.NewMockClient()
	}

	return &serviceImpl{
		homeworkRepo: homeworkRepo.NewRepository(),
		classRepo:    classRepo.NewRepository(),
		ocrClient:    ocrClient,
		aiClient:     aiClient,
	}
}

func (s *serviceImpl) Create(teacherID int64, req *request.CreateHomeworkRequest) error {
	// 1. 检查班级权限
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

	// 2. 构建 Homework 对象
	hw := &model.Homework{
		Title:     req.Title,
		ClassID:   req.ClassID,
		CreatorID: teacherID,
		Deadline:  req.Deadline,
		Questions: make([]model.Question, 0, len(req.Questions)),
	}

	// 3. 构建 Questions
	for i, q := range req.Questions {
		hw.Questions = append(hw.Questions, model.Question{
			QuestionType:  q.QuestionType,
			Content:       q.Content,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Score:         q.Score,
			OrderNum:      i + 1,
		})
	}

	// 4. 保存
	return s.homeworkRepo.Create(hw)
}

func (s *serviceImpl) Delete(operatorID int64, req *request.DeleteHomeworkRequest) error {
	hw, err := s.homeworkRepo.GetByID(req.HomeworkID)
	if err != nil {
		return err
	}
	if hw == nil {
		return errors.New("homework not found")
	}
	if hw.CreatorID != operatorID {
		return errors.New("permission denied")
	}

	return s.homeworkRepo.Delete(req.HomeworkID)
}

func (s *serviceImpl) Update(operatorID int64, req *request.UpdateHomeworkRequest) error {
	// 1. 检查权限
	hw, err := s.homeworkRepo.GetByID(req.HomeworkID)
	if err != nil {
		return err
	}
	if hw == nil {
		return errors.New("homework not found")
	}
	if hw.CreatorID != operatorID {
		return errors.New("permission denied")
	}

	// 2. 更新基础字段
	hw.Title = req.Title
	hw.Deadline = req.Deadline

	// 3. 构建新的题目列表
	// GORM Replace 逻辑：
	// - 如果 Question 有 ID，GORM 会尝试 UPDATE
	// - 如果 Question ID 为 0，GORM 会 INSERT
	// - 如果旧列表里有 ID 而新列表里没有，GORM 会 DELETE (或置空外键)
	newQuestions := make([]model.Question, 0, len(req.Questions))
	for i, q := range req.Questions {
		newQuestions = append(newQuestions, model.Question{
			ID:            q.ID, // 带上 ID
			HomeworkID:    hw.ID,
			QuestionType:  q.QuestionType,
			Content:       q.Content,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Score:         q.Score,
			OrderNum:      i + 1,
		})
	}
	hw.Questions = newQuestions

	return s.homeworkRepo.Update(hw)
}

func (s *serviceImpl) GetByID(id int64) (*model.Homework, error) {
	return s.homeworkRepo.GetByID(id)
}

func (s *serviceImpl) ListByCreator(creatorID int64) ([]*model.Homework, error) {
	return s.homeworkRepo.ListByCreator(creatorID)
}

func (s *serviceImpl) ListByClass(studentID, classID int64) ([]*model.Homework, error) {
	// 校验学生是否在班级中
	isMember, err := s.classRepo.IsMember(classID, studentID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this class")
	}

	return s.homeworkRepo.ListByClass(classID)
}

func (s *serviceImpl) Submit(studentID int64, req *request.SubmitHomeworkRequest) error {
	// 1. 查作业
	hw, err := s.homeworkRepo.GetByID(req.HomeworkID)
	if err != nil {
		return err
	}
	if hw == nil {
		return errors.New("homework not found")
	}
	// 检查截止时间
	if hw.Deadline != nil && time.Now().After(*hw.Deadline) {
		return errors.New("homework deadline exceeded")
	}

	// 2. 准备 Submission
	oldSub, _ := s.homeworkRepo.GetSubmission(req.HomeworkID, studentID)

	// 并发控制：如果作业正在批改中，禁止重复提交
	if oldSub != nil && oldSub.Status == "submitted" {
		return errors.New("submission is being processed, please wait")
	}

	sub := &model.Submission{
		HomeworkID:  req.HomeworkID,
		StudentID:   studentID,
		SubmittedAt: time.Now(),
	}
	if oldSub != nil {
		sub.ID = oldSub.ID
	}

	// 3. 保存初始提交 (状态为 submitted，未批改)
	var details []model.SubmissionDetail

	// 建立题目映射
	qMap := make(map[int64]model.Question)
	for _, q := range hw.Questions {
		qMap[q.ID] = q
	}

	for _, d := range req.Details {
		_, exists := qMap[d.QuestionID]
		if !exists {
			continue
		}

		detail := model.SubmissionDetail{
			QuestionID:     d.QuestionID,
			SubmissionID:   sub.ID, // 显式赋值
			AnswerImageURL: d.ImageURL,
			AnswerContent:  d.TextContent,
			// Score, IsCorrect, Comment 初始为空
		}
		details = append(details, detail)
	}

	sub.Details = details
	sub.Status = "submitted"
	sub.TotalScore = 0

	// 先保存到数据库
	if err := s.homeworkRepo.SaveSubmission(sub); err != nil {
		return err
	}

	// 4. 异步执行 OCR 和 AI 评分
	go s.gradeSubmissionAsync(sub, hw)

	return nil
}

// gradeSubmissionAsync 后台批改逻辑
func (s *serviceImpl) gradeSubmissionAsync(sub *model.Submission, hw *model.Homework) {
	// 重新建立题目映射
	qMap := make(map[int64]model.Question)
	for _, q := range hw.Questions {
		qMap[q.ID] = q
	}

	var totalScore int
	for i := range sub.Details {
		detail := &sub.Details[i]
		q, exists := qMap[detail.QuestionID]
		if !exists {
			continue
		}

		// A. OCR 识别
		if detail.AnswerImageURL != "" && detail.AnswerContent == "" {
			text, err := s.ocrClient.RecognizeBasic(detail.AnswerImageURL)
			if err == nil {
				detail.AnswerContent = text
			} else {
				detail.AnswerContent = "[OCR Failed]"
			}
		}

		// B. AI 评分
		if detail.AnswerContent != "" {
			res, err := s.aiClient.GradeQuestion(q.Content, q.CorrectAnswer, detail.AnswerContent, q.Score)
			if err == nil {
				detail.Score = res.Score
				detail.IsCorrect = res.IsCorrect
				detail.Comment = res.Comment
			} else {
				detail.Comment = "AI评分失败: " + err.Error()
			}
		}

		totalScore += detail.Score
	}

	sub.TotalScore = totalScore
	sub.Status = "graded"

	// 保存批改结果
	// 注意：这里再次调用 SaveSubmission 会触发“删旧插新”，这是安全的。
	_ = s.homeworkRepo.SaveSubmission(sub)
}

func (s *serviceImpl) GetSubmission(homeworkID, studentID int64) (*model.Submission, error) {
	return s.homeworkRepo.GetSubmission(homeworkID, studentID)
}

func (s *serviceImpl) ListSubmissions(teacherID, homeworkID int64) ([]*model.Submission, error) {
	hw, err := s.homeworkRepo.GetByID(homeworkID)
	if err != nil {
		return nil, err
	}
	if hw == nil {
		return nil, errors.New("homework not found")
	}
	if hw.CreatorID != teacherID {
		return nil, errors.New("permission denied")
	}
	return s.homeworkRepo.ListSubmissions(homeworkID)
}

func (s *serviceImpl) GradeSubmission(teacherID int64, req *request.ManualGradeRequest) error {
	sub, err := s.homeworkRepo.GetSubmissionByID(req.SubmissionID)
	if err != nil {
		return err
	}
	if sub == nil {
		return errors.New("submission not found")
	}

	// 校验权限
	hw, _ := s.homeworkRepo.GetByID(sub.HomeworkID)
	if hw.CreatorID != teacherID {
		return errors.New("permission denied")
	}

	// 更新分数
	// 注意：sub.Details 是切片，修改其中的元素
	// 为了高效查找，建个 map
	detailMap := make(map[int64]*model.SubmissionDetail)
	for i := range sub.Details {
		detailMap[sub.Details[i].ID] = &sub.Details[i]
	}

	for _, update := range req.Details {
		if d, ok := detailMap[update.DetailID]; ok {
			d.Score = update.Score
			d.Comment = update.Comment
		}
	}

	// 重新计算总分
	var totalScore int
	for _, d := range sub.Details {
		totalScore += d.Score
	}
	sub.TotalScore = totalScore
	sub.Status = "graded"
	sub.Feedback = "人工复核完成"

	return s.homeworkRepo.SaveSubmission(sub)
}
