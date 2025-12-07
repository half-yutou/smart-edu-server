package ai

import "strings"

type mockClient struct{}

func NewMockClient() Client {
	return &mockClient{}
}

func (c *mockClient) GradeQuestion(q, std, stu string, fullScore int) (*GradingResult, error) {
	// 简单模拟
	// 如果是选择题（标准答案短，且不含空格），做精确匹配
	// 这里假设标准答案长度小于等于5且不包含中文标点通常是选择题
	if len(std) > 0 && len(std) <= 5 {
		isCorrect := strings.EqualFold(strings.TrimSpace(std), strings.TrimSpace(stu))
		score := 0
		if isCorrect {
			score = fullScore
		}
		return &GradingResult{
			Score:     score,
			IsCorrect: isCorrect,
			Comment:   "Mock AI: 自动判定",
		}, nil
	}

	// 简答题，默认给一半分
	return &GradingResult{
		Score:     fullScore / 2,
		IsCorrect: false,
		Comment:   "Mock AI: 待人工复核",
	}, nil
}