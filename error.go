package wechatpay

import (
	"encoding/json"
	"fmt"
)

// Error is detail error message
type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`

	err error
}

// NewError return a new error
func NewError(status int, message []byte) error {
	e := &Error{Status: status}
	if err := json.Unmarshal(message, e); err != nil {
		return NewInternalError(err)
	}

	return e
}

// NewInternalError return a internal error
func NewInternalError(err error) error {
	e := &Error{Code: "Internal", err: err}
	return e
}

// Error implement Error function for err
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.err != nil {
		return fmt.Sprintf("code: %s, message: %s", e.Code, e.err)
	}

	return fmt.Sprintf("staus: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
}
