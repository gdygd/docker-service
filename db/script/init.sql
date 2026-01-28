-- Database 생성
CREATE DATABASE IF NOT EXISTS `docker` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `docker`;

-- 사용자
CREATE TABLE `USERS` (
	`username`            VARCHAR(50)  NOT NULL, -- 사용자ID
	`hashed_password`     VARCHAR(256) NULL,     -- 비밀번호
	`full_name`           VARCHAR(50)  NULL,     -- 사용자이름
	`email`               VARCHAR(40)  NULL,     -- 이메일
	`password_changed_at` DATETIME     NULL,     -- 비밀번호 변경일시
	`created_at`          DATETIME     NULL      -- 생성일시
);

-- 사용자 기본키
ALTER TABLE `USERS`
	ADD CONSTRAINT `PK_USERS`
	PRIMARY KEY (
	    `username`
	);

-- 세션
CREATE TABLE `SESSIONS` (
	`ID`            VARCHAR(255) NOT NULL, -- ID
	`username`      VARCHAR(50)  NULL,     -- 사용자ID
	`refresh_token` VARCHAR(512) NULL,     -- 리프레시 토큰
	`user_agent`    VARCHAR(255) NULL,     -- 사용자 에이전트
	`client_ip`     VARCHAR(20)  NULL,     -- 접속 IP
	`is_blocked`    DECIMAL(1)   NULL,     -- 블락 여부
	`expires_at`    DATETIME     NULL,     -- 만료 시간
	`created_at`    DATETIME     NULL      -- 생성 시간
);

-- 세션 기본키
ALTER TABLE `SESSIONS`
	ADD CONSTRAINT `PK_SESSIONS`
	PRIMARY KEY (
	    `ID`
	);

-- 세션 -> 사용자 FK
ALTER TABLE `SESSIONS`
	ADD CONSTRAINT `FK_USERS_TO_SESSIONS`
	FOREIGN KEY (
	    `username`
	)
	REFERENCES `USERS` (
	    `username`
	);
