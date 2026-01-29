# Docker Container Management API Documentation

**Base URL:** `http://{server}:{port}`
**Content-Type:** `application/json`

---

## 공통 응답 형식

### 성공 응답
```json
{
  "success": true,
  "data": { ... }
}
```

### 에러 응답
```json
{
  "success": false,
  "message": "에러 메시지"
}
```

---

## Authentication APIs

---

## 1. POST /user

새로운 사용자를 생성합니다.

### Request
```
POST /user
Content-Type: application/json

{
  "username": "gildong",
  "password": "123123",
  "full_name": "YunGilDong",
  "email": "gildong@email.com"
}
```

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `username` | string | Yes | 사용자 ID |
| `password` | string | Yes | 비밀번호 |
| `full_name` | string | Yes | 사용자 이름 |
| `email` | string | Yes | 이메일 주소 |

### Response
```json
{
  "success": true,
  "data": {
    "username": "gildong",
    "full_name": "YunGilDong",
    "email": "gildong@email.com",
    "created_at": "2026-01-28 21:33:51"
  }
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `username` | string | 사용자 ID |
| `full_name` | string | 사용자 이름 |
| `email` | string | 이메일 주소 |
| `created_at` | string | 생성 시간 |

---

## 2. POST /login

사용자 로그인을 수행하고 인증 토큰을 발급합니다.

### Request
```
POST /login
Content-Type: application/json

