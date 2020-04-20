package internal

type ClientError struct {
	message string
}

func NewClientError(message string) ClientError {
	return ClientError{
		message: message,
	}
}

func (err ClientError) Error() string {
	return err.message
}

func IsClientError(err error) bool {
	_, ok := err.(ClientError)
	return ok
}