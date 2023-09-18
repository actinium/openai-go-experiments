################################################################################
# Build
################################################################################
.PHONY: build
build:
	@go build -o chat cmd/cli/main.go
	@go build -o streaming-chat cmd/streaming_cli/main.go
	@go build -o imagine cmd/image_cli/main.go


################################################################################
# Clean
################################################################################
.PHONY: clean
clean:
	@rm -rf chat streaming-chat imagine
