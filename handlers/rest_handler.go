package handlers

import (
	"net/http"

	"github.com/deeraj-kumar/exam-audit/auditengine"
	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(r *gin.Engine)
}

type handlerImpl struct {
	auditEngine auditengine.ExamAuditHandler
}

func NewHandler(ae auditengine.ExamAuditHandler) Handler {
	return &handlerImpl{
		auditEngine: ae,
	}
}

func (h *handlerImpl) RegisterRoutes(r *gin.Engine) {
	r.POST("/submit-answer", h.SubmitAnswer)
	r.GET("/audit-answer", h.AuditAnswer)
}

func (h *handlerImpl) SubmitAnswer(c *gin.Context) {
	var req model.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.auditEngine.SubmitAnswer(req.StudentID, req.ExamID, req.QuestionID, req.Ans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *handlerImpl) AuditAnswer(c *gin.Context) {
	instructorId := c.Query("instructorId")
	examID := c.Query("examID")

	if instructorId == "" || examID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instructorId and examID are required"})
		return
	}

	auditResp, err := h.auditEngine.AuditAnswer(instructorId, examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, auditResp)
}
