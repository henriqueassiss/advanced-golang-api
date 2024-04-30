package errorMsg

import "errors"

var (
	ErrTableIsPopulated   = errors.New("run-time: table is already populated")
	ErrInvalidRequestData = errors.New("run-time: invalid request data")
)
