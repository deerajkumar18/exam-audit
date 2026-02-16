package auditengine

import (
	"fmt"

	"github.com/deeraj-kumar/exam-audit/config"
	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/deeraj-kumar/exam-audit/service"
	"github.com/deeraj-kumar/exam-audit/util"
)

type ExamAuditHandler interface {
	SubmitAnswer(studentId, examID, questionID, ans string) error
	AuditAnswer(instructorId, examID string) (model.AuditReportResponse, error)
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

func (ea *examAuditHandler) AuditAnswer(instructorId, examID string) (model.AuditReportResponse, error) {
	exams, err := util.ReadExamJSONData(config.Cfg.WorkingDir + "/data/exam_details.json")
	if err != nil {
		return model.AuditReportResponse{}, fmt.Errorf("read exam data failed: %w", err)
	}
	var selectedExam model.Exam
	for _, e := range exams.Exams {
		if e.ExamID == examID {
			selectedExam = e
			break
		}
	}

	students, err := util.ReadStudentsJSONData(config.Cfg.WorkingDir + "/data/students_details.json")
	if err != nil {
		return model.AuditReportResponse{}, fmt.Errorf("read students failed: %w", err)
	}

	answers, err := ea.service.QueryEdittedAnswersByExam(selectedExam, students.Students)
	if err != nil {
		return model.AuditReportResponse{}, fmt.Errorf("failed to query editted answers by exam %s , err - %v", selectedExam.ExamID, err)
	}

	grouped := util.GenerateFlattenedTable(answers)

	adj := util.GenerateAuditReport(grouped)

	return model.AuditReportResponse{ExamID: examID, Report: adj}, nil
}
