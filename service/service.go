package service

import (
	"encoding/json"
	"fmt"
	"os"

	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/deeraj-kumar/exam-audit/service/contract"
	fabricutils "github.com/deeraj-kumar/exam-audit/service/fabricUtils"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type FabricService interface {
	SetAnswer(studentId, examID, questionID, ans string) error
	QueryEdittedAnswersByExam(exam model.Exam, students []model.Student) ([]model.Answer, error)
	Close()
}

type fabricService struct {
	gateway  *client.Gateway
	contract contract.Contract
}

func NewFabricService(peerEndpoint, peerTLSCertPath, certPath, keyPath, mspID, channelName, chaincodeName string) (FabricService, error) {

	// dir, err := os.Getwd()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get the current working directory : %v", err)
	// }
	// dir += "/config/msp/"
	// // Load TLS cert for peer (Gateway uses TLS)
	// peerCertBytes, err := os.ReadFile(dir + peerTLSCertPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed reading peer TLS cert: %w", err)
	// }

	// cp := x509.NewCertPool()
	// if !cp.AppendCertsFromPEM(peerCertBytes) {
	// 	return nil, fmt.Errorf("failed appending peer TLS cert")
	// }
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get the current working directory : %v", err)
	}
	dir += "/config/msp/"

	cp, err := fabricutils.GetCertPool(peerTLSCertPath, dir)
	if err != nil {
		return nil, err
	}
	// grpcConn, err := grpc.NewClient(
	// 	peerEndpoint,
	// 	grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, "")),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed creating grpc connection: %w", err)
	// }

	// // Load identity
	// id, err := util.LoadX509Identity(dir+certPath, mspID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed loading identity: %w", err)
	// }
	grpcClient, err := fabricutils.GetGrpcClient(peerEndpoint, cp)
	if err != nil {
		return nil, err
	}

	// certPem, err := os.ReadFile(dir + certPath)
	// if err != nil {
	// 	return nil, err
	// }
	// cert, err := identity.CertificateFromPEM(certPem)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse certificate: %w", err)
	// }
	// id, err := identity.NewX509Identity(
	// 	"Org1MSP",
	// 	cert,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	id, err := fabricutils.GetId(certPath, dir)
	if err != nil {
		return nil, err
	}

	// keyDir := dir + keyPath
	// keyPem, err := os.ReadFile(keyDir)
	// if err != nil {
	// 	return nil, err
	// }
	// privateKey, err := identity.PrivateKeyFromPEM(keyPem)
	// if err != nil {
	// 	return nil, err
	// }

	// signer, err := identity.NewPrivateKeySign(privateKey)
	// if err != nil {
	// 	return nil, err
	// }

	signer, err := fabricutils.GetSigner(keyPath, dir)
	if err != nil {
		return nil, err
	}

	// Connect Gateway
	// gw, err := client.Connect(
	// 	id,
	// 	client.WithSign(signer),
	// 	client.WithClientConnection(grpcClient),
	// 	client.WithEvaluateTimeout(5*time.Second),
	// 	client.WithEndorseTimeout(15*time.Second),
	// 	client.WithSubmitTimeout(15*time.Second),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("gateway connect failed: %w", err)
	// }
	gw, err := fabricutils.GetFabricGateway(id, signer, grpcClient)
	if err != nil {
		return nil, err
	}

	//network := gw.GetNetwork(channelName)
	c := fabricutils.GetContract(gw, channelName, chaincodeName)

	return &fabricService{
		gateway:  gw,
		contract: c,
	}, nil
}

func (s *fabricService) Close() {
	if s.gateway != nil {
		_ = s.gateway.Close()
	}
}

func (s *fabricService) SetAnswer(studentId, examID, questionID, ans string) error {
	if s.contract == nil {
		return fmt.Errorf("contract not initialized")
	}

	compositeKey := fmt.Sprintf("Answer~%s~%s~%s", examID, questionID, studentId)

	_, err := s.contract.SubmitTransaction("SetAnswer", compositeKey, ans)
	if err != nil {
		return fmt.Errorf("failed submitting SetAnswer: %w", err)
	}
	return nil
}

func (s *fabricService) QueryEdittedAnswersByExam(exam model.Exam, students []model.Student) ([]model.Answer, error) {
	var answers []model.Answer
	if s.contract == nil {
		return nil, fmt.Errorf("contract not initialized")
	}
	examID := exam.ExamID
	for _, q := range exam.Questions {
		var answerHistoryRecords []model.AnswerHistory
		for _, std := range students {
			key := fmt.Sprintf("Answer~%s~%s~%s", examID, q.QuestionID, std.StudentID)
			transactionResp, err := s.contract.EvaluateTransaction("GetAnswerRevisionHistory", key)
			if err != nil {
				return nil, fmt.Errorf("failed to get the answer revision history for the key %s , due to %v", key, err)
			}
			if len(transactionResp) == 0 {
				return nil, fmt.Errorf("transaction response data can't be empty , key %s", key)
			}
			fmt.Printf("Ledger raw response for History Query , key - %s , resp - %s", key, string(transactionResp))
			if err := json.Unmarshal(transactionResp, &answerHistoryRecords); err != nil {
				return nil, fmt.Errorf("failed to unmarshal %v , err - %v ", transactionResp, err)
			}
			for _, record := range answerHistoryRecords {
				answers = append(answers, model.Answer{Ans: record.Value, QuestionID: q.QuestionID, StudentID: std.StudentID, SubmittedAt: record.Timestamp})
			}

		}

	}
	return answers, nil
}
