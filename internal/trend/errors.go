package trend

import "errors"

// ErrInvalidWindow is returned when a non-positive window duration is provided.
var ErrInvalidWindow = errors.New("trend: window duration must be positive")
