package apierrors

// ProblemDetails is a minimal RFC 7807-style problem object for JSON responses.
type ProblemDetails struct {
	Type     string         `json:"type"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	Extras   map[string]any `json:"extras,omitempty"`
}

// // AppError is an application error attached with c.Error for the global handler.
// type AppError struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// }

// func (e *AppError) Error() string {
// 	return e.Message
// }
