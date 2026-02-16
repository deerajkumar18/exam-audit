package model

type Config struct {
	FabricParams struct {
		PeerEP          string `mapstructure:"peer_endpoint"`
		PeerTlsCertPath string `mapstructure:"peer_tls_cert_path"`
		ChannelName     string `mapstructure:"channel_name"`
		ChaincodeName   string `mapstructure:"chaincode_name"`
		MspID           string `mapstructure:"mspid"`
		FabricIdentity  struct {
			CertPath string `mapstructure:"cert_path"`
			KeyPath  string `mapstructure:"keypath"`
		} `mapstructure:"fabric_identity"`
	} `mapstructure:"fabric_params"`
	SuspicionScoreThreshold float64 `mapstructure:"suspicion_score_threshold"`
	WorkingDir              string  `mapstructure:"working_dir"`
}

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
	QuestionID  string `json:"questionID"`
	Ans         string `json:"ans"`
	StudentID   string `json:"studentID"`
	SubmittedAt int64  `json:"submittedAt"`
}

type AnswerRevision struct {
	SubmittedAt int64  `json:"submittedAt"`
	Ans         string `json:"ans"`
}

type AdjacencyItem struct {
	StudentA string  `json:"studentA"`
	StudentB string  `json:"studentB"`
	Score    float64 `json:"score"`
	//Reason   string  `json:"reason,omitempty"`
}

type AdjacencyList []AdjacencyItem

type AuditReportResponse struct {
	ExamID string        `json:"examID" binding:"required"`
	Report AdjacencyList `json:"report" binding:"required"`
}

type AnswerHistoryRecord struct {
	Records []AnswerHistory `json:"answerHistoryRecord"`
}

type AnswerHistory struct {
	TxID      string `json:"txId"`
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
	IsDelete  bool   `json:"isDelete"`
}

type SubmitAnswerRequest struct {
	StudentID  string `json:"studentId" binding:"required"`
	ExamID     string `json:"examId" binding:"required"`
	QuestionID string `json:"questionId" binding:"required"`
	Ans        string `json:"ans" binding:"required"`
}
