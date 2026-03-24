package mdb

import "context"

func (q *MariaDbHandler) DeleteUserSession(ctx context.Context, id string) error {
	ado := q.GetDB()

	query := `DELETE FROM sessions WHERE ID = ?`

	_, err := ado.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
