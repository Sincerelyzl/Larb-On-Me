package utils

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data",omitempty`
}

func NewErrorResponse(statusCode int, message string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewSuccessResponse(statusCode int, message string, data any) SuccessResponse {
	if data == nil {
		data = struct{}{}
	}
	return SuccessResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
