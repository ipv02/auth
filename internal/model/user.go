package model

import (
	"database/sql"
	"time"
)

// UserCreate модель для конвертации из протомодели в модель бизнес-логики
type UserCreate struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Role            int32  `json:"role"`
}

// UserGet модель для конвертации из протомодели в модель бизнес-логики
type UserGet struct {
	ID        int64
	Name      string
	Email     string
	UserRole  int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserUpdate модель для конвертации из протомодели в модель бизнес-логики
type UserUpdate struct {
	ID    int64
	Name  *string
	Email *string
	Role  int32
}
