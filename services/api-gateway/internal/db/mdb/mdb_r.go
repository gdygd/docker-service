package mdb

import (
	"context"
)

func (q *MariaDbHandler) ReadSysdate(ctx context.Context) (string, error) {
	db := q.GetDB()

	query := `
	select now() as dt from dual
	`

	rows, err := db.QueryContext(ctx, query)
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
