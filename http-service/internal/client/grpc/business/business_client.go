package business

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "http-service/gen"
	"http-service/internal/config"
	"log"
	"time"
)

const connectionError = "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:9091: connect: connection refused\""

type BusinessClient struct {
	GRPCClient gen.BusinessLogicClient // Wrap GRPCclient
}

func CreateBusinessClient(cfg *config.Config) *BusinessClient {
	conn, err := grpc.NewClient(cfg.BusinessAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to business server: %v", err)
	}

	client := &BusinessClient{
		GRPCClient: gen.NewBusinessLogicClient(conn),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	testEntry := createTestEntry()

	resp, err := client.GRPCClient.Process(ctx, testEntry)

	if err != nil {
		fmt.Println(err.Error())

		if err.Error() == connectionError {
			log.Println("Business server connection error")
			return nil
		} else {
			log.Printf("Business server connected, but test message failed: %v", err)
			return nil
		}
	}

	log.Printf("Handshake successful, response: %v", resp)

	return client

}

func createTestEntry() (testEntry *gen.OperationRequest) {
	testOperationCalc := &gen.Operation{
		Type:  "Test calc",
		Op:    "Test +",
		Var:   "Test var x",
		Left:  "Test 2",
		Right: "Test 3",
	}

	testOperationPrint := &gen.Operation{
		Type: "Test print",
		Var:  "test x",
	}

	testEntry = &gen.OperationRequest{Operations: []*gen.Operation{testOperationCalc, testOperationPrint}}
	return testEntry
}
