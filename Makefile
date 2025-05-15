PROJECT_ROOT := .
PROTO_LOG_DIR := proto/logger
PROTO_LOG_FILE := logger.proto

PROTO_BUSINESS_DIR := proto/business
PROTO_BUSINESS_FILE := business.proto


SERVICES := log-service http-service business-service

.PHONY: generate clean test

generate:
	@echo "gRPC-script is generating..."
	@for service in $(SERVICES); do \
		mkdir -p $$service/gen; \
		mkdir -p $$service/gen/logger; \
		mkdir -p $$service/gen/business; \
		protoc --go_out=./$$service/gen/logger \
		       --go-grpc_out=./$$service/gen/logger \
		       --proto_path=$(PROTO_LOG_DIR) \
		       $(PROTO_LOG_DIR)/$(PROTO_LOG_FILE); \
		protoc --go_out=./$$service/gen/business \
		       --go-grpc_out=./$$service/gen/business \
		       --proto_path=$(PROTO_BUSINESS_DIR) \
		       $(PROTO_BUSINESS_DIR)/$(PROTO_BUSINESS_FILE); \
		echo "Generated for $$service"; \
	done
	@echo "Generation completed."



clean:
	@echo "Clean generated files..."
	@for service in $(SERVICES); do \
		rm -rf $$service/gen; \
		echo "Cleaned for $$service"; \
	done
