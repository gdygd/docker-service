-- Database 생성
CREATE DATABASE IF NOT EXISTS `docker` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `docker`;

-- 사용자
CREATE TABLE `USERS` (
	`USER_NM`   VARCHAR(50)  NOT NULL, -- 사용자 이름
	`PASSWD`    VARCHAR(256) NULL,     -- 비밀번호
	`EMAIL`     VARCHAR(40)  NULL,     -- 이메일
	`CHG_DT`    DATE         NULL,     -- 변경 일자
	`CREATE_DT` DATE         NULL      -- 생성 일자
);

-- 사용자 기본키
ALTER TABLE `USERS`
	ADD CONSTRAINT `PK_USERS`
	PRIMARY KEY (
	    `USER_NM`
	);

-- 세션
CREATE TABLE `SESSIONS` (
	`ID`         VARCHAR(255) NOT NULL, -- ID
	`USER_NM`    VARCHAR(50)  NULL,     -- 사용자 이름
	`REF_TOKEN`  VARCHAR(512) NULL,     -- 리프레시 토큰
	`USER_AGENT` VARCHAR(255) NULL,     -- 사용자 에이전트
	`CLIENT_IP`  VARCHAR(20)  NULL,     -- 접속 IP
	`BLOCK_YN`   DECIMAL(1)   NULL,     -- 블락 여부
	`EXP_DT`     DATETIME     NULL,     -- 만료 시간
	`CREATE_DT`  DATETIME     NULL      -- 생성 시간
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
	    `USER_NM`
	)
	REFERENCES `USERS` (
	    `USER_NM`
	);
