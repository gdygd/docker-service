test:
	go test -v -cover 

build:
	go build -o ./bin/delivery-service ./cmd/main.go

.PHONY: build