{
  "username": "gildong",
  "password": "123123"
}
```

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `username` | string | Yes | 사용자 ID |
| `password` | string | Yes | 비밀번호 |

### Response
```json
{
  "session_id": "4b2efe8c-b526-4ff3-a2fa-c1ddcb9601e9",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "access_token_expires_at": "2026-01-28T22:14:11.180108176+09:00",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token_expires_at": "2026-01-29T21:59:11.180211342+09:00",
  "user": {
    "username": "gildong",
    "email": "gildong@email.com",
    "password_changed_at": "2026-01-28T21:33:51Z",
    "created_at": "2026-01-28T21:33:51Z"
  }
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `session_id` | string | 세션 ID (UUID) |
| `access_token` | string | JWT 액세스 토큰 |
| `access_token_expires_at` | string | 액세스 토큰 만료 시간 |
| `refresh_token` | string | JWT 리프레시 토큰 |
| `refresh_token_expires_at` | string | 리프레시 토큰 만료 시간 |
| `user` | object | 사용자 정보 |

#### User Object
| Field | Type | Description |
|-------|------|-------------|
| `username` | string | 사용자 ID |
| `email` | string | 이메일 주소 |
| `password_changed_at` | string | 비밀번호 변경 시간 |
| `created_at` | string | 계정 생성 시간 |

---

## 3. POST /logout

사용자 로그아웃을 수행하고 세션을 삭제합니다.

### Request
```
POST /logout
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | Yes | 로그인 시 발급받은 리프레시 토큰 |

### Response
```json
{
  "success": true,
  "data": "logged out successfully"
}
```

---

## 4. POST /token/renew_access

리프레시 토큰을 사용하여 새로운 액세스 토큰을 발급합니다.

### Request
```
POST /token/renew_access
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | Yes | 로그인 시 발급받은 리프레시 토큰 |

### Response
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "access_token_expires_at": "2026-01-28T22:32:51.341015248+09:00"
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | 새로 발급된 JWT 액세스 토큰 |
| `access_token_expires_at` | string | 액세스 토큰 만료 시간 |

---

## Docker APIs

---

## 5. GET /hosts

등록된 Docker 호스트 목록을 조회합니다.

### Request
```
GET /hosts
```

### Response
```json
{
  "success": true,
  "data": [
    {
      "host": "119server",
      "addr": "tcp://10.1.0.119:2376"
    },
    {
      "host": "dev-server",
      "addr": "tcp://10.1.0.120:2376"
    }
  ]
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `host` | string | 호스트 식별 이름 |
| `addr` | string | Docker Daemon 주소 (tcp:// 또는 unix) |

---

## 6. GET /ps2/:host

특정 호스트의 컨테이너 목록을 조회합니다.

### Request
```
GET /ps2/{host}
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `host` | string | Yes | 호스트 식별 이름 (예: `119server`) |

### Example
```
GET /ps2/119server
```

### Response
```json
{
  "success": true,
  "data": [
    {
      "id": "a1b2c3d4e5f6",
      "name": "nginx-web",
      "image": "nginx:latest",
      "state": "running",
      "status": "Up 2 hours"
    },
    {
      "id": "f6e5d4c3b2a1",
      "name": "redis-cache",
      "image": "redis:7",
      "state": "exited",
      "status": "Exited (0) 1 hour ago"
    }
  ]
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `id` | string | 컨테이너 ID (12자리) |
| `name` | string | 컨테이너 이름 |
| `image` | string | 이미지 이름 |
| `state` | string | 상태 (`running`, `exited`, `paused`, etc.) |
| `status` | string | 상태 설명 |

---

## 7. GET /inspect2/:host/:id

특정 호스트의 컨테이너 상세 정보를 조회합니다.

### Request
```
GET /inspect2/{host}/{id}
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `host` | string | Yes | 호스트 식별 이름 |
| `id` | string | Yes | 컨테이너 ID 또는 이름 |

### Example
```
GET /inspect2/119server/a1b2c3d4e5f6
```

### Response
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4e5f6789...",
    "name": "/nginx-web",
    "image": "sha256:abc123...",
    "created": "2024-01-15T10:30:00.000000000Z",
    "platform": "linux",
    "restart_count": 0,
    "state": {
      "status": "running",
      "running": true,
      "paused": false,
      "restarting": false,
      "exit_code": 0,
      "started_at": "2024-01-15T10:30:05.000000000Z",
      "finished_at": "0001-01-01T00:00:00Z"
    },
    "config": {
      "hostname": "a1b2c3d4e5f6",
      "user": "",
      "env": [
        "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        "NGINX_VERSION=1.25.3"
      ],
      "cmd": ["nginx", "-g", "daemon off;"],
      "entrypoint": ["/docker-entrypoint.sh"],
      "working_dir": "",
      "labels": {
        "maintainer": "NGINX Docker Maintainers"
      }
    },
    "network": {
      "ip_address": "172.17.0.2",
      "gateway": "172.17.0.1",
      "mac_address": "02:42:ac:11:00:02",
      "ports": {
        "80/tcp": [
          {
            "host_ip": "0.0.0.0",
            "host_port": "8080"
          }
        ]
      },
      "networks": {
        "bridge": {
          "network_id": "abc123...",
          "ip_address": "172.17.0.2",
          "gateway": "172.17.0.1",
          "mac_address": "02:42:ac:11:00:02"
        }
      }
    },
    "mounts": [
      {
        "type": "bind",
        "name": "",
        "source": "/host/path/html",
        "destination": "/usr/share/nginx/html",
        "mode": "rw",
        "rw": true
      }
    ]
  }
}
```

### Response Fields

#### Root
| Field | Type | Description |
|-------|------|-------------|
| `id` | string | 컨테이너 전체 ID |
| `name` | string | 컨테이너 이름 |
| `image` | string | 이미지 ID |
| `created` | string | 생성 시간 (ISO 8601) |
| `platform` | string | 플랫폼 (linux/windows) |
| `restart_count` | int | 재시작 횟수 |
| `state` | object | 상태 정보 |
| `config` | object | 설정 정보 |
| `network` | object | 네트워크 정보 |
| `mounts` | array | 마운트 정보 |

#### State
| Field | Type | Description |
|-------|------|-------------|
| `status` | string | 상태 (`running`, `exited`, `paused`) |
| `running` | bool | 실행 중 여부 |
| `paused` | bool | 일시정지 여부 |
| `restarting` | bool | 재시작 중 여부 |
| `exit_code` | int | 종료 코드 |
| `started_at` | string | 시작 시간 |
| `finished_at` | string | 종료 시간 |

#### Config
| Field | Type | Description |
|-------|------|-------------|
| `hostname` | string | 호스트명 |
| `user` | string | 실행 사용자 |
| `env` | array | 환경변수 목록 |
| `cmd` | array | 실행 명령어 |
| `entrypoint` | array | 엔트리포인트 |
| `working_dir` | string | 작업 디렉토리 |
| `labels` | object | 라벨 |

#### Network
| Field | Type | Description |
|-------|------|-------------|
| `ip_address` | string | IP 주소 |
| `gateway` | string | 게이트웨이 |
| `mac_address` | string | MAC 주소 |
| `ports` | object | 포트 바인딩 정보 |
| `networks` | object | 연결된 네트워크 목록 |

#### Mounts
| Field | Type | Description |
|-------|------|-------------|
| `type` | string | 마운트 타입 (`bind`, `volume`, `tmpfs`) |
| `name` | string | 볼륨 이름 (volume인 경우) |
| `source` | string | 호스트 경로 |
| `destination` | string | 컨테이너 경로 |
| `mode` | string | 모드 (`rw`, `ro`) |
| `rw` | bool | 읽기/쓰기 가능 여부 |

---

## 8. POST /start2

특정 호스트의 컨테이너를 시작합니다.

### Request
```
POST /start2
Content-Type: application/json

