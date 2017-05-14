package engine

import "errors"

// Engine errors
var (
	ErrNoDocInit     = errors.New("No document initiliazer found")
	ErrNoConstructor = errors.New("No constructor")
	ErrNoClassFound  = errors.New("No class found")
)
