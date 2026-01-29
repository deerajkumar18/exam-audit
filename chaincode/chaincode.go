package chaincode

import (
	"encoding/json"
	"fmt"

	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AnswerContract struct {
	contractapi.Contract
}

func (c *AnswerContract) SetAnswer(ctx contractapi.TransactionContextInterface, key string, answer string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	return ctx.GetStub().PutState(key, []byte(answer))
}

// ====== GetAnswerRevisionHistory ======
func (c *AnswerContract) GetAnswerRevisionHistory(
	ctx contractapi.TransactionContextInterface,
	examID string,
	questionID string,
	studentID string,
) ([]byte, error) {

	if examID == "" || questionID == "" || studentID == "" {
		return nil, fmt.Errorf("all IDs must be provided")
	}

	key := fmt.Sprintf("Answer~%s~%s~%s", examID, questionID, studentID)

	// Get iterator for full history
	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve history: %v", err)
	}
	defer iter.Close()

	var history []model.AnswerHistory

	for iter.HasNext() {
		mod, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("failed iterating history: %v", err)
		}

		record := model.AnswerHistory{
			TxID:      mod.TxId,
			Timestamp: mod.Timestamp.Seconds,
			Value:     string(mod.Value),
			IsDelete:  mod.IsDelete,
		}

		history = append(history, record)
	}

	return json.Marshal(history)
}
