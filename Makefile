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
	@go build -o say cmd/text_to_speech_cli/main.go
	@go build -o list-models cmd/models_cli/main.go

.PHONY: generate
generate:
	@go generate ./...

################################################################################
# Run
################################################################################
.PHONY: run-chat-server
run-chat-server:
	@go run cmd/chat_web/main.go

.PHONY: run-translation-server
run-translation-server:
	@go run cmd/translation_web/main.go

.PHONY: run-color-server
run-color-server:
	@go run cmd/color_web/main.go

.PHONY: run-image-server
run-image-server: run-image-server-dalle-3

.PHONY: run-image-server-dalle-2
run-image-server-dalle-2:
	@go run cmd/image_web/main.go

.PHONY: run-image-server-dalle-3
run-image-server-dalle-3:
	@go run cmd/image_web_dalle_3/main.go

.PHONY: run-tts-server
run-tts-server:
	@go run cmd/text_to_speech_web/main.go

################################################################################
# Clean
################################################################################
.PHONY: clean
clean:
	@rm -rf chat streaming-chat imagine ai2ai create-embedding say
