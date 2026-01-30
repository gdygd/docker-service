GENPROTO_DIR := $(shell go list -f '{{ .Dir }}' -m google.golang.org/genproto)
PROTOBUF_DIR := $(shell go list -f '{{ .Dir }}' -m google.golang.org/protobuf)
GATEWAY_DIR  := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2)

OUT_DIR := pb

# 테스트 실행
test:
	go test -v -cover 

start:
	cd cmd && go run main.go
	
# 빌드 실행
build:
	go build -o ./bin/auth-service ./cmd/main.go

# proto 컴파일 (grpc + gateway + annotations)
proto:
	@echo "Generating Go files from proto..."
	rm -f $(OUT_DIR)/*.go
	protoc \
		-I$(PROTOBUF_DIR) \
		-I$(GENPROTO_DIR) \
		-I$(GATEWAY_DIR) \
		--proto_path=proto \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		proto/*.proto

.PHONY: build proto test
