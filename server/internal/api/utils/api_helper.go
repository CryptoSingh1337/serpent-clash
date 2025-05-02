package utils

type Response[T any] struct {
	Data  T      `json:"data"`
	Error *Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

func CreateResponse[T any](data T, err *Error) Response[T] {
	if err != nil {
		return Response[T]{
			Error: err,
		}
	}
	return Response[T]{
		Data: data,
	}
}
