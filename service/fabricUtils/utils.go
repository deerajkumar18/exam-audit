package fabricutils

import (
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/deeraj-kumar/exam-audit/service/contract"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var GetCertPool = func(peerTLSCertPath, dir string) (*x509.CertPool, error) {
	peerCertBytes, err := os.ReadFile(dir + peerTLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading peer TLS cert: %w", err)
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(peerCertBytes) {
		return nil, fmt.Errorf("failed appending peer TLS cert")
	}
	return cp, nil
}

var GetGrpcClient = func(peerEndpoint string, cp *x509.CertPool) (*grpc.ClientConn, error) {
	grpcConn, err := grpc.NewClient(
		peerEndpoint,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating grpc connection: %w", err)
	}
	return grpcConn, nil
}

var GetId = func(certPath, dir string) (*identity.X509Identity, error) {
	certPem, err := os.ReadFile(dir + certPath)
	if err != nil {
		return nil, err
	}
	cert, err := identity.CertificateFromPEM(certPem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	id, err := identity.NewX509Identity(
		"Org1MSP",
		cert,
	)
	if err != nil {
		return nil, err
	}

	return id, nil

}

var GetSigner = func(keyPath, dir string) (identity.Sign, error) {
	keyDir := dir + keyPath
	keyPem, err := os.ReadFile(keyDir)
	if err != nil {
		return nil, err
	}
	privateKey, err := identity.PrivateKeyFromPEM(keyPem)
	if err != nil {
		return nil, err
	}

	signer, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

var GetFabricGateway = func(id identity.Identity, signer identity.Sign, grpcClient *grpc.ClientConn) (*client.Gateway, error) {
	gw, err := client.Connect(
		id,
		client.WithSign(signer),
		client.WithClientConnection(grpcClient),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(15*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("gateway connect failed: %w", err)
	}
	return gw, nil
}

var GetContract = func(gw *client.Gateway, channelName, chainCodeName string) contract.Contract {
	c := gw.GetNetwork(channelName).GetContract(chainCodeName)
	return contract.NewFabricContract(c)
}
