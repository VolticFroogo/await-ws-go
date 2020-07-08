package awaitws

import "errors"

var (
	ErrNoID    = errors.New("awaitws: no id in input message")
	ErrorBadID = errors.New("awaitws: id must be an int")
)
