package awaitws

import "errors"

var (
	// ErrNoID will be thrown if there is no ID in an input message.
	ErrNoID = errors.New("awaitws: no id in input message")

	// ErrorBadID will be thrown if the ID is a bad type (not int).
	ErrorBadID = errors.New("awaitws: id must be an int")
)
