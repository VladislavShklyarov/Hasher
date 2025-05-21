package business

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gen "http-service/gen"
	"log"
	"time"
)

type BusinessClient struct {
	GRPCClient gen.BusinessLogicClient // Wrap GRPCclient
}

//type BusinessClientInterface interface {
//	Process(ctx context.Context, req *gen.Request) (*gen.Response, error)
//}
//
//var _ BusinessClientInterface = (*BusinessClient)(nil)

func CreateBusinessClient() *BusinessClient {
	conn, err := grpc.NewClient("Localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to business server: %v", err)
	}

	client := &BusinessClient{
		GRPCClient: gen.NewBusinessLogicClient(conn),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testEntry := createTestEntry()

	resp, err := client.GRPCClient.Process(ctx, testEntry)

	if err != nil {
		log.Fatalf("Business server connected, but test operation failed: %v", err)
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
