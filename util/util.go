package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/deeraj-kumar/exam-audit/config"
	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

// ReadExamJSONData reads filename and unmarshals to []model.Exam
func ReadExamJSONData(filename string) (model.Exams, error) {
	log.Printf("exam details file path - %s", filename)
	f, err := os.Open(filename)
	if err != nil {
		return model.Exams{}, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return model.Exams{}, err
	}
	var exams model.Exams
	if err := json.Unmarshal(b, &exams); err != nil {
		return model.Exams{}, err
	}
	return exams, nil
}

// ReadStudentsJSONData reads filename and unmarshals to []model.Student
func ReadStudentsJSONData(filename string) (model.Students, error) {
	f, err := os.Open(filename)
	if err != nil {
		return model.Students{}, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return model.Students{}, err
	}
	var students model.Students
	if err := json.Unmarshal(b, &students); err != nil {
		return model.Students{}, err
	}
	return students, nil
}

// GenerateFlattenedTable groups answers by student id
func GenerateFlattenedTable(editedAnswerList []model.Answer) map[string]map[string][]model.AnswerRevision {
	out := make(map[string]map[string][]model.AnswerRevision)
	for _, ans := range editedAnswerList {
		if _, ok := out[ans.StudentID]; !ok {
			out[ans.StudentID] = make(map[string][]model.AnswerRevision)
		}
		out[ans.StudentID][ans.QuestionID] = append(out[ans.StudentID][ans.QuestionID], model.AnswerRevision{SubmittedAt: ans.SubmittedAt, Ans: ans.Ans})
	}
	return out
}

// GenerateAuditReport compares every pair of students and produces adjacency list
func GenerateAuditReport(studentAnswersMap map[string]map[string][]model.AnswerRevision) (record model.AdjacencyList) {
	susScoreThreshold := config.Cfg.SuspicionScoreThreshold

	studentIDs := make([]string, 0, len(studentAnswersMap))
	for sid := range studentAnswersMap {
		studentIDs = append(studentIDs, sid)
	}

	for i := 0; i < len(studentIDs); i++ {
		aID := studentIDs[i]
		std1AnsMap := studentAnswersMap[aID]
		for qID := range std1AnsMap {
			std1AnswerRevisions := std1AnsMap[qID]
			for j := i + 1; j < len(studentIDs); j++ {
				bID := studentIDs[j]
				std2AnswerRevisions := studentAnswersMap[bID][qID]

				score := questionScore(std1AnswerRevisions, std2AnswerRevisions)
				log.Printf("Audit score: %f between Students : %s - %s", score, aID, bID)
				if score <= 0 || score <= susScoreThreshold {
					continue
				}
				adj := model.AdjacencyItem{
					StudentA: aID,
					StudentB: bID,
					Score:    score,
				}
				record = append(record, adj)
			}
		}
	}
	return
}

func questionScore(stdA, stdB []model.AnswerRevision) float64 {
	finalAnsA := stdA[len(stdA)-1].Ans
	finalAnsB := stdB[len(stdB)-1].Ans

	var stdATimeStamps, stdBTimeStamps []int64
	for _, t := range stdA {
		stdATimeStamps = append(stdATimeStamps, t.SubmittedAt)
	}

	for _, t := range stdB {
		stdBTimeStamps = append(stdBTimeStamps, t.SubmittedAt)
	}

	var stdAEdits, stdBEdits []string
	for _, t := range stdA {
		stdAEdits = append(stdAEdits, t.Ans)
	}

	for _, t := range stdB {
		stdBEdits = append(stdBEdits, t.Ans)
	}

	as := answerSimilarity(finalAnsA, finalAnsB)
	ts := timeCorrelation(stdATimeStamps, stdBTimeStamps)
	es := editPatternScore(stdAEdits, stdBEdits)

	return (0.5 * as) + (0.3 * ts) + (0.2 * es)
}

func answerSimilarity(a, b string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	if a == b {
		return 1.0
	}

	maxLen := math.Max(float64(len(a)), float64(len(b)))
	minLen := math.Min(float64(len(a)), float64(len(b)))

	return minLen / maxLen
}

func timeCorrelation(aTime, bTime []int64) float64 {
	minLen := int(math.Min(float64(len(aTime)), float64(len(bTime))))
	if len(aTime) > minLen {
		aTime = aTime[minLen:]
	}

	if len(bTime) > minLen {
		bTime = bTime[minLen:]
	}

	var sum float64
	for i := minLen - 1; i >= 0; i-- {
		diff := math.Abs(float64(aTime[i] - bTime[i]))
		score := math.Max(0, 1-(diff/60.0))
		sum += score
	}
	return sum / float64(minLen)
}

func editPatternScore(aEdits, bEdits []string) float64 {
	minLen := int(math.Min(float64(len(aEdits)), float64(len(bEdits))))
	if minLen == 0 {
		return 0
	}

	if len(aEdits) > minLen {
		aEdits = aEdits[minLen:]
	}

	if len(bEdits) > minLen {
		bEdits = bEdits[minLen:]
	}

	var match float64
	for i := minLen - 1; i >= 0; i-- {
		if aEdits[i] == bEdits[i] {
			match++
		}
	}
	return match / float64(minLen)
}

func LoadX509Identity(certPath, mspID string) (*identity.X509Identity, error) {
	pemBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading cert file: %w", err)
	}

	cert, err := identity.CertificateFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed parsing cert: %w", err)
	}

	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		return nil, fmt.Errorf("failed creating identity: %w", err)
	}

	return id, nil
}

func LoadSigner(keyPath string) (identity.Sign, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading private key: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed loading private key: %w", err)
	}

	signer, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed creating signer: %w", err)
	}

	return signer, nil
}
