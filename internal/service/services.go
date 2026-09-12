package service

import (
	"database/sql"

	"chat_api/internal/repository"
)

// Services aggregates all application use-case services.
type Services struct {
Users *UsersService
Rooms *RoomsService
Messages *MessagesService
}

// Build wires repositories into services (composition root helper).
func Build(db *sql.DB) *Services {
	return &Services{
Users: newUsersService(repository.NewUsersRepository(db)),
Rooms: newRoomsService(repository.NewRoomsRepository(db)),
Messages: newMessagesService(repository.NewMessagesRepository(db)),
}
}