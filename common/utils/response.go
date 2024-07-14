package utils

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func NewErrorResponse(statusCode int, message string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewSuccessResponse(statusCode int, message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
