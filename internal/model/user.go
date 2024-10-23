package model

import (
	"database/sql"
	"time"
)

type UserCreate struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int32
}

type UserGet struct {
	ID        int64
	Name      string
	Email     string
	UserRole  int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserUpdate struct {
	ID    int64
	Name  string
	Email string
	Role  int32
}
