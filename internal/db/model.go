package db

import (
	"time"
)

type USER struct {
	USER_NM   string
	PASSWD    string
	EMAIL     string
	CHG_DT    time.Time
	CREATE_DT time.Time
}
