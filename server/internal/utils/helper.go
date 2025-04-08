package utils

import (
	"encoding/json"
	"math"
)

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

func LerpAngle(a, b, t float64) float64 {
	diff := b - a
	// Handle wrapping from -π to π
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	return a + diff*math.Min(t, 1.0)
}

func EuclideanDistance(a, b, x, y float64) float64 {
	dy := y - b
	dx := x - a
	return math.Sqrt(dy*dy + dx*dx)
}

func RemoveFromSlice[T any](slice []T, index int) []T {
	slice[index] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice
}
