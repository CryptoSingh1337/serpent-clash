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

func ToJsonS(data any) (string, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func ToJsonB(data any) ([]byte, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	return val, nil
}

func FromJsonS[T any](data string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func FromJsonB[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
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
