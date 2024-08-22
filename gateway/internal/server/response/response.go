package response

// Response struct is a JSON encoded response from the API.
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewErrorResponse constructs the Response that represents a request
// failure from the API.
func NewErrorResponse(msg string) Response[*struct{}] {
	return Response[*struct{}]{
		Error:   msg,
		Success: false,
	}
}

// NewSuccessResponse constructs the Response that represents a request
// success from the API.
func NewSuccessResponse[T any](msg string, data T) *Response[T] {
	return &Response[T]{
		Message: msg,
		Success: true,
		Data:    data,
	}
}
