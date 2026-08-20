package domain

import "errors"

var (
	ErrEmptyMatrix           = errors.New("matrix must contain at least one row")
	ErrInvalidMatrixShape    = errors.New("all matrix rows must have the same number of columns")
	ErrInvalidMatrixValue    = errors.New("matrix values must be finite numbers")
	ErrMatrixNotFactorizable = errors.New("matrix cannot be factorized")
	ErrStatisticsUnavailable = errors.New("statistics service unavailable")
	ErrInvalidRotationAngle  = errors.New("rotation angle must be 90, 180, or 270 degrees")
)
