package model

type ReturnData[T any] struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *T
}

func NewReturnData[T any](code int, success bool, message string, data *T) *ReturnData[T] {
	return &ReturnData[T]{
		Code:    code,
		Success: success,
		Message: message,
		Data:    data,
	}
}
