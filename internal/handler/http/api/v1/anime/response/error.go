package response

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
	}
}
