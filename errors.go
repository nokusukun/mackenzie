package mackenzie

import (
	"errors"
	"fmt"
)

var (
	ErrMackenzie                               = errors.New("mackenzie error")
	ErrCallerMustBeFunction                    = fmt.Errorf("%w: caller must be a function", ErrMackenzie)
	ErrCallerMustHaveAtLeastOneArgument        = fmt.Errorf("%w: caller must have at least one argument", ErrMackenzie)
	ErrCallerMustHaveAtLeastOneReturnValue     = fmt.Errorf("%w: caller must have at least one return value", ErrMackenzie)
	ErrCallerMustHaveNoMoreThanTwoReturnValues = fmt.Errorf("%w: caller must have no more than two return values", ErrMackenzie)
	ErrCallerMustReturnTAsItsFirstMethod       = fmt.Errorf("%w: caller must return T as it's first method", ErrMackenzie)
	ErrCallerMustReturnAnErrorAsItsLastMethod  = fmt.Errorf("%w: caller must return an error as it's last method", ErrMackenzie)
	ErrCallInvalid                             = fmt.Errorf("%w: call error", ErrMackenzie)
	ErrIncorrectNumberOfArguments              = fmt.Errorf("%w: incorrect number of arguments", ErrCallInvalid)
	ErrIncorrectTypeForArgument                = func(expected interface{}, actual interface{}) error {
		return fmt.Errorf("%w: incorrect type for argument, expected %T, got %T", ErrCallInvalid, expected, actual)
	}
)
