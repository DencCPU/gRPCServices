package sharederrors

import "errors"

// JWT errros
var (
	EXPIRED_TOKEN = errors.New("the token has expired")
)
