package repository

type Response struct {
	Status  string
	Message string
	Data    interface{}
}

func Err(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}

func Success(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}
