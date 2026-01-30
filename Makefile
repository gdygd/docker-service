
IMAGE_NAME=docker-service
APIGW_BIN_DIR = ./services/api-gateway/bin
AUTH_BIN_DIR = ./services/auth_service/bin
API_BIN_DIR = ./services/docker_service/bin

SERVICE_PATH = ./services/docker_service
AUTH_SERVICE_PATH = ./services/auth_service
APIGW_SERVICE_PATH = ./services/api-gateway

# ---------------------------------
# Build
# ---------------------------------
build-all: build-api build-auth build-gw

build-api:
	cd $(SERVICE_PATH) && go build -o ./bin/docker-service ./cmd/main.go

build-auth:
	cd $(AUTH_SERVICE_PATH) && go build -o ./bin/auth-service ./cmd/main.go

build-gw:
	cd $(APIGW_SERVICE_PATH) && go build -o ./bin/api-gateway ./cmd/main.go

# ---------------------------------
# Start (build & run)
# ---------------------------------
startgw: build-gw
	cd $(APIGW_SERVICE_PATH)/bin && ./api-gateway

startauth: build-auth
	cd $(AUTH_BIN_DIR) && ./auth-service

startapi: build-api
	cd $(API_BIN_DIR) && ./docker-service

allstart: build-all startgw startauth startapi
# 	@echo "Starting all services..."
# 	$(APIGW_SERVICE_PATH)/api-gateway &
# 	$(AUTH_BIN_DIR)/auth-service &
# 	$(API_BIN_DIR)/docker-service &
# 	@echo "All services started in background"

allstop:
	@echo "Stopping all services..."
	@pkill -f "./api-gateway" || true
	@pkill -f "./auth-service" || true
	@pkill -f "./docker-service" || true
	@echo "All services stopped"

# 포트 기반 종료 (더 확실함)
allstop-port:
	@echo "Stopping all services by port..."	
	@fuser -k 9091/tcp 2>/dev/null || true
	@fuser -k 9190/tcp 2>/dev/null || true
	@fuser -k 9081/tcp 2>/dev/null || true
	@fuser -k 9082/tcp 2>/dev/null || true
	@fuser -k 9083/tcp 2>/dev/null || true
	@echo "All services stopped"

# ---------------------------------
# Test
# ---------------------------------
test:
	go test -v -cover ./...

# ---------------------------------
# Docker
# ---------------------------------
run:
	docker run --rm \
		-p 9083:9083 \
		$(IMAGE_NAME):latest

# ---------------------------------
# Clean
# ---------------------------------
clean:
	rm -rf $(BIN_DIR)

.PHONY: build-all build-api build-auth build-gw \
        startgw startauth startapi allstart allstop \
        test run clean