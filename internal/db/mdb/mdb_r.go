package mdb

import (
	"context"

	"docker_service/internal/db"
)

func (q *MariaDbHandler) ReadSysdate(ctx context.Context) (string, error) {
	ado := q.GetDB()

	query := `
	select now() as dt from dual
	`

	rows, err := ado.QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	strDateTime := ""
	if rows.Next() {
		if err := rows.Scan(
			&strDateTime,
		); err != nil {
			return "", err
		}
	}
	if err := rows.Close(); err != nil {
		return "", err
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return strDateTime, nil
}

func (q *MariaDbHandler) ReadUser(ctx context.Context, username string) (db.User, error) {
	ado := q.GetDB()
	var u db.User

	query := `
	SELECT username
		 , hashed_password
		 , full_name
		 , email
		 , password_changed_at
		 , created_at
		 FROM users
	WHERE username = ?
	`

	rows, err := ado.QueryContext(ctx, query, username)
	if err != nil {
		return db.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&u.Username,
			&u.HashedPassword,
			&u.FullName,
			&u.Email,
			&u.PasswordChangedAt,
			&u.CreatedAt,
		); err != nil {
			return db.User{}, err
		}
	}
	if err := rows.Close(); err != nil {
		return db.User{}, err
	}
	if err := rows.Err(); err != nil {
		return db.User{}, err
	}

	return u, nil
}
