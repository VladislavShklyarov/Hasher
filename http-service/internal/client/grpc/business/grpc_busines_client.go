package business

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "http-service/gen/business"
	"log"
	"time"
)

type BusinessClient struct {
	BusinessLogicClient gen.BusinessLogicClient
}

type BusinessClientInterface interface {
	Process(ctx context.Context, req *gen.Request) (*gen.Response, error)
}

var _ BusinessClientInterface = (*BusinessClient)(nil)

func CreateBusinessClient() *BusinessClient {
	conn, err := grpc.NewClient("Localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to business service: %v", err)
	}

	client := &BusinessClient{
		BusinessLogicClient: gen.NewBusinessLogicClient(conn),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testEntry := createTestEntry()

	resp, err := client.BusinessLogicClient.Process(ctx, testEntry)

	if err != nil {
		log.Fatalf("Business service connected, but test operation failed: %v", err)
	}
	log.Printf("Handshake successful, response: %v", resp)

	return client

}

func createTestEntry() (testEntry *gen.Request) {
	testOperationCalc := &gen.Operation{
		Type:  "calc",
		Op:    "+",
		Var:   "x",
		Left:  "2",
		Right: "3",
	}

	testOperationPrint := &gen.Operation{
		Type: "print",
		Var:  "x",
	}

	testEntry = &gen.Request{Operations: []*gen.Operation{testOperationCalc, testOperationPrint}}
	return testEntry
}
