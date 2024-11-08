package model

// User модель для работы c redis
type User struct {
	ID          int64  `redis:"id"`
	Name        string `redis:"name"`
	Email       string `redis:"email"`
	UserRole    int32  `redis:"role"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at"`
}

// UserCreate модель для работы c redis
type UserCreate struct {
	Name            string `redis:"name"`
	Email           string `redis:"email"`
	Password        string `redis:"password"`
	PasswordConfirm string `redis:"password_confirm"`
	Role            int32  `redis:"role"`
}

// UserUpdate модель для для работы c redis
type UserUpdate struct {
	ID    int64   `redis:"id"`
	Name  *string `redis:"name"`
	Email *string `redis:"email"`
	Role  int32   `redis:"role"`
}