{
  "id": "a1b2c3d4e5f6",
  "host": "119server"
}
```

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | 컨테이너 ID 또는 이름 |
| `host` | string | Yes | 호스트 식별 이름 |

### Response (Success)
```
HTTP/1.1 200 OK
""
```

### Response (Error)
```json
{
  "success": false,
  "message": "Error response from daemon: container already started"
}
```

---

## 9. POST /stop2

특정 호스트의 컨테이너를 중지합니다.

### Request
```
POST /stop2
Content-Type: application/json

{
  "id": "a1b2c3d4e5f6",
  "host": "119server"
}
```

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | 컨테이너 ID 또는 이름 |
| `host` | string | Yes | 호스트 식별 이름 |

### Response (Success)
```
HTTP/1.1 200 OK
""
```

### Response (Error)
```json
{
  "success": false,
  "message": "Error response from daemon: container not running"
}
```

---

## 10. GET /stat2/:host/:id

특정 호스트의 컨테이너 리소스 사용량을 조회합니다.

### Request
```
GET /stat2/{host}/{id}
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `host` | string | Yes | 호스트 식별 이름 |
| `id` | string | Yes | 컨테이너 ID 또는 이름 |

### Example
```
GET /stat2/119server/a1b2c3d4e5f6
```

### Response
```json
{
  "success": true,
  "data": {
    "id": "",
    "name": "",
    "cpu_percent": 2.35,
    "memory_usage": "128.50 MiB",
    "memory_limit": "4.00 GiB",
    "memory_percent": 3.14,
    "network_rx": "1.25 MB",
    "network_tx": "512.00 KB"
  }
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `id` | string | 컨테이너 ID |
| `name` | string | 컨테이너 이름 |
| `cpu_percent` | float | CPU 사용률 (%) |
| `memory_usage` | string | 메모리 사용량 (포맷팅됨) |
| `memory_limit` | string | 메모리 제한 (포맷팅됨) |
| `memory_percent` | float | 메모리 사용률 (%) |
| `network_rx` | string | 네트워크 수신량 (포맷팅됨) |
| `network_tx` | string | 네트워크 송신량 (포맷팅됨) |

---

## 11. GET /stat3/:host

특정 호스트의 **모든 컨테이너** 리소스 사용량을 일괄 조회합니다.

### Request
```
GET /stat3/{host}
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `host` | string | Yes | 호스트 식별 이름 (예: `119server`) |

### Example
```
GET /stat3/119server
```

