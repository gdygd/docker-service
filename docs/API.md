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

## 1. GET /hosts

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

## 2. GET /ps2/:host

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

## 3. GET /inspect2/:host/:id

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

## 4. POST /start2

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

## 5. POST /stop2

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

## 6. GET /stat2/:host/:id

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

## 7. GET /events (SSE)

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

# 컨테이너 리소스 사용량 조회
curl -X GET http://localhost:9083/stat2/119server/nginx-web

# SSE 이벤트 스트림 수신
curl -N -H "Accept: text/event-stream" http://localhost:9083/events
```
