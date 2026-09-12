package domain

import "time"

type Rooms struct {
Id int64 `db:"id" json:"id"`
Name string `db:"name" json:"name"`
}

var _ = time.Time{}