### Response
```json
{
  "success": true,
  "data": {
    "58287566f213": {
      "id": "58287566f213",
      "name": "registry-ui-https",
      "cpu_percent": 0,
      "memory_usage": "2.91 MiB",
      "memory_limit": "3.83 GiB",
      "memory_percent": 0.07,
      "network_rx": "401.12 KB",
      "network_tx": "9.64 KB"
    },
    "7adb406b3f36": {
      "id": "7adb406b3f36",
      "name": "docker-service",
      "cpu_percent": 0,
      "memory_usage": "0.00 B",
      "memory_limit": "0.00 B",
      "memory_percent": 0,
      "network_rx": "0 B",
      "network_tx": "0 B"
    },
    "8d603732e1fc": {
      "id": "8d603732e1fc",
      "name": "docker_mariadb",
      "cpu_percent": 0.01,
      "memory_usage": "53.00 MiB",
      "memory_limit": "3.83 GiB",
      "memory_percent": 1.35,
      "network_rx": "325.20 KB",
      "network_tx": "367.22 KB"
    }
  }
}
```

### Response Structure
응답의 `data` 필드는 컨테이너 ID를 키로 하는 객체(Map)입니다.

### Response Fields (각 컨테이너)
| Field | Type | Description |
|-------|------|-------------|
| `id` | string | 컨테이너 ID (12자리) |
| `name` | string | 컨테이너 이름 |
| `cpu_percent` | float | CPU 사용률 (%) |
| `memory_usage` | string | 메모리 사용량 (포맷팅됨) |
| `memory_limit` | string | 메모리 제한 (포맷팅됨) |
| `memory_percent` | float | 메모리 사용률 (%) |
| `network_rx` | string | 네트워크 수신량 (포맷팅됨) |
| `network_tx` | string | 네트워크 송신량 (포맷팅됨) |

### Notes
- 실행 중이지 않은 컨테이너는 `memory_usage`, `memory_limit`가 `"0.00 B"`로 표시됩니다
- 3초 timeout이 적용되어 있으며, timeout 발생 시 수집된 결과까지만 반환됩니다

---

## 12. GET /events (SSE)

Docker 컨테이너 이벤트를 실시간으로 수신하는 Server-Sent Events (SSE) 엔드포인트입니다.

### Request
```
GET /events
Accept: text/event-stream
```

### Response
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

### Event Format
```
event: container-event
data: {"host":"...","type":"...","action":"...","actor_id":"...","actor_name":"...","timestamp":...,"attrs":{...}}
```

### Event Fields
| Field | Type | Description |
|-------|------|-------------|
| `host` | string | Docker 호스트 이름 |
| `type` | string | 이벤트 타입 (`container`, `network`, `image`, `volume`, `daemon`) |
| `action` | string | 이벤트 액션 |
| `actor_id` | string | 대상 ID (컨테이너 ID, 네트워크 ID 등) |
| `actor_name` | string | 대상 이름 |
| `timestamp` | int | Unix timestamp (초) |
| `attrs` | object | 추가 속성 정보 |

### Event Types & Actions

#### Container Events
| Action | Description |
|--------|-------------|
| `create` | 컨테이너 생성 |
| `start` | 컨테이너 시작 |
| `stop` | 컨테이너 중지 |
| `restart` | 컨테이너 재시작 |
| `die` | 컨테이너 종료 |
| `kill` | 컨테이너 강제 종료 |
| `pause` | 컨테이너 일시정지 |
| `unpause` | 컨테이너 일시정지 해제 |
| `destroy` | 컨테이너 삭제 |

#### Network Events
| Action | Description |
|--------|-------------|
| `create` | 네트워크 생성 |
| `connect` | 컨테이너가 네트워크에 연결 |
| `disconnect` | 컨테이너가 네트워크에서 분리 |
| `destroy` | 네트워크 삭제 |

#### Image Events
| Action | Description |
|--------|-------------|
| `pull` | 이미지 풀 |
| `push` | 이미지 푸시 |
| `tag` | 이미지 태그 |
| `delete` | 이미지 삭제 |

### Attributes (attrs)
| Attribute | Description |
|-----------|-------------|
| `name` | 컨테이너/네트워크 이름 |
| `image` | 이미지 이름 |
| `exitCode` | 종료 코드 (die 이벤트) |
| `execDuration` | 실행 시간 (ms) |
| `signal` | kill 시그널 |
| `container` | 컨테이너 ID (네트워크 이벤트) |
| `com.docker.compose.project` | Docker Compose 프로젝트명 |
| `com.docker.compose.service` | Docker Compose 서비스명 |

