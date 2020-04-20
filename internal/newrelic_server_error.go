package internal

type newRelicServerError struct {}

func (e newRelicServerError) Error () string {
	return "A New Relic server error occurred"
}

var NewRelicServerError = newRelicServerError{}
