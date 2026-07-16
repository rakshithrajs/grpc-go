package models

import "database/sql"

type Postgres struct {
	Db *sql.DB
}
