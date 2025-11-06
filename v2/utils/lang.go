package utils

import "encoding/json"

func Ptr[T any](v T) *T {
	return &v
}
func VoZ[T any](v *T) T {
	if v != nil {
		return *v
	}
	return *new(T)
}

type Option[T any] struct {
	Value T
	Valid bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{Value: v, Valid: true}
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.Valid {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Valid = false
		return nil
	}
	o.Valid = true
	return json.Unmarshal(data, &o.Value)
}
