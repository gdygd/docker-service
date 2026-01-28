package db

import (
	"database/sql"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

type User struct {
	Username          string       `json:"username"`
	HashedPassword    string       `json:"hashed_password"`
	FullName          string       `json:"full_name"`
	Email             string       `json:"email"`
	PasswordChangedAt sql.NullTime `json:"password_changed_at"`
	CreatedAt         sql.NullTime `json:"created_at"`
}

type CreateSessionParams struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	RefreshToken string       `json:"refresh_token"`
	UserAgent    string       `json:"user_agent"`
	ClientIp     string       `json:"client_ip"`
	IsBlocked    bool         `json:"is_blocked"`
	ExpiresAt    sql.NullTime `json:"expires_at"`
}

type Session struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	RefreshToken string       `json:"refresh_token"`
	UserAgent    string       `json:"user_agent"`
	ClientIp     string       `json:"client_ip"`
	IsBlocked    bool         `json:"is_blocked"`
	ExpiresAt    sql.NullTime `json:"expires_at"`
	CreatedAt    sql.NullTime `json:"created_at"`
}
