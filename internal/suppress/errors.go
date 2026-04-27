package suppress

import "errors"

// ErrInvalidWindow is returned when a non-positive suppression window is given.
var ErrInvalidWindow = errors.New("suppress: window must be greater than zero")
