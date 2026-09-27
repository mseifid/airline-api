package agency

import "errors"

var ErrInvalidDepositAmount = errors.New("deposit amount must be greater than zero")