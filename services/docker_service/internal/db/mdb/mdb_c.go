package mdb

import (
	"context"

	"docker_service/internal/db"
)

func (q *MariaDbHandler) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	ado := q.GetDB()

	query := `
	INSERT INTO users (
		username,
		hashed_password,
		full_name,
		email,
		password_changed_at,
		created_at

	) VALUES (
		?, ?, ?, ?, now(), now()
	)
	RETURNING username, hashed_password, full_name, email, password_changed_at, created_at
	`

	row := ado.QueryRow(query,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var u db.User
	err := row.Scan(
		&u.Username,
		&u.HashedPassword,
		&u.FullName,
		&u.Email,
		&u.PasswordChangedAt,
		&u.CreatedAt,
	)
	if err != nil {
		return db.User{}, err
	}
	return u, err
}

func (q *MariaDbHandler) CreateUserSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	ado := q.GetDB()

	query := `
	INSERT INTO sessions (ID, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at)
	VALUES(?, ?, ?, ?, ?, ?, ?, now()) RETURNING ID, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
	`

	row := ado.QueryRow(query,
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var se db.Session
	err := row.Scan(
		&se.ID,
		&se.Username,
		&se.RefreshToken,
		&se.UserAgent,
		&se.ClientIp,
		&se.IsBlocked,
		&se.ExpiresAt,
		&se.CreatedAt,
	)
	if err != nil {
		return db.Session{}, err
	}
	return se, err
}