### Example Events

#### Container Die Event
```
event: container-event
data: {"host":"119server","type":"container","action":"die","actor_id":"4a55acf11f30b628d58ce47a19bde7d30e144314b530f90a2501524e7f3f9cd4","actor_name":"final_project-redis-1","timestamp":1769573841,"attrs":{"com.docker.compose.project":"final_project","com.docker.compose.service":"redis","execDuration":"1001","exitCode":"0","image":"redis:alpine","name":"final_project-redis-1"}}
```

#### Network Connect Event
```
event: container-event
data: {"host":"119server","type":"network","action":"connect","actor_id":"486c597a3abbd36a9d78d8003f1e5366c2606689e0d3db99148beaf52ce7510e","actor_name":"final_project_default","timestamp":1769573852,"attrs":{"container":"4a55acf11f30b628d58ce47a19bde7d30e144314b530f90a2501524e7f3f9cd4","name":"final_project_default"}}
```

#### Container Start Event
```
event: container-event
data: {"host":"119server","type":"container","action":"start","actor_id":"4a55acf11f30b628d58ce47a19bde7d30e144314b530f90a2501524e7f3f9cd4","actor_name":"final_project-redis-1","timestamp":1769573852,"attrs":{"com.docker.compose.project":"final_project","com.docker.compose.service":"redis","image":"redis:alpine","name":"final_project-redis-1"}}
```

### JavaScript Client Example
```javascript
const eventSource = new EventSource('http://localhost:9083/events');

eventSource.addEventListener('container-event', (event) => {
  const data = JSON.parse(event.data);
  console.log(`[${data.host}] ${data.type}/${data.action}: ${data.actor_name}`);

  // 컨테이너 종료 이벤트 처리
  if (data.type === 'container' && data.action === 'die') {
    console.log(`Exit code: ${data.attrs.exitCode}`);
  }
});

eventSource.onerror = (error) => {
  console.error('SSE connection error:', error);
};

// 연결 종료
// eventSource.close();
```

### cURL Example
```bash
curl -N -H "Accept: text/event-stream" http://localhost:9083/events
```

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | 성공 |
| 400 | 잘못된 요청 (필수 파라미터 누락 등) |
| 500 | 서버 에러 (Docker Daemon 연결 실패 등) |

---

## cURL Examples

```bash
# 사용자 생성
curl -X POST http://localhost:9083/user \
  -H "Content-Type: application/json" \
  -d '{"username":"gildong","password":"123123","full_name":"YunGilDong","email":"gildong@email.com"}'

# 로그인
curl -X POST http://localhost:9083/login \
  -H "Content-Type: application/json" \
  -d '{"username":"gildong","password":"123123"}'

# 로그아웃
curl -X POST http://localhost:9083/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}'

# 액세스 토큰 갱신
curl -X POST http://localhost:9083/token/renew_access \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}'

# 호스트 목록 조회
curl -X GET http://localhost:9083/hosts

# 컨테이너 목록 조회
curl -X GET http://localhost:9083/ps2/119server

# 컨테이너 상세 조회
curl -X GET http://localhost:9083/inspect2/119server/nginx-web

# 컨테이너 시작
curl -X POST http://localhost:9083/start2 \
  -H "Content-Type: application/json" \
  -d '{"id":"nginx-web","host":"119server"}'

# 컨테이너 중지
curl -X POST http://localhost:9083/stop2 \
  -H "Content-Type: application/json" \
  -d '{"id":"nginx-web","host":"119server"}'

# 컨테이너 리소스 사용량 조회 (단일)
curl -X GET http://localhost:9083/stat2/119server/nginx-web

# 컨테이너 리소스 사용량 조회 (전체)
curl -X GET http://localhost:9083/stat3/119server

# SSE 이벤트 스트림 수신
curl -N -H "Accept: text/event-stream" http://localhost:9083/events
```
