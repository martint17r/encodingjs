package encodingjs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/robertkrimen/otto"
)

var (
	// ErrArrayExpected is returned when the JS variable is not an array (i.e. for a go slice)
	ErrArrayExpected = errors.New("array expected")
	// ErrObjectExpected is returned when the JS variable is not an object (i.e. for a go map or struct)
	ErrObjectExpected = errors.New("object expected")
)

// UnsupportedTypeError is returned when channels and functions are given as targets
type UnsupportedTypeError struct {
	rtype reflect.Type
}

func (ute *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unknown type: %s", ute.rtype)
}

// InvalidValueError is returned when the js variable does not match the requirements from go
type InvalidValueError struct {
	Wanted string
	Got    otto.Value
}

func (ive *InvalidValueError) Error() string {
	return fmt.Sprintf("invalid value (wanted %s): %s", ive.Wanted, ive.Got)
}
