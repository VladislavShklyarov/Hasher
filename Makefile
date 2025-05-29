PROJECT_ROOT := .
PROTO_DIR := proto/
PROTO_FILE := gen.proto
OUTPUT = cover_total.out

SERVICES := log-service http-service business-service dashboard-service

.PHONY: generate clean test coverage

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


.PHONY: test coverage

test:
	@for service in $(SERVICES); do \
		echo "Running tests with coverage in $$service..."; \
		cd $$service && \
		pkgs=$$(go list ./... | grep -v '/gen'); \
		if [ -n "$$pkgs" ]; then \
			go test -v -cover -coverprofile=cover.out $$pkgs; \
		else \
			echo "No packages to test in $$service (after excluding gen)"; \
		fi; \
		cd - > /dev/null; \
	done

coverage: test
	@echo "Combining coverage profiles into $(OUTPUT)..."
	@rm -f $(OUTPUT)
	@first=1; \
	for service in $(SERVICES); do \
		if [ -f "$$service/cover.out" ]; then \
			if [ $$first -eq 1 ]; then \
				head -n 1 "$$service/cover.out" > $(OUTPUT); \
				tail -n +2 "$$service/cover.out" | sed 's|^|./|' >> $(OUTPUT); \
				first=0; \
			else \
				tail -n +2 "$$service/cover.out" | sed 's|^|./|' >> $(OUTPUT); \
			fi; \
		else \
			echo "Warning: $$service/cover.out not found!"; \
		fi; \
	done
	@echo "Combined coverage profile created at $(OUTPUT)"
	@echo
	@echo "=== Function coverage summary ==="
	@go tool cover -func=$(OUTPUT)
	@echo
	@echo "Generating HTML report as coverage.html..."
	@go tool cover -html=$(OUTPUT) -o coverage.html
	@echo "HTML report saved to coverage.html"