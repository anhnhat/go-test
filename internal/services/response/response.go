package response

type Response struct {
	Data       any    `json:"data,omitempty"`
	ErrMessage string `json:"err_message,omitempty"`
}

func NewResponse(data any) Response {
	return Response{
		Data: data,
	}
}

func NewErrorResponse(errorMessage string) Response {
	return Response{
		Data:       nil,
		ErrMessage: errorMessage,
	}
}

func NewInvalidRequestResponse() Response {
	return Response{
		Data:       nil,
		ErrMessage: "invalid request",
	}
}
