package errors

import (
	"errors"
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrConfig represents a configuration error
	ErrConfig ErrorType = "config"
	// ErrAPI represents an API error
	ErrAPI ErrorType = "api"
	// ErrValidation represents a validation error
	ErrValidation ErrorType = "validation"
	// ErrNetwork represents a network error
	ErrNetwork ErrorType = "network"
)

// Error represents a custom error with context
type Error struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error
func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewError creates a new Error
func NewError(errType ErrorType, message string, cause error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Cause:   err,
	}
}

// IsType checks if the error is of a specific type
func IsType(err error, errType ErrorType) bool {
	var e *Error
	if errors.As(err, &e) { // 使用 errors.As
		return e.Type == errType
	}
	return false
}

// GetContext returns the context value for a key
func GetContext(err error, key string) (interface{}, bool) {
	var e *Error
	if errors.As(err, &e) { // 使用 errors.As
		if e.Context != nil {
			if val, ok := e.Context[key]; ok {
				return val, true
			}
		}
	}
	return nil, false
}

func main() {
	// 示例用法
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, ErrAPI, "API request failed").WithContext("request_id", "12345")

	fmt.Println(wrappedErr) // 输出：api: API request failed: original error

	// 使用 IsType 检查错误类型
	if IsType(wrappedErr, ErrAPI) {
		fmt.Println("Error is of type API")
	}

	if !IsType(wrappedErr, ErrConfig) {
		fmt.Println("Error is not of type Config")
	}

	// 获取上下文信息
	if reqID, ok := GetContext(wrappedErr, "request_id"); ok {
		fmt.Println("Request ID:", reqID) // 输出：Request ID: 12345
	}

	if errors.Is(wrappedErr, originalErr) {
		fmt.Println("wrappedErr contains originalErr")
	}
}
