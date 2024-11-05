package pkg

// Validator интерфейс для валидации запросов
type Validator interface {
	ValidateRequest() error
}
