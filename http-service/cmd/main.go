package main

import (
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
	"http-service/internal/app"
	grpcBiz "http-service/internal/client/grpc/business"
	grpcLog "http-service/internal/client/grpc/log"
	"http-service/internal/config"
	"http-service/internal/server"
	"http-service/internal/signals"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	clients := &app.Clients{
		LogClient:      grpcLog.CreateLogClient(cfg),
		BusinessClient: grpcBiz.CreateBusinessClient(cfg),
	}

	go server.RunHttpServer(clients, cfg)

	signals.WaitForShutdown(ctx, cancel)
}

// ProcessDataSwagger godoc
// @Summary      Обработка бизнес-операций
// @Description  Принимает JSON с последовательностью операций (`calc`, `print`), преобразует во внутренние Protobuf-сообщения и передаёт в бизнес-сервис и лог-сервис по gRPC.
//
//	Поддерживаются операции с числовыми значениями и ссылками на ранее сохранённые переменные.
//	Пример:
//	{
//	  "operations": [
//	    { "type": "calc", "op": "+", "var": "x", "left": 10, "right": 5 },
//	    { "type": "calc", "op": "*", "var": "y", "left": "x", "right": 3 },
//	    { "type": "print", "var": "y" }
//	  ]
//	}
//
// @Tags         operations
// @Accept       json
// @Produce      json
// @Param request body requestJSON true "Список операций. Поля left и right могут быть числом или строкой (переменной)."
// @Success      200 {object} CompositeResponse "Операции успешно обработаны"
// @Failure      400 {object} CompositeResponse "Некорректный запрос (например, отсутствует поле или неверный формат)"
// @Failure      500 {object} CompositeResponse "Внутренняя ошибка сервера при обработке запроса"
// @Failure      503 {object} CompositeResponse "gRPC-сервисы недоступны"
// @Router       /process [post]
func ProcessDataSwagger() {}

type CompositeResponse struct {
	Success            bool            `json:"success"`
	Status             int             `json:"status"`
	Message            string          `json:"message"`
	LogID              string          `json:"log_id,omitempty"`
	ResultID           string          `json:"result_id,omitempty"`
	LogError           string          `json:"log_error,omitempty"`
	ProcessError       string          `json:"process_error,omitempty"`
	Items              []VariableValue `json:"items,omitempty"`
	ProcessingDuration string          `json:"processing_duration"`
}

type VariableValue struct {
	state         protoimpl.MessageState `protogen:"open. v1"`
	Var           string                 `protobuf:"bytes,1,opt,name=var,proto3" json:"var,omitempty"`
	Value         int64                  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

type requestJSON struct {
	Operations []operationJSON `json:"operations"`
}

type operationJSON struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"` // у операций print может отсутствовать op, сделаем опциональным
	Var   string      `json:"var"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

// DeleteLogSwagger godoc
// @Summary      Удалить лог по идентификатору и имени файла
// @Description  Выполняет gRPC-запрос к лог-сервису для удаления лог-сообщения по указанным параметрам `id` и `filename`.
//
//	В случае успеха возвращает сообщение об успешном удалении. В противном случае возвращает описание ошибки.
//
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        id       query     string  true  "Уникальный идентификатор лога (например, rozNzBFDWy)"
// @Param        filename query     string  true  "Имя файла, в котором содержится лог (например, http_logs.json)"
// @Success      200 {object} DeleteResponse "Успешное удаление лога"
// @Failure      400 {object} DeleteResponse "Ошибка валидации: отсутствуют обязательные параметры id или filename"
// @Failure      500 {object} DeleteResponse "Внутренняя ошибка при выполнении gRPC-запроса или лог не найден"
// @Router       /deleteLog [delete]
func DeleteLogSwagger() {}

// ReadLogSwagger godoc
// @Summary      Получить лог по идентификатору и имени файла
// @Description  Обрабатывает HTTP GET-запрос и выполняет gRPC-вызов к лог-сервису для получения структурированного лог-сообщения.
//
//	Требует обязательные query-параметры `id` и `filename`. В случае успеха возвращает JSON-представление лог-сообщения,
//	иначе возвращает описание ошибки. Обрабатываются следующие ошибки: отсутствие параметров, ошибка gRPC, отсутствие лога.
//
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        id       query     string  true  "Уникальный идентификатор лога (например, rozNzBFDWy)"
// @Param        filename query     string  true  "Имя файла, из которого необходимо извлечь лог (например, http_logs.json)"
// @Success      200 {object} ReadResponse "Успешное чтение лога. Поле 'log' содержит структурированное сообщение."
// @Failure      400 {object} ReadResponse "Ошибка валидации: отсутствует один или оба обязательных параметра (id, filename)"
// @Failure      500 {object} ReadResponse "Внутренняя ошибка: сбой gRPC-запроса или лог/файл не найден на стороне сервиса"
// @Router       /getLog [get]
func ReadLogSwagger() {}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ReadResponse struct {
	Success bool     `json:"success"`
	Log     LogEntry `json:"log,omitempty"`
	Error   string   `json:"error,omitempty"`
}

type LogEntry struct {
	state         protoimpl.MessageState `protogen:"open. v1"`
	ServiceName   string                 `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Level         string                 `protobuf:"bytes,2,opt,name=level,proto3" json:"level,omitempty"`
	Message       *StructuredMessage     `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Metadata      map[string]string      `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	TimestampSend int64                  `protobuf:"varint,5,opt,name=timestamp_send,json=timestampSend,proto3" json:"timestamp_send,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

type StructuredMessage struct {
	state         protoimpl.MessageState `protogen:"open. v1"`
	Method        string                 `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Path          string                 `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Body          []*operationJSON       `protobuf:"bytes,3,rep,name=body,proto3" json:"body,omitempty"`
	Result        OperationResponse      `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

type OperationResponse struct {
	state          protoimpl.MessageState `protogen:"open. v1"`
	LogID          LogID                  `protobuf:"bytes,1,opt,name=LogID,proto3" json:"LogID,omitempty"`
	Items          []*VariableValue       `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	Warning        *string                `protobuf:"bytes,3,opt,name=warning,proto3,oneof" json:"warning,omitempty"`
	ProcessingTime Duration               `protobuf:"bytes,4,opt,name=processing_time,json=processingTime,proto3" json:"processing_time,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

type LogID struct {
	state         protoimpl.MessageState `protogen:"open. v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

type Duration struct {
	state         protoimpl.MessageState `protogen:"open. v1"`
	Seconds       int64                  `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	Nanos         int32                  `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}
