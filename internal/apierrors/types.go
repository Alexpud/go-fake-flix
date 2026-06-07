package apierrors

// ProblemDetails is a minimal RFC 7807-style problem object for JSON responses.
type ProblemDetails struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail,omitempty"`
	Code     string       `json:"code,omitempty"`
	Instance string       `json:"instance,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CreateProblemDetails(code string, message string, statusCode int) *ProblemDetails {
	return &ProblemDetails{
		Title:    "An error happened while processing the request",
		Status:   statusCode,
		Detail:   message,
		Type:     "ProblemDetails",
		Instance: "",
		Code:     code,
	}
}
