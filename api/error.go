package api

import "strconv"

type apiError struct {
	Error   error
	Message string
	Code    int
	Request string
}

func (apiError apiError) Formatted() string {
	return apiError.Error.Error() + ";" + apiError.Message + ";" + strconv.Itoa(apiError.Code) + ";" + apiError.Request
}
