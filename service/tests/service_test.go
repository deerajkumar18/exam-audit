package fabricsvctest

import (
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/deeraj-kumar/exam-audit/service"
	"github.com/deeraj-kumar/exam-audit/service/contract"
	fabricutils "github.com/deeraj-kumar/exam-audit/service/fabricUtils"
	"github.com/deeraj-kumar/exam-audit/service/mocks"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	mockPeerEP        = "localhost:7051"
	mockMspID         = "msp1"
	mockChannelName   = "mychannel"
	mockChaincodeName = "samplecc"
)

func TestSetAnswer_Success(t *testing.T) {
	mockContract := new(mocks.Contract)

	originalGetCertPool := fabricutils.GetCertPool
	fabricutils.GetCertPool = func(peerTLSCertPath, dir string) (*x509.CertPool, error) {
		return &x509.CertPool{}, nil
	}
	defer func() {
		fabricutils.GetCertPool = originalGetCertPool
	}()

	originalGetFabricGrpcClient := fabricutils.GetGrpcClient
	fabricutils.GetGrpcClient = func(peerEndpoint string, cp *x509.CertPool) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}
	defer func() {
		fabricutils.GetGrpcClient = originalGetFabricGrpcClient
	}()

	oldGetFabricGateway := fabricutils.GetFabricGateway
	fabricutils.GetFabricGateway = func(id identity.Identity, signer identity.Sign, grpcClient *grpc.ClientConn) (*client.Gateway, error) {
		return &client.Gateway{}, nil
	}
	defer func() {
		fabricutils.GetFabricGateway = oldGetFabricGateway
	}()

	oldGetFabricId := fabricutils.GetId
	fabricutils.GetId = func(certPath, dir string) (*identity.X509Identity, error) {
		return &identity.X509Identity{}, nil
	}
	defer func() {
		fabricutils.GetId = oldGetFabricId
	}()

	oldGetFabricSigner := fabricutils.GetSigner
	fabricutils.GetSigner = func(keyPath, dir string) (identity.Sign, error) {
		var x identity.Sign
		return x, nil
	}
	defer func() {
		fabricutils.GetSigner = oldGetFabricSigner
	}()

	originalGetContract := fabricutils.GetContract
	fabricutils.GetContract = func(gw *client.Gateway, channelName, chainCodeName string) contract.Contract {
		return mockContract
	}
	defer func() {
		fabricutils.GetContract = originalGetContract
	}()

	mockContract.
		On("SubmitTransaction", "SetAnswer", "Answer~exam1~q1~s1", "A").
		Return([]byte("OK"), nil)

	fabricSvc, err := service.NewFabricService(mockPeerEP, "", "", "", mockMspID, mockChannelName, mockChaincodeName)
	assert.Nil(t, err)

	serviceErr := fabricSvc.SetAnswer("s1", "exam1", "q1", "A")
	assert.Nil(t, serviceErr)
}

func TestSetAnswer_Fail(t *testing.T) {
	mockContract := new(mocks.Contract)

	originalGetCertPool := fabricutils.GetCertPool
	fabricutils.GetCertPool = func(peerTLSCertPath, dir string) (*x509.CertPool, error) {
		return &x509.CertPool{}, nil
	}
	defer func() {
		fabricutils.GetCertPool = originalGetCertPool
	}()

	originalGetFabricGrpcClient := fabricutils.GetGrpcClient
	fabricutils.GetGrpcClient = func(peerEndpoint string, cp *x509.CertPool) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}
	defer func() {
		fabricutils.GetGrpcClient = originalGetFabricGrpcClient
	}()

	oldGetFabricGateway := fabricutils.GetFabricGateway
	fabricutils.GetFabricGateway = func(id identity.Identity, signer identity.Sign, grpcClient *grpc.ClientConn) (*client.Gateway, error) {
		return &client.Gateway{}, nil
	}
	defer func() {
		fabricutils.GetFabricGateway = oldGetFabricGateway
	}()

	oldGetFabricId := fabricutils.GetId
	fabricutils.GetId = func(certPath, dir string) (*identity.X509Identity, error) {
		return &identity.X509Identity{}, nil
	}
	defer func() {
		fabricutils.GetId = oldGetFabricId
	}()

	oldGetFabricSigner := fabricutils.GetSigner
	fabricutils.GetSigner = func(keyPath, dir string) (identity.Sign, error) {
		var x identity.Sign
		return x, nil
	}
	defer func() {
		fabricutils.GetSigner = oldGetFabricSigner
	}()

	originalGetContract := fabricutils.GetContract
	fabricutils.GetContract = func(gw *client.Gateway, channelName, chainCodeName string) contract.Contract {
		return mockContract
	}
	defer func() {
		fabricutils.GetContract = originalGetContract
	}()

	mockContract.
		On("SubmitTransaction", "SetAnswer", "Answer~exam1~q1~s1", "A").
		Return([]byte("OK"), fmt.Errorf("some error"))

	fabricSvc, err := service.NewFabricService(mockPeerEP, "", "", "", mockMspID, mockChannelName, mockChaincodeName)
	assert.Nil(t, err)

	serviceErr := fabricSvc.SetAnswer("s1", "exam1", "q1", "A")
	assert.NotNil(t, serviceErr)
	assert.Contains(t, serviceErr.Error(), "some error")
}
