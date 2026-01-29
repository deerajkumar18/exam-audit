package contract

import "github.com/hyperledger/fabric-gateway/pkg/client"

type Contract interface {
	SubmitTransaction(name string, args ...string) ([]byte, error)
	EvaluateTransaction(name string, args ...string) ([]byte, error)
}

type fabricContract struct {
	contract *client.Contract
}

func NewFabricContract(c *client.Contract) Contract {
	return &fabricContract{contract: c}
}

func (f *fabricContract) SubmitTransaction(name string, args ...string) ([]byte, error) {
	return f.contract.SubmitTransaction(name, args...)
}

func (f *fabricContract) EvaluateTransaction(name string, args ...string) ([]byte, error) {
	return f.contract.EvaluateTransaction(name, args...)
}
