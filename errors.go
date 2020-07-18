package mergo

import "errors"

var (
	// ErrNilDestination happens when trying to merge into a nil value.
	ErrNilDestination = errors.New("unable to merge into nil destination")
	// ErrNilSource happens when trying to merge from a nil value.
	ErrNilSource = errors.New("unable to merge a nil source")
	// ErrNonPointerDestination happens when destination is not a pointer.
	ErrNonPointerDestination = errors.New("destination must be a pointer")
	// ErrIncompatibleTypes happens when dst and src are not compatible.
	ErrIncompatibleTypes = errors.New("destination and source are not compatible")
)
