package user_v1

import (
	"github.com/pkg/errors"
)

// ValidateRequest валидация CreateUserRequest
func (req *CreateUserRequest) ValidateRequest() error {
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

// ValidateRequest валидация UpdateUserRequest
func (req *UpdateUserRequest) ValidateRequest() error {
	if req.Role == UserRole_UNKNOWN {
		return errors.New("role is unknown")
	}

	return nil
}

// ValidateRequest валидация DeleteUserRequest
func (req *DeleteUserRequest) ValidateRequest() error {
	if req.Id == 0 {
		return errors.New("id is required")
	}
	return nil
}

// ValidateRequest валидация GetUserRequest
func (req *GetUserRequest) ValidateRequest() error {
	if req.Id == 0 {
		return errors.New("id is required")
	}
	return nil
}
