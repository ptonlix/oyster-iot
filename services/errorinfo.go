package services

import "errors"

var (
	ErrMQTTToken = errors.New("Payload Token missing")
	ErrMQTTMsg   = errors.New("Payload Msg missing")
	ErrDevType   = errors.New("Device Type Error")
)
