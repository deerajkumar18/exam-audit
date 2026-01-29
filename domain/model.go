package model

// Exam, Question, Student and Answer models

type Exams struct {
	Exams []Exam `json:"exams"`
}
type Exam struct {
	ExamID    string     `json:"examID"`
	Questions []Question `json:"questions"`
}

type Question struct {
	QuestionID string `json:"questionID"`
	Question   string `json:"question"`
}

type Students struct {
	Students []Student `json:"students"`
}

type Student struct {
	StudentID   string `json:"studentID"`
	StudentName string `json:"studentName"`
}

type Answer struct {
	QuestionID  string   `json:"questionID"`
	Ans         string   `json:"ans"` // the answer text (kept)
	StudentID   string   `json:"studentID"`
	Edits       []string `json:"edits"`       // edit pattern history (optional, helpful for scoring)
	SubmittedAt int64    `json:"submittedAt"` // timestamps per edit (UNIX seconds) - matches scoring functions
}

// AdjacencyItem describes suspicion score edge
type AdjacencyItem struct {
	StudentA string  `json:"studentA"`
	StudentB string  `json:"studentB"`
	Score    float64 `json:"score"`
	Reason   string  `json:"reason,omitempty"`
}

// AdjacencyList is a map keyed by studentID -> slice of adjacency items
type AdjacencyList map[string][]AdjacencyItem

type AnswerHistoryRecord struct {
	Records []AnswerHistory `json:"answerHistoryRecord"`
}

type AnswerHistory struct {
	TxID      string `json:"txId"`
	Timestamp int64    `json:"timestamp"`
	Value     string `json:"value"`
	IsDelete  bool   `json:"isDelete"`
}

// SubmitAnswerRequest payload
type SubmitAnswerRequest struct {
	StudentID  string `json:"studentId" binding:"required"`
	ExamID     string `json:"examId" binding:"required"`
	QuestionID string `json:"questionId" binding:"required"`
	Ans        string `json:"ans" binding:"required"`
}
