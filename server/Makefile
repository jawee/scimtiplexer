all: build test

GOOS := $(shell go env GOOS)

build:
	@echo "Building..."

	@CGO_ENABLED=1 GOOS=$(GOOS) go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go 

run-api:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

test-coverage:
	rm -r testresults || true
	mkdir testresults
	go test -coverprofile ./testresults/cover.out -o ./testresults ./...
	go tool cover -html=./testresults/cover.out -o ./testresults/cover.html
	echo "Coverage report is saved in testresults/cover.html"

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch

