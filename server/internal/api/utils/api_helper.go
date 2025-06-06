package utils

import "encoding/json"

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

func ToJsonB(data any) ([]byte, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	return val, nil
}
