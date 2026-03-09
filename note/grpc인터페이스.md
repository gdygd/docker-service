# gRPC 인터페이스 설계

## 1. 서비스 정의

```protobuf
syntax = "proto3";
package containeragent;
option go_package = "proto/agentpb";

service AgentService {
    // Agent -> SaaS: 양방향 스트리밍 (데이터 전송 + 원격 제어 수신)
    rpc DataStream(stream AgentMessage) returns (stream ServerMessage);

    // Agent -> SaaS: 연결 시 핸드셰이크
    rpc Register(RegisterRequest) returns (RegisterResponse);

    // Agent -> SaaS: Heartbeat (연결 상태 확인)
    rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
}
```

---

## 2. 데이터 타입

```protobuf
enum DataType {
    CONTAINER_LIST = 0;
    CONTAINER_INSPECT = 1;
    CONTAINER_STATS = 2;
    CONTAINER_EVENT = 3;
}

enum CommandType {
    ACK = 0;
    START_CONTAINER = 1;
    STOP_CONTAINER = 2;
    RESTART_CONTAINER = 3;
}
```

---

## 3. Agent -> SaaS 메시지

### 3-1. AgentMessage (스트리밍 데이터)

```protobuf
message AgentMessage {
    string agent_key = 1;
    DataType type = 2;
    string host = 3;
    int64 timestamp = 4;
    oneof data {
        ContainerListData list_data = 10;
        ContainerInspectData inspect_data = 11;
        ContainerStatsData stats_data = 12;
        ContainerEventData event_data = 13;
    }
}
```

### 3-2. Register (핸드셰이크)

```protobuf
message RegisterRequest {
    string agent_key = 1;
    string agent_name = 2;
    repeated HostInfo hosts = 3;
}

message HostInfo {
    string name = 1;
    string addr = 2;
}

message RegisterResponse {
    bool success = 1;
    string message = 2;
    string agent_id = 3;
}
```

### 3-3. Heartbeat

```protobuf
message HeartbeatRequest {
    string agent_key = 1;
    int64 timestamp = 2;
}

message HeartbeatResponse {
    bool success = 1;
    int64 server_time = 2;
}
```

---

## 4. SaaS -> Agent 메시지

```protobuf
message ServerMessage {
    CommandType command = 1;
    string target_container = 2;
    string host = 3;
}
```

---

## 5. 데이터 메시지 상세

### 5-1. ContainerListData

```protobuf
message ContainerListData {
    repeated ContainerInfo containers = 1;
}

message ContainerInfo {
    string id = 1;
    string name = 2;
    string image = 3;
    string state = 4;
    string status = 5;
}
```

### 5-2. ContainerStatsData

```protobuf
message ContainerStatsData {
    repeated ContainerStats stats = 1;
}

message ContainerStats {
    string id = 1;
    string name = 2;
    double cpu_percent = 3;
    uint64 memory_usage = 4;
    uint64 memory_limit = 5;
    double memory_percent = 6;
    uint64 network_rx = 7;
    uint64 network_tx = 8;
}
```

### 5-3. ContainerInspectData

```protobuf
message ContainerInspectData {
    repeated ContainerInspect inspects = 1;
}

message ContainerInspect {
    string id = 1;
    string name = 2;
    string image = 3;
    string created = 4;
    string platform = 5;
    ContainerState state = 10;
    ContainerConfig config = 11;
    ContainerNetwork network = 12;
    repeated MountPoint mounts = 13;
}

message ContainerState {
    string status = 1;
    bool running = 2;
    bool paused = 3;
    bool restarting = 4;
    int32 exit_code = 5;
    string started_at = 6;
    string finished_at = 7;
}

message ContainerConfig {
    string hostname = 1;
    string user = 2;
    repeated string env = 3;
    repeated string cmd = 4;
    repeated string entrypoint = 5;
    string working_dir = 6;
    map<string, string> labels = 7;
}

message ContainerNetwork {
    string ip_address = 1;
    string gateway = 2;
    string mac_address = 3;
    map<string, PortBindings> ports = 4;
    map<string, NetworkEndpoint> networks = 5;
}

message PortBindings {
    repeated PortBinding bindings = 1;
}

message PortBinding {
    string host_ip = 1;
    string host_port = 2;
}

message NetworkEndpoint {
    string network_id = 1;
    string ip_address = 2;
    string gateway = 3;
    string mac_address = 4;
}

message MountPoint {
    string type = 1;
    string name = 2;
    string source = 3;
    string destination = 4;
    string mode = 5;
    bool rw = 6;
}
```

### 5-4. ContainerEventData

```protobuf
message ContainerEventData {
    string type = 1;       // container, network, image, volume
    string action = 2;     // start, stop, die, create, destroy
    string actor_id = 3;
    string actor_name = 4;
    int64 timestamp = 5;
    map<string, string> attrs = 6;
}
```

---

## 6. 통신 흐름

```
Agent                                    SaaS Server
  │                                          │
  │──── Register(agent_key, hosts) ─────────>│  핸드셰이크
  │<─── RegisterResponse(agent_id) ──────────│
  │                                          │
  │════ DataStream (양방향) ════════════════>│  스트리밍 시작
  │                                          │
  │──── AgentMessage(list_data) ────────────>│  컨테이너 목록
  │──── AgentMessage(inspect_data) ─────────>│  Inspect 정보
  │──── AgentMessage(stats_data) ───────────>│  Stats 정보
  │──── AgentMessage(event_data) ───────────>│  이벤트
  │                                          │
  │<─── ServerMessage(START_CONTAINER) ──────│  원격 제어 (향후)
  │                                          │
  │──── Heartbeat ──────────────────────────>│  30초 주기
  │<─── HeartbeatResponse ──────────────────│
  │                                          │
```

---

## 7. 인증 방식

gRPC Metadata를 통한 agent_key 인증:

```
Metadata:
  authorization: Bearer ak_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

- Register, DataStream, Heartbeat 모든 RPC에 agent_key 필수
- Server Interceptor에서 agent_key 검증
- 유효하지 않으면 `codes.Unauthenticated` 반환
