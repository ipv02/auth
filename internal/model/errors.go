package model

import "github.com/pkg/errors"

// ErrorUserNotFound глобальная переменная хранящая ошибку с сообщением
var ErrorUserNotFound = errors.New("user not found")
