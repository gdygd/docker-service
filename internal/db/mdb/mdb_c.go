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
		email
	) VALUES (
		?, ?, ?, ?
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
