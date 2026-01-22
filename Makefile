
IMAGE_NAME=docker-service

test:
	go test -v -cover 

start:	
	cd cmd && go run main.go
build:
	go build -o ./bin/docker-service ./cmd/main.go


run:
	docker run --rm \		
		-p 9083:9083 \
		$(IMAGE_NAME):latest

.PHONY: build