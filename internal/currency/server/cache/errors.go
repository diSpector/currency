package cache

import "github.com/pkg/errors"

var (
	ErrNotFound      = errors.New(`not found in cache`)
	ErrFoundInAbsent = errors.New(`found in absent`)
)
