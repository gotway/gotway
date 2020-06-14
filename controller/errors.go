package controller

import "errors"

// ErrAlreadyRegistered error for service already registered
var ErrAlreadyRegistered = errors.New("Service is already registered")
