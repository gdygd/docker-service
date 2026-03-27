package mdb

import (
	"context"

	"docker_service/internal/db"
	"docker_service/internal/logger"
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

func (q *MariaDbHandler) ReadUserSession(ctx context.Context, id string) (db.Session, error) {
	ado := q.GetDB()
	var se db.Session

	query := `
	SELECT ID          
		 , username    
		 , refresh_token
		 , user_agent   
		 , client_ip    
		 , block_yn   
		 , expires_at   
		 , created_at   
	FROM sessions
	WHERE ID = ?
	`

	rows, err := ado.QueryContext(ctx, query, id)
	if err != nil {
		return db.Session{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&se.ID,
			&se.Username,
			&se.RefreshToken,
			&se.UserAgent,
			&se.ClientIp,
			&se.IsBlocked,
			&se.ExpiresAt,
			&se.CreatedAt,
		); err != nil {
			return db.Session{}, err
		}
	}
	if err := rows.Close(); err != nil {
		return db.Session{}, err
	}
	if err := rows.Err(); err != nil {
		return db.Session{}, err
	}

	return se, nil
}

func (q *MariaDbHandler) ReadHost(ctx context.Context) ([]db.Host, error) {
	ado := q.GetDB()

	query := `
	select a.host_id, a.hostname, a.host_address, ifnull(a.mode, 1) mode from host_info a
	`

	rows, err := ado.QueryContext(ctx, query)
	if err != nil {
		logger.Log.Error("ReadHost#1 error %v", err)
		return nil, err
	}
	defer rows.Close()

	var rst []db.Host = []db.Host{}

	for rows.Next() {
		row := db.Host{}
		if err := rows.Scan(
			&row.HostId,
			&row.HostName,
			&row.HostAddress,
			&row.Mode,
		); err != nil {
			logger.Log.Error("ReadHost#2 error %v", err)
			return nil, err
		}
		rst = append(rst, row)
	}
	if err := rows.Close(); err != nil {
		logger.Log.Error("ReadHost#3 error %v", err)
		return nil, err
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("ReadHost#4 error %v", err)
		return nil, err
	}
	return rst, nil
}

func (q *MariaDbHandler) ReadHostInfo(ctx context.Context, hostid int) (db.Host, error) {
	ado := q.GetDB()

	query := `
	select a.host_id, a.hostname, a.host_address 
	from host_info a
	where a.host_id = ?
	`

	var rst db.Host
	rows, err := ado.QueryContext(ctx, query, hostid)
	if err != nil {
		logger.Log.Error("ReadHost#1 error %v", err)
		return rst, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&rst.HostId,
			&rst.HostName,
			&rst.HostAddress,
		); err != nil {
			logger.Log.Error("ReadHost#2 error %v", err)
			return rst, err
		}
	}
	if err := rows.Close(); err != nil {
		logger.Log.Error("ReadHost#3 error %v", err)
		return rst, err
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("ReadHost#4 error %v", err)
		return rst, err
	}
	return rst, nil
}

// func (q *MariaDbHandler) DeleteUserSession(ctx context.Context, id string) error {
// 	ado := q.GetDB()

// 	query := `DELETE FROM sessions WHERE ID = ?`

// 	_, err := ado.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
