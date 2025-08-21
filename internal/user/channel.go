package user

import (
	"github.com/google/uuid"
)

type Channel struct {
	ID uuid.UUID
	Name string
	Owner *User
	Users []*User
}