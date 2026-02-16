package audit_test

import (
	"testing"
	"time"

	"github.com/deeraj-kumar/exam-audit/auditengine"
	"github.com/deeraj-kumar/exam-audit/config"
	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/deeraj-kumar/exam-audit/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuditHandler_HighCollision(t *testing.T) {
	if err := config.LoadConfig(); err != nil {
		return
	}
	mockFabricService := new(mocks.FabricService)
	mockAns := []model.Answer{
		{
			QuestionID:  "q1",
			Ans:         "A",
			StudentID:   "s1",
			SubmittedAt: time.Now().Unix(),
		},
		{
			QuestionID:  "q1",
			Ans:         "C",
			StudentID:   "s2",
			SubmittedAt: time.Now().Unix(),
		},
		{
			QuestionID:  "q1",
			Ans:         "C",
			StudentID:   "s1",
			SubmittedAt: time.Now().Add(10 * time.Second).Unix(),
		},
	}
	mockFabricService.On("QueryEdittedAnswersByExam", mock.Anything, mock.Anything).Return(mockAns, nil)
	h := auditengine.NewExamAuditHandler(mockFabricService)
	resp, err := h.AuditAnswer("1", "exam170126")
	assert.Nil(t, err)
	t.Logf("report - %v", resp.Report)
}

func TestAuditHandler_ModerateCollision(t *testing.T) {
	if err := config.LoadConfig(); err != nil {
		return
	}
	mockFabricService := new(mocks.FabricService)

	mockAns := []model.Answer{
		{
			QuestionID:  "q1",
			Ans:         "A",
			StudentID:   "s1",
			SubmittedAt: time.Now().Unix(),
		},
		{
			QuestionID:  "q1",
			Ans:         "C",
			StudentID:   "s2",
			SubmittedAt: time.Now().Unix(),
		},
		{
			QuestionID:  "q1",
			Ans:         "C",
			StudentID:   "s1",
			SubmittedAt: time.Now().Add(40 * time.Second).Unix(),
		},
	}
	mockFabricService.On("QueryEdittedAnswersByExam", mock.Anything, mock.Anything).Return(mockAns, nil)
	h := auditengine.NewExamAuditHandler(mockFabricService)
	resp, err := h.AuditAnswer("1", "exam170126")
	assert.Nil(t, err)
	t.Logf("report - %v", resp.Report)
}

func TestAuditHandler_NoCollision(t *testing.T) {
	if err := config.LoadConfig(); err != nil {
		return
	}
	mockFabricService := new(mocks.FabricService)

	mockAns := []model.Answer{
		{
			QuestionID:  "q1",
			Ans:         "A",
			StudentID:   "s1",
			SubmittedAt: time.Now().Unix(),
		},
		{
			QuestionID:  "q1",
			Ans:         "C",
			StudentID:   "s2",
			SubmittedAt: time.Now().Unix(),
		},
		{
			QuestionID:  "q1",
			Ans:         "D",
			StudentID:   "s1",
			SubmittedAt: time.Now().Add(40 * time.Second).Unix(),
		},
	}
	mockFabricService.On("QueryEdittedAnswersByExam", mock.Anything, mock.Anything).Return(mockAns, nil)
	h := auditengine.NewExamAuditHandler(mockFabricService)
	resp, err := h.AuditAnswer("1", "exam170126")
	assert.Nil(t, err)
	assert.Empty(t, resp.Report)
}
