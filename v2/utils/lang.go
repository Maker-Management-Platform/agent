package utils

func Ptr[T any](v T) *T {
	return &v
}
func VoZ[T any](v *T) T {
	if v != nil {
		return *v
	}
	return *new(T)
}
