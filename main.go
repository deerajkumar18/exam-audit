package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deeraj-kumar/exam-audit/auditengine"
	"github.com/deeraj-kumar/exam-audit/handlers"
	"github.com/deeraj-kumar/exam-audit/service"
	"github.com/gin-gonic/gin"
)

func main() {
	peerEndpoint := getenv("FABRIC_PEER_ENDPOINT", "localhost:7051")
	peerTlsCertPath := getenv("FABRIC_PEER_TLS_CERT_PATH", "peer0.org1.example.com/tls/ca.cert")
	certPath := getenv("FABRIC_CERT_PATH", "admin@org1/signcerts/cert.pem")
	keyPath := getenv("FABRIC_KEY_PATH", "admin@org1/keystore/priv_sk")
	mspID := getenv("FABRIC_MSPID", "Org1MSP")
	channelName := getenv("FABRIC_CHANNEL", "mychannel")
	chaincodeName := getenv("FABRIC_CHAINCODE", "exam")

	fabricSvc, err := service.NewFabricService(peerEndpoint, peerTlsCertPath, certPath, keyPath, mspID, channelName, chaincodeName)
	if err != nil {
		log.Fatalf("failed to initialize fabric service: %v", err)
	}
	defer fabricSvc.Close()

	examAuditHandler := auditengine.NewExamAuditHandler(fabricSvc)

	r := gin.Default()
	h := handlers.NewHandler(examAuditHandler)
	h.RegisterRoutes(r)

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()
	log.Println("server running on :8080")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
