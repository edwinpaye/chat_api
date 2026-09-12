package domain

import "time"

type Messages struct {
Id int64 `db:"id" json:"id"`
RoomId int64 `db:"room_id" json:"roomId"`
AuthorId string `db:"author_id" json:"authorId"`
Body string `db:"body" json:"body"`
Status string `db:"status" json:"status,omitempty"`
CreatedAt time.Time `db:"created_at" json:"createdAt"`
UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

