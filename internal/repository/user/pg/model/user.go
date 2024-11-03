package model

import (
	"database/sql"
	"time"
)

// User модель для работы в репо слое
type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	UserRole  int32        `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
