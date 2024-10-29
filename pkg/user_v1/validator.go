package user_v1

import (
	"fmt"
	"github.com/pkg/errors"
)

// Validate валидация CreateUserRequest
func (req *CreateUserRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password != req.PasswordConfirm {
		return errors.New("passwords do not match")
	}
	return nil
}

// Validate валидация UpdateUserRequest
func (req *UpdateUserRequest) Validate() error {
	fmt.Println("--- ROLE", req.Role)
	if req.Role == UserRole_UNKNOWN {
		return errors.New("role is unknown")
	}

	return nil
}

// Validate валидация DeleteUserRequest
func (req *DeleteUserRequest) Validate() error {
	if req.Id == 0 {
		return errors.New("id is required")
	}
	return nil
}

// Validate валидация GetUserRequest
func (req *GetUserRequest) Validate() error {
	if req.Id == 0 {
		return errors.New("id is required")
	}
	return nil
}
