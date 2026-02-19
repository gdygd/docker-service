# REST API 엔드포인트 설계

**Base URL:** `https://api.container-agent.com/api/v1`

---

## 1. Auth API

| Method | Endpoint | 설명 | 인증 |
|--------|----------|------|------|
| POST | `/auth/signup` | 회원가입 | - |
| POST | `/auth/login` | 로그인 | - |
| POST | `/auth/logout` | 로그아웃 | JWT |
| POST | `/auth/refresh` | 토큰 갱신 | Refresh Token |

### POST /auth/signup
```json
// Request
{
  "email": "user@example.com",
  "password": "password123",
  "name": "홍길동"
}

// Response 201
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "홍길동",
    "plan": "free",
    "created_at": "2026-02-04T10:00:00Z"
  }
}
```

### POST /auth/login
```json
// Request
{
  "email": "user@example.com",
  "password": "password123"
}

// Response 200
{
  "success": true,
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "access_token_expires_at": "2026-02-04T11:00:00Z",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "홍길동",
      "plan": "free"
    }
  }
}
```

### POST /auth/logout
```json
// Request (Header: Authorization: Bearer <access_token>)
{
  "refresh_token": "eyJhbG..."
}

// Response 200
{
  "success": true,
  "data": "logged out successfully"
}
```

### POST /auth/refresh
```json
// Request
{
  "refresh_token": "eyJhbG..."
}

// Response 200
{
  "success": true,
  "data": {
    "access_token": "eyJhbG...",
    "access_token_expires_at": "2026-02-04T12:00:00Z"
  }
}
```

---

## 2. Agent 관리 API

| Method | Endpoint | 설명 | 인증 |
|--------|----------|------|------|
| POST | `/agents` | Agent 등록 (agent_key 발급) | JWT |
| GET | `/agents` | 내 Agent 목록 조회 | JWT |
| GET | `/agents/:slug` | Agent 상세 조회 | JWT |
| PUT | `/agents/:slug` | Agent 정보 수정 | JWT |
| DELETE | `/agents/:slug` | Agent 삭제 | JWT |
| POST | `/agents/:slug/regenerate-key` | agent_key 재발급 | JWT |

### POST /agents
```json
// Request
{
  "agent_name": "119서버 에이전트",
  "slug": "119agent",
  "description": "119번 서버 컨테이너 모니터링"
}

// Response 201
{
  "success": true,
  "data": {
    "id": "uuid",
    "agent_name": "119서버 에이전트",
    "slug": "119agent",
    "agent_key": "ak_a1b2c3d4e5f6g7h8i9j0...",
    "status": "inactive",
    "created_at": "2026-02-04T10:00:00Z"
  }
}
```

### GET /agents
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "agent_name": "119서버 에이전트",
      "slug": "119agent",
      "status": "active",
      "last_seen": "2026-02-04T10:30:00Z",
      "host_count": 2,
      "container_count": 15
    },
    {
      "id": "uuid",
      "agent_name": "120서버 에이전트",
      "slug": "120agent",
      "status": "disconnected",
      "last_seen": "2026-02-04T09:00:00Z",
      "host_count": 1,
      "container_count": 8
    }
  ]
}
```

### GET /agents/:slug
```json
// Response 200
{
  "success": true,
  "data": {
    "id": "uuid",
    "agent_name": "119서버 에이전트",
    "slug": "119agent",
    "description": "119번 서버 컨테이너 모니터링",
    "status": "active",
    "last_seen": "2026-02-04T10:30:00Z",
    "hosts": [
      {
        "name": "119server",
        "addr": "tcp://10.1.0.119:2376",
        "status": "active"
      },
      {
        "name": "dev-server",
        "addr": "tcp://10.1.0.120:2376",
        "status": "active"
      }
    ],
    "created_at": "2026-02-04T10:00:00Z"
  }
}
```

### PUT /agents/:slug
```json
// Request
{
  "agent_name": "119서버 에이전트 (업데이트)",
  "description": "설명 변경"
}

