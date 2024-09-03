package entry

import "errors"

var (
	ErrEmailAlreadyEntered = errors.New("email has already been entered")
	ErrReCaptchaFiled      = errors.New("reCaptcha response not accepted")
	ErrFormMissing         = errors.New("form missing information")
)
