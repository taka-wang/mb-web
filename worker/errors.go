package worker

import "errors"

var (
	// ErrMarshal is the error when marshalling to JSON string failed.
	ErrMarshal = errors.New("Fail to marshal!")

	// ErrUnmarshal is the error when unmarshalling JSON string to structure failed.
	ErrUnmarshal = errors.New("Fail to unmarshal!")

	// ErrRequestNotSupport is the error when the request is not supported.
	ErrRequestNotSupport = errors.New("Request not support!")

	// ErrResponseNotSupport is the error when the response is not supported.
	ErrResponseNotSupport = errors.New("Response not support!")

	// ErrInvalidMessageLength is the error when the length of message is invalid.
	ErrInvalidMessageLength = errors.New("Invalid message length!")

	// ErrFailToGetFromMap is the error when retrieving item failed.
	ErrFailToGetFromMap = errors.New("Fail to retrieve item from map.")
)
