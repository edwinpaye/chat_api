package domain

import "time"

// Users — application users
type Users struct {
Id string `db:"id" json:"id"`
Username string `db:"username" json:"username"`
Email string `db:"email" json:"email"`
Status *bool `db:"status" json:"status,omitempty"`
Online *bool `db:"online" json:"online,omitempty"`
LastSeen *time.Time `db:"last_seen" json:"lastSeen,omitempty"`
CreatedAt time.Time `db:"created_at" json:"createdAt"`
UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

