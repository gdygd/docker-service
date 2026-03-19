-- Database 생성
CREATE DATABASE IF NOT EXISTS `docker` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `docker`;

-- 사용자
CREATE TABLE `USERS` (
	`username`            VARCHAR(100) NOT NULL, -- 사용자ID
	`hashed_password`     VARCHAR(256) NULL,     -- 해시 비밀번호
	`full_name`           VARCHAR(100) NULL,     -- 사용자이름
	`email`               VARCHAR(40)  NULL,     -- 이메일
	`password_changed_at` DATETIME     NULL,     -- 비밀번호 변경일시
	`created_at`          DATETIME     NULL      -- 생성일시
);

-- 사용자
ALTER TABLE `USERS`
	ADD CONSTRAINT `PK_USERS` -- 사용자 기본키
	PRIMARY KEY (
	    `username` -- 사용자ID
	);

-- 세션
CREATE TABLE `SESSIONS` (
	`id`            VARCHAR(255) NOT NULL, -- ID
	`username`      VARCHAR(100) NULL,     -- 사용자ID
	`refresh_token` VARCHAR(512) NULL,     -- 리프레시 토큰
	`user_agent`    VARCHAR(255) NULL,     -- 사용자 에이전트
	`client_ip`     VARCHAR(20)  NULL,     -- 접속 IP
	`block_yn`      DECIMAL(1)   NULL,     -- 블락 여부
	`expires_at`    DATETIME     NULL,     -- 만료 시간
	`created_at`    DATETIME     NULL      -- 생성 시간
);

-- 세션
ALTER TABLE `SESSIONS`
	ADD CONSTRAINT `PK_SESSIONS` -- 세션 기본키
	PRIMARY KEY (
	    `id` -- ID
	);

-- 컨테이너 정보
CREATE TABLE `container_info` (
	`id`             INT          NOT NULL, -- ID
	`container_id`   VARCHAR(32)  NOT NULL, -- 컨테이너 ID
	`container_name` VARCHAR(100) NULL,     -- 컨테이너 명
	`image`          VARCHAR(100) NULL,     -- 이미지
	`state`          VARCHAR(64)  NULL,     -- 상태
	`status`         VARCHAR(128) NULL,     -- 상태 메시지
	`changed_at`     DATETIME     NULL      -- 변경 일시
);

-- 컨테이너 정보
ALTER TABLE `container_info`
	ADD CONSTRAINT `PK_container_info` -- 컨테이너 정보 기본키
	PRIMARY KEY (
	    `id`,           -- ID
	    `container_id`  -- 컨테이너 ID
	);

-- 도커 호스트
CREATE TABLE `docker_host` (
	`id`            INT          NOT NULL, -- ID
	`agent_key`     VARCHAR(256) NOT NULL, -- 에이전트 키
	`hostname`      VARCHAR(256) NULL,     -- 호스트명
	`agent_address` VARCHAR(255) NULL,     -- 에이전트 주소
	`updated_at`    DATETIME     NULL      -- 갱신 일시
);

-- 도커 호스트
ALTER TABLE `docker_host`
	ADD CONSTRAINT `PK_docker_host` -- 도커 호스트 기본키
	PRIMARY KEY (
	    `id` -- ID
	);

-- 컨테이너 인스펙트
CREATE TABLE `container_inspect` (
	`id`             INT          NOT NULL, -- ID
	`container_id`   VARCHAR(32)  NOT NULL, -- 컨테이너 ID
	`container_name` VARCHAR(100) NULL,     -- 컨테이너 명
	`image`          VARCHAR(100) NULL,     -- 이미지
	`platform`       VARCHAR(64)  NULL,     -- 플랫폼
	`restart_count`  INT          NULL,     -- 재시작 횟수
	`state_info`     JSON         NULL,     -- 상태 정보
	`config_info`    JSON         NULL,     -- 설정 정보
	`network_info`   JSON         NULL,     -- 네트워크 정보
	`mount_info`     JSON         NULL,     -- 마운트 정보
	`changed_at`     DATETIME     NULL      -- 변경 일시
);

-- 컨테이너 인스펙트
ALTER TABLE `container_inspect`
	ADD CONSTRAINT `PK_container_inspect` -- 컨테이너 인스펙트 기본키
	PRIMARY KEY (
	    `id`,           -- ID
	    `container_id`  -- 컨테이너 ID
	);

