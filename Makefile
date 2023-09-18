################################################################################
# Build
################################################################################
.PHONY: build
build:
	@go build -o chat cmd/chat_cli/main.go
	@go build -o streaming-chat cmd/chat_streaming_cli/main.go
	@go build -o imagine cmd/image_cli/main.go


################################################################################
# Clean
################################################################################
.PHONY: clean
clean:
	@rm -rf chat streaming-chat imagine
