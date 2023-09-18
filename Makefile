################################################################################
# Build
################################################################################
.PHONY: build
build:
	@go build -o chat cmd/chat_cli/main.go
	@go build -o streaming-chat cmd/chat_streaming_cli/main.go
	@go build -o imagine cmd/image_cli/main.go
	@go build -o ai2ai cmd/ai2ai_cli/main.go
	@go build -o create-embedding cmd/embeddings_cli/main.go

################################################################################
# Run
################################################################################
.PHONY: run-color-server
run-color-server:
	@go run cmd/color_web/main.go

.PHONY: run-imagine-server
run-imagine-server:
	@go run cmd/image_web/main.go

################################################################################
# Clean
################################################################################
.PHONY: clean
clean:
	@rm -rf chat streaming-chat imagine ai2ai create-embedding
