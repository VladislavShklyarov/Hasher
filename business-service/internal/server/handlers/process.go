package server

import (
	"business-service/gen"
	logGRPC "business-service/internal/clients/grpc/log"
	"business-service/internal/clients/kafka"
	"business-service/internal/config"
	"business-service/internal/logic"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/durationpb"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

type BusinessLogicManager struct {
	gen.UnimplementedBusinessLogicServer
	GRPCClient *logGRPC.LogClient
}

func (blm *BusinessLogicManager) Process(ctx context.Context, req *gen.OperationRequest) (*gen.OperationResponse, error) {

	cfg := config.Load()

	operations := req.GetOperations()

	aliveVars, graph := logic.FindAliveVariables(operations)

	err := logic.ExportToDOT(operations, aliveVars, graph, "graph.dot")

	if err != nil {
		fmt.Println("Error during export:", err)
	} else {
		cmd := exec.Command("dot", "-Tpng", "graph.dot", "-o", "graph.png")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error generating PNG:", err)
		} else {
			fmt.Println("graph.png successfully generated")
			kafka.PublishAlgoGraph(cfg.KafkaBroker, cfg.KafkaTopic, "graph.png")
		}
	}

	start := time.Now()
	fmt.Println("Программа запущена")
	resultItems, brokenItems := logic.Process(operations, aliveVars)

	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %s\n", elapsed)

	resp := &gen.OperationResponse{
		Items: resultItems,
	}

	entry := formLogEntry(req, resp)

	if blm.GRPCClient == nil {
		fmt.Println("GRPCClient is nil — log clients not initialized... Proceeding without it")
	}
	if blm.GRPCClient != nil && !isNil(blm.GRPCClient.LoggerClient) {
		responseLog, err := blm.GRPCClient.LogDataGRPC(ctx, entry)
		if err != nil {
			fmt.Println("Something went wrong during logging: ", err.Error())
		}
		fmt.Println("Log of the result: " + responseLog.GetId())

		resp.LogID = responseLog

	}

	resp.ProcessingTime = durationpb.New(elapsed)

	if len(brokenItems) != 0 {
		warningText := fmt.Sprintf("WARNING: variable(s) %s called for print before calculation", strings.Join(brokenItems, ", "))
		resp.Warning = &warningText
	}

	fmt.Println(resp.GetWarning())

	return resp, nil

}

func formLogEntry(req *gen.OperationRequest, opsResp *gen.OperationResponse) *gen.LogEntry {

	return &gen.LogEntry{
		ServiceName: "business-server",
		Level:       "INFO",
		Message: &gen.StructuredMessage{
			Path:   req.GetLogID().GetId(),
			Result: opsResp,
		},
		Metadata:      nil,
		TimestampSend: 0,
	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}