-- 컨테이너 리소스 이력
CREATE TABLE `container_stats` (
	`id`             INT          NOT NULL, -- ID
	`container_id`   VARCHAR(32)  NOT NULL, -- 컨테이너 ID
	`collected_at`   DATETIME     NOT NULL, -- 수집 일시
	`container_name` VARCHAR(100) NULL,     -- 컨테이너 명
	`cpu_percent`    DECIMAL(5,2) NULL,     -- cpu 사용률
	`memory_usage`   INT          NULL,     -- 메모리 사용량
	`memory_limit`   INT          NULL,     -- 메모리 제한
	`memory_percent` DECIMAL(5,2) NULL,     -- 메모리 사용률
	`network_rx`     INT          NULL,     -- 네트워크 수신량
	`COL`            INT          NULL      -- 네트워크 송신량
);

-- 컨테이너 리소스 이력
ALTER TABLE `container_stats`
	ADD CONSTRAINT `PK_container_stats` -- 컨테이너 리소스 이력 기본키
	PRIMARY KEY (
	    `id`,           -- ID
	    `container_id`, -- 컨테이너 ID
	    `collected_at`  -- 수집 일시
	);

-- 컨테이너 이벤트 이력
CREATE TABLE `container_event` (
	`id`              INT          NOT NULL, -- ID
	`container_id`    VARCHAR(32)  NOT NULL, -- 컨테이너 ID
	`received_at`     DATETIME     NOT NULL, -- 수신 일시
	`hostname`        VARCHAR(100) NULL,     -- 호스트명
	`type`            DECIMAL(2)   NULL,     -- 타입
	`action`          VARCHAR(64)  NULL,     -- 액션
	`actor_id`        VARCHAR(128) NULL,     -- 대상 ID
	`actor_name`      VARCHAR(256) NULL,     -- 대상 명
	`event_timestamp` BIGINT       NULL,     -- 이벤트 발생시각
	`attrs`           JSON         NULL      -- 부가 속성
);

-- 컨테이너 이벤트 이력
ALTER TABLE `container_event`
	ADD CONSTRAINT `PK_container_event` -- 컨테이너 이벤트 이력 기본키
	PRIMARY KEY (
	    `id`,           -- ID
	    `container_id`, -- 컨테이너 ID
	    `received_at`   -- 수신 일시
	);

-- 세션
ALTER TABLE `SESSIONS`
	ADD CONSTRAINT `FK_USERS_TO_SESSIONS` -- 사용자 -> 세션
	FOREIGN KEY (
	    `username` -- 사용자ID
	)
	REFERENCES `USERS` ( -- 사용자
	    `username` -- 사용자ID
	);

-- 컨테이너 정보
ALTER TABLE `container_info`
	ADD CONSTRAINT `FK_docker_host_TO_container_info` -- 도커 호스트 -> 컨테이너 정보
	FOREIGN KEY (
	    `id` -- ID
	)
	REFERENCES `docker_host` ( -- 도커 호스트
	    `id` -- ID
	);

-- 컨테이너 인스펙트
ALTER TABLE `container_inspect`
	ADD CONSTRAINT `FK_docker_host_TO_container_inspect` -- 도커 호스트 -> 컨테이너 인스펙트
	FOREIGN KEY (
	    `id` -- ID
	)
	REFERENCES `docker_host` ( -- 도커 호스트
	    `id` -- ID
	);

-- 컨테이너 리소스 이력
ALTER TABLE `container_stats`
	ADD CONSTRAINT `FK_docker_host_TO_container_stats` -- 도커 호스트 -> 컨테이너 리소스 이력
	FOREIGN KEY (
	    `id` -- ID
	)
	REFERENCES `docker_host` ( -- 도커 호스트
	    `id` -- ID
	);

-- 컨테이너 이벤트 이력
ALTER TABLE `container_event`
	ADD CONSTRAINT `FK_docker_host_TO_container_event` -- 도커 호스트 -> 컨테이너 이벤트 이력
	FOREIGN KEY (
	    `id` -- ID
	)
	REFERENCES `docker_host` ( -- 도커 호스트
	    `id` -- ID
	);