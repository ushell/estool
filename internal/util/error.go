package util

import "encoding/json"

type Err struct {
	Code int
	Message string
}

func (e *Err) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func NewError(code int, message string) error {
	return &Err{
		Code: code,
		Message: message,
	}
}