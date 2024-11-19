CLI_NAME = "counters"
SERVER_NAME = "counters-server"

schema-docs:
	@echo "Generating schema documentation..."
	go run utils/gen_schema.go
	@echo "Schema documentation generated."

cli:
	@echo "Building CLI..."
	go build -o ~/go/bin/$(CLI_NAME) ./cmd
	@echo "CLI built."

build-server:
	@echo "Building Server..."
	go build -o ~/go/bin/$(SERVER_NAME) ./server
	@echo "Server built."

.PHONY: server
server:
	@echo "Running Server..."
	go run ./server