// Response 200
{
  "success": true,
  "data": {
    "id": "uuid",
    "agent_name": "119서버 에이전트 (업데이트)",
    "slug": "119agent",
    "description": "설명 변경"
  }
}
```

### POST /agents/:slug/regenerate-key
```json
// Response 200
{
  "success": true,
  "data": {
    "agent_key": "ak_new_key_xxxxxxxxxxxxxxxx...",
    "message": "기존 키는 즉시 무효화됩니다. Agent 설정을 업데이트하세요."
  }
}
```

---

## 3. 모니터링 데이터 API

| Method | Endpoint | 설명 | 인증 |
|--------|----------|------|------|
| GET | `/agents/:slug/hosts` | 호스트 목록 | JWT |
| GET | `/agents/:slug/containers` | 컨테이너 목록 (최신) | JWT |
| GET | `/agents/:slug/containers/:id` | 컨테이너 상세 (Inspect) | JWT |
| GET | `/agents/:slug/stats` | 전체 Stats (최신) | JWT |
| GET | `/agents/:slug/stats/:id` | 개별 컨테이너 Stats | JWT |
| GET | `/agents/:slug/stats/history` | Stats 히스토리 (그래프용) | JWT |
| GET | `/agents/:slug/events` | 이벤트 목록 | JWT |
| GET | `/agents/:slug/events/stream` | SSE 실시간 이벤트 | JWT |

### GET /agents/:slug/hosts
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "name": "119server",
      "addr": "tcp://10.1.0.119:2376",
      "status": "active",
      "container_count": 10
    }
  ]
}
```

### GET /agents/:slug/containers
```
Query Parameters:
  - host: 호스트 필터 (선택)
  - state: 상태 필터 (running, exited, all) (기본: all)
```
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": "a1b2c3d4e5f6",
      "name": "nginx-web",
      "image": "nginx:latest",
      "state": "running",
      "status": "Up 2 hours",
      "host": "119server"
    }
  ]
}
```

### GET /agents/:slug/containers/:id
```json
// Response 200
{
  "success": true,
  "data": {
    "id": "a1b2c3d4e5f6...",
    "name": "nginx-web",
    "image": "nginx:latest",
    "created": "2026-01-15T10:30:00Z",
    "platform": "linux",
    "host": "119server",
    "state": {
      "status": "running",
      "running": true,
      "exit_code": 0,
      "started_at": "2026-01-15T10:30:05Z"
    },
    "config": {
      "hostname": "a1b2c3d4e5f6",
      "env": ["PATH=/usr/local/bin:/usr/bin"],
      "cmd": ["nginx", "-g", "daemon off;"]
    },
    "network": {
      "ip_address": "172.17.0.2",
      "gateway": "172.17.0.1",
      "ports": {
        "80/tcp": [{"host_ip": "0.0.0.0", "host_port": "8080"}]
      }
    },
    "mounts": [
      {
        "type": "bind",
        "source": "/host/html",
        "destination": "/usr/share/nginx/html",
        "mode": "rw",
        "rw": true
      }
    ]
  }
}
```

### GET /agents/:slug/stats
```
Query Parameters:
  - host: 호스트 필터 (선택)
```
```json
// Response 200
{
  "success": true,
  "data": {
    "a1b2c3d4e5f6": {
      "id": "a1b2c3d4e5f6",
      "name": "nginx-web",
      "host": "119server",
      "cpu_percent": 2.35,
      "memory_usage": 134742016,
      "memory_limit": 4294967296,
      "memory_percent": 3.14,
      "network_rx": 1310720,
      "network_tx": 524288,
      "collected_at": "2026-02-04T10:30:00Z"
    }
  }
}
```

### GET /agents/:slug/stats/history
```
Query Parameters:
  - container_id: 컨테이너 ID (필수)
  - host: 호스트명 (필수)
  - from: 시작 시간 (ISO 8601) (기본: 1시간 전)
  - to: 종료 시간 (ISO 8601) (기본: 현재)
  - interval: 집계 간격 (1m, 5m, 15m, 1h) (기본: 1m)
