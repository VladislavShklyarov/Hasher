package server

import (
	"business-service/gen"
	logGRPC "business-service/internal/client/grpc/log"
	"business-service/internal/logic"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type BusinessLogicManager struct {
	gen.UnimplementedBusinessLogicServer
	GRPCClient *logGRPC.LogClient
}

func (blm *BusinessLogicManager) Process(ctx context.Context, req *gen.OperationRequest) (*gen.OperationResponse, error) {
	log.Printf("Revieved request: %+v", req)

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
		}
	}

	start := time.Now()
	fmt.Println("Программа запущена")
	resultItems, brokenItems := logic.Process(operations, aliveVars, graph)

	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %s\n", elapsed)

	resp := &gen.OperationResponse{
		Items: resultItems,
	}

	entry := formLogEntry(req, resp)

	fmt.Println("Вход на 44 строке:" + entry.String())

	responseLog, err := blm.GRPCClient.LogDataGRPC(ctx, entry)
	fmt.Println("Log of the result: " + responseLog.GetId())

	resp.LogID = responseLog

	if err != nil {
		fmt.Println("Something went wrong during logging: ", err.Error())
	}
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
