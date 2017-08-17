package errors

import "errors"

var (
	ErrApiParams = errors.New("Api params are invalid")
	ErrLimitPage = errors.New("Limit and Page must be greater than zero")

	ErrCantGetAccountsFromConfig = errors.New("Cannot get accounts from config")
)
