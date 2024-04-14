package model

import "database/sql"

type Page struct {
	Limit  sql.NullInt32
	Offset sql.NullInt32
}
