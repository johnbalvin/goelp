package trace

import (
	"errors"
)

var (
	ErrEmpty      = errors.New("err Empty restult")
	ErrPermisions = errors.New("err Yelp requested captcha, try with another IP")
	ErrParameter  = errors.New("err not correct parameters")
	ErrMaxAttempt = errors.New("err Max attemps")
	ErrStatusCode = errors.New("err Not a correct status code")
)
