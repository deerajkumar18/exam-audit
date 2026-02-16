package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AnswerContract struct {
	contractapi.Contract
}

type Answer struct {
	AnsString string `json:"ans"`
}

type AnswerSubmissionDetail struct {
	TxID      string `json:"txId"`
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
	IsDelete  bool   `json:"isDelete"`
}

func (t *AnswerContract) SubmissionExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", key, err)
	}

	return assetBytes != nil, nil
}

func (c *AnswerContract) SetAnswer(ctx contractapi.TransactionContextInterface, key string, answer string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	ans := &Answer{
		AnsString: answer,
	}
	bytes, err := json.Marshal(ans)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, bytes)
}

func (c *AnswerContract) GetAnswerRevisionHistory(
	ctx contractapi.TransactionContextInterface, key string) ([]AnswerSubmissionDetail, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	log.Printf("GetAnswerRevisionHistory : Key - %s", key)

	exists, err := c.SubmissionExists(ctx, key)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("key %s does not have a world state existing in the ledger", key)
	}

	// Get iterator for full history
	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve history: %v", err)
	}
	defer iter.Close()

	var submissionRecord []AnswerSubmissionDetail

	for iter.HasNext() {
		resp, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("failed iterating history: %v", err)
		}

		var answer Answer
		if len(resp.Value) > 0 {
			err = json.Unmarshal(resp.Value, &answer)
			if err != nil {
				return nil, err
			}
		}

		record := AnswerSubmissionDetail{
			TxID:      resp.TxId,
			Timestamp: resp.Timestamp.Seconds,
			Value:     answer.AnsString,
			IsDelete:  resp.IsDelete,
		}

		submissionRecord = append(submissionRecord, record)
	}

	return submissionRecord, nil
}

func main() {
	cc, err := contractapi.NewChaincode(
		&AnswerContract{},
	)

	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
