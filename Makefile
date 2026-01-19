test:
	go test -v -cover 

start:	
	cd cmd && go run main.go
build:
	go build -o ./bin/docker-service ./cmd/main.go

.PHONY: build