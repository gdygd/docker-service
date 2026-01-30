# =================================
# docker_service Makefile
# =================================

SERVICE_NAME = docker-service
IMAGE_NAME = docker-service
REGISTRY_URL ?= 10.1.0.119:5000

# ---------------------------------
# Development
# ---------------------------------
.PHONY: start test build clean

start:
	cd cmd && go run main.go

test:
	go test -v -cover ./...

build:	
	go build -o ./bin/$(SERVICE_NAME) ./cmd/main.go
	cd ./bin
	##./$(SERVICE_NAME)

clean:
	rm -rf ./bin/

# ---------------------------------
# Docker
# ---------------------------------
.PHONY: docker-build docker-push docker-run

docker-build:
	docker build -t $(IMAGE_NAME):latest \
		-t $(REGISTRY_URL)/$(IMAGE_NAME):latest \
		-f Dockerfile .

docker-push:
	docker push $(REGISTRY_URL)/$(IMAGE_NAME):latest

docker-run:
	docker run --rm \
		-p 9083:9083 \
		-v $(PWD)/certs:/certs:ro \
		$(IMAGE_NAME):latest

# ---------------------------------
# All-in-one
# ---------------------------------
.PHONY: deploy

deploy: docker-build docker-push
	@echo "✅ $(SERVICE_NAME) deployed to $(REGISTRY_URL)"
