package ai

type GradingResult struct {
	Score     int
	IsCorrect bool
	Comment   string
}

type Client interface {
	// GradeQuestion 批改单个题目
	GradeQuestion(questionContent, standardAnswer, studentAnswer string, fullScore int) (*GradingResult, error)
}