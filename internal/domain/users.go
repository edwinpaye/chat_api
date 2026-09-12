package domain

import "time"

// Users — application users
type Users struct {
Id string `db:"id" json:"id"`
Username string `db:"username" json:"username"`
Email string `db:"email" json:"email"`
Status map[string]any `db:"status" json:"status,omitempty"`
CreatedAt time.Time `db:"created_at" json:"createdAt"`
UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

