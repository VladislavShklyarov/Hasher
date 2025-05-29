PROJECT_ROOT := .
PROTO_DIR := proto/
PROTO_FILE := gen.proto



SERVICES := log-service http-service business-service dashboard-service

.PHONY: generate clean test

generate:
	@echo "gRPC-script is generating..."
	@for service in $(SERVICES); do \
		mkdir -p $$service/gen; \
		protoc --go_out=./$$service/gen/ \
		       --go-grpc_out=./$$service/gen/ \
		       --proto_path=$(PROTO_DIR) \
		       $(PROTO_DIR)/$(PROTO_FILE); \
		echo "Generated for $$service"; \
	done
	@echo "Generation completed."



clean:
	@echo "Clean generated files..."
	@for service in $(SERVICES); do \
		rm -rf $$service/gen; \
		echo "Cleaned for $$service"; \
	done
