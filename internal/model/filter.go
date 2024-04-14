package model

import "database/sql"

type Filter struct {
	TagId     sql.NullInt32
	FeatureId sql.NullInt32
}
