.PHONY: test

test:
	go test -v -cover -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o test_coverage.html
	@echo "HTML coverage report generated: test_coverage.html"
