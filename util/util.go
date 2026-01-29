package util

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"

	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

// ReadExamJSONData reads filename and unmarshals to []model.Exam
func ReadExamJSONData(filename string) (model.Exams, error) {
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
func GenerateFlattenedTable(editedAnswerList []model.Answer) map[string][]model.Answer {
	out := make(map[string][]model.Answer)
	for _, ans := range editedAnswerList {
		out[ans.StudentID] = append(out[ans.StudentID], ans)
	}
	return out
}

// GenerateAuditReport compares every pair of students and produces adjacency list
func GenerateAuditReport(studentAnswersMap map[string][]model.Answer) model.AdjacencyList {
	adj := make(model.AdjacencyList)
	// list of student ids
	studentIDs := make([]string, 0, len(studentAnswersMap))
	for sid := range studentAnswersMap {
		studentIDs = append(studentIDs, sid)
	}

	for i := 0; i < len(studentIDs); i++ {
		aID := studentIDs[i]
		aAnswers := studentAnswersMap[aID]
		// only consider students with more than 1 answer record
		if len(aAnswers) <= 1 {
			continue
		}
		for j := i + 1; j < len(studentIDs); j++ {
			bID := studentIDs[j]
			bAnswers := studentAnswersMap[bID]
			if len(bAnswers) <= 1 {
				continue
			}
			score := questionScore(aAnswers, bAnswers)
			if score <= 0 {
				continue
			}
			edgeA := model.AdjacencyItem{
				StudentA: aID,
				StudentB: bID,
				Score:    score,
				Reason:   fmt.Sprintf("similarity-%.3f", score),
			}
			edgeB := model.AdjacencyItem{
				StudentA: bID,
				StudentB: aID,
				Score:    score,
				Reason:   fmt.Sprintf("similarity-%.3f", score),
			}
			adj[aID] = append(adj[aID], edgeA)
			adj[bID] = append(adj[bID], edgeB)
		}
	}
	return adj
}

func questionScore(stdA, stdB []model.Answer) float64 {
	finalAnsA := stdA[len(stdA)-1]
	finalAnsB := stdB[len(stdB)-1]

	stdATimeStamps := []int64{}
	for i := 0; i < len(stdA); i++ {
		stdATimeStamps = append(stdATimeStamps, stdA[i].SubmittedAt)
	}
	stdBTimeStamps := []int64{}
	for i := 0; i < len(stdB); i++ {
		stdBTimeStamps = append(stdBTimeStamps, stdB[i].SubmittedAt)
	}

	as := answerSimilarity(finalAnsA.Ans, finalAnsB.Ans)
	ts := timeCorrelation(stdATimeStamps, stdBTimeStamps)
	es := editPatternScore(finalAnsA.Edits, finalAnsB.Edits)

	// Weighted combination (tuneable)
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
	if len(aTime) == 0 || len(bTime) == 0 {
		return 0
	}

	minLen := int(math.Min(float64(len(aTime)), float64(len(bTime))))

	var sum float64
	for i := 0; i < minLen; i++ {
		diff := math.Abs(float64(aTime[i] - bTime[i]))
		// normalize: assuming differences > 60 seconds are "not correlated"
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

	var match float64
	for i := 0; i < minLen; i++ {
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
