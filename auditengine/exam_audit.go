package auditengine

import (
	"fmt"

	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/deeraj-kumar/exam-audit/service"
	"github.com/deeraj-kumar/exam-audit/util"
)

const (
	examDetailsJSON     = "/home/kkdee/go/src/github.com/deerajkumar18/exam-audit/data/exam_details.json"
	studentsDetailsJSON = "/home/kkdee/go/src/github.com/deerajkumar18/exam-audit/data/students_details.json"
)

type ExamAuditHandler interface {
	SubmitAnswer(studentId, examID, questionID, ans string) error
	AuditAnswer(instructorId, examID string) (model.AdjacencyList, error)
}

type examAuditHandler struct {
	service service.FabricService
}

func NewExamAuditHandler(svc service.FabricService) ExamAuditHandler {
	return &examAuditHandler{service: svc}
}

func (ea *examAuditHandler) SubmitAnswer(studentId, examID, questionID, ans string) error {
	if err := ea.service.SetAnswer(studentId, examID, questionID, ans); err != nil {
		return fmt.Errorf("failed to submit answer . student id - %s , question id - %s , exam id - %s , err - %v", studentId, examID, questionID, err)
	}
	return nil
}

func (ea *examAuditHandler) AuditAnswer(instructorId, examID string) (model.AdjacencyList, error) {
	exams, err := util.ReadExamJSONData(examDetailsJSON)
	if err != nil {
		return nil, fmt.Errorf("read exam data failed: %w", err)
	}
	var selectedExam model.Exam
	for _, e := range exams.Exams {
		if e.ExamID == examID {
			selectedExam = e
			break
		}
	}

	students, err := util.ReadStudentsJSONData(studentsDetailsJSON)
	if err != nil {
		return nil, fmt.Errorf("read students failed: %w", err)
	}

	answers, err := ea.service.QueryEdittedAnswersByExam(selectedExam, students.Students)
	if err != nil {
		return nil, fmt.Errorf("failed to query editted answers by exam %s , err - %v", selectedExam.ExamID, err)
	}

	grouped := util.GenerateFlattenedTable(answers)

	adj := util.GenerateAuditReport(grouped)

	return adj, nil
}