```
```json
// Response 200
{
  "success": true,
  "data": {
    "container_id": "a1b2c3d4e5f6",
    "container_name": "nginx-web",
    "interval": "1m",
    "points": [
      {
        "time": "2026-02-04T10:00:00Z",
        "cpu_percent": 2.35,
        "memory_usage": 134742016,
        "memory_percent": 3.14,
        "network_rx": 1310720,
        "network_tx": 524288
      },
      {
        "time": "2026-02-04T10:01:00Z",
        "cpu_percent": 1.82,
        "memory_usage": 135200000,
        "memory_percent": 3.15,
        "network_rx": 1320000,
        "network_tx": 530000
      }
    ]
  }
}
```

### GET /agents/:slug/events
```
Query Parameters:
  - host: 호스트 필터 (선택)
  - type: 이벤트 타입 필터 (container, network, image) (선택)
  - from: 시작 시간 (선택)
  - limit: 개수 제한 (기본: 50, 최대: 200)
```
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "host": "119server",
      "type": "container",
      "action": "start",
      "actor_id": "a1b2c3d4e5f6",
      "actor_name": "nginx-web",
      "timestamp": 1707040200,
      "attrs": {
        "image": "nginx:latest",
        "name": "nginx-web"
      }
    }
  ]
}
```

### GET /agents/:slug/events/stream (SSE)
```
Headers:
  Accept: text/event-stream
  Authorization: Bearer <access_token>

Response:
  Content-Type: text/event-stream
```
```
event: container-event
data: {"host":"119server","type":"container","action":"start","actor_id":"a1b2c3d4e5f6","actor_name":"nginx-web","timestamp":1707040200}

event: container-event
data: {"host":"119server","type":"container","action":"die","actor_id":"f6e5d4c3b2a1","actor_name":"redis-cache","timestamp":1707040260}
```

---

## 4. 원격 제어 API (향후)

| Method | Endpoint | 설명 | 인증 |
|--------|----------|------|------|
| POST | `/agents/:slug/containers/:id/start` | 컨테이너 시작 | JWT |
| POST | `/agents/:slug/containers/:id/stop` | 컨테이너 중지 | JWT |
| POST | `/agents/:slug/containers/:id/restart` | 컨테이너 재시작 | JWT |

```json
// POST /agents/:slug/containers/:id/start
// Request
{
  "host": "119server"
}

// Response 200
{
  "success": true,
  "data": "container start command sent"
}
```

---

## 5. 공통 에러 응답

```json
// 400 Bad Request
{
  "success": false,
  "message": "invalid request body"
}

// 401 Unauthorized
{
  "success": false,
  "message": "invalid or expired token"
}

// 403 Forbidden
{
  "success": false,
  "message": "access denied to this agent"
}

// 404 Not Found
{
  "success": false,
  "message": "agent not found"
}

// 429 Too Many Requests
{
  "success": false,
  "message": "rate limit exceeded"
}

// 500 Internal Server Error
{
  "success": false,
  "message": "internal server error"
}
```

---

## 6. 인증 헤더

모든 인증이 필요한 API는 아래 헤더를 포함해야 합니다:

```
Authorization: Bearer <access_token>
```

---

## 7. Plan별 제한 (향후)

| 항목 | Free | Pro | Enterprise |
|------|------|-----|------------|
| Agent 수 | 1 | 5 | 무제한 |
| 데이터 보관 | 7일 | 30일 | 90일 |
| Stats 수집 주기 | 60초 | 30초 | 10초 |
| 원격 제어 | X | O | O |
| API Rate Limit | 100/분 | 1000/분 | 10000/분 |
