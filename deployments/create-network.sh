#!/bin/bash
# ===========================================
# Docker Network 생성 스크립트
# ===========================================

NETWORK_NAME="docker-mng"

# 네트워크 존재 여부 확인
if docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
    echo "Network '$NETWORK_NAME' already exists."
else
    echo "Creating network '$NETWORK_NAME'..."
    docker network create --driver bridge "$NETWORK_NAME"
    echo "Network '$NETWORK_NAME' created successfully."
fi
