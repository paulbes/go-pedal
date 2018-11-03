package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// nolint
const (
	NotFound int = iota
	Unmarshal
	Marshal
	IO
)

type errors struct {
	err error
	msg string
	typ int
}

// New creates a new error with the provided information
func New(err error, msg string, typ int) error {
	return &errors{
		err: err,
		msg: msg,
		typ: typ,
	}
}

// codeString converts the int const
// to a string representation
func (e *errors) codeString() string {
	var typ string
	switch e.typ {
	case NotFound:
		typ = "notfound"
	case Marshal:
		typ = "marshal"
	case Unmarshal:
		typ = "unmarshal"
	case IO:
		typ = "io"
	default:
		typ = "unknown"
	}
	return typ
}

// Error implements the error interface
func (e *errors) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.codeString(), e.msg, e.err)
}

// StatusCode implements the go-kit StatusCoder interface
// so that the type of error can be converted to its
// corresponding http error code
func (e *errors) StatusCode() int {
	var code int
	switch e.typ {
	case NotFound:
		code = http.StatusNotFound
	case Unmarshal:
		code = http.StatusBadRequest
	case Marshal:
		fallthrough
	case IO:
		fallthrough
	default:
		code = http.StatusInternalServerError
	}
	return code
}

// MarshalJSON implements the json.Marshaller interface
// so that the error can be marshalled
func (e *errors) MarshalJSON() ([]byte, error) {
	content := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Type    string `json:"type"`
	}{
		Message: e.Error(),
		Code:    e.StatusCode(),
		Type:    e.codeString(),
	}
	data, err := json.Marshal(&content)
	if err != nil {
		return nil, err
	}
	return data, nil
}
