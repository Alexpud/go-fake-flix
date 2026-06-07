package common

type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}
