package error

type MinioError struct {
	message string
}

func (e *MinioError) Error() string {
	return e.message
}
