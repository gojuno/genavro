package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ErrorCode string

const (
	CODE_BAD_REQUEST ErrorCode = "ERR_BAD_REQUEST"
	CODE_MQ_TIMEOUT  ErrorCode = "ERR_MQ_TIMEOUT"
	CODE_INTERNAL    ErrorCode = "ERR_INTERNAL"
	CODE_FORBIDDEN   ErrorCode = "ERR_FORBIDDEN"
	CODE_DEADLINE    ErrorCode = "ERR_DEADLINE"
	CODE_NOT_FOUND   ErrorCode = "ERR_NOT_FOUND"

	// not returned to requesting side ever, see comments below
	CODE_REQUEST_SKIPPED ErrorCode = "ERR_REQUEST_SKIPPED"
)

var (
	BAD_REQUEST Error = New(CODE_BAD_REQUEST, http.StatusBadRequest)
	MQ_TIMEOUT  Error = New(CODE_MQ_TIMEOUT, http.StatusInternalServerError)
	INTERNAL    Error = New(CODE_INTERNAL, http.StatusInternalServerError)
	Canceled    Error = New("ERR_CANCELED", http.StatusInternalServerError)
	FORBIDDEN   Error = New(CODE_FORBIDDEN, http.StatusForbidden)
	TOO_LARGE   Error = New(CODE_FORBIDDEN, http.StatusRequestEntityTooLarge)
	DEADLINE    Error = New(CODE_DEADLINE, http.StatusInternalServerError)
	NOT_FOUND   Error = New(CODE_NOT_FOUND, http.StatusBadRequest)

	// This error indicates that nats async handler (e.g. subscribed with ms_core.Subscribe(..))
	// on this particular microservice instance decided not to handle given request, and thus no
	// reply should be sent from this handler (even if reply subject provided).
	//
	// This does not mean request won't be handled - it's probably supposed to be handled on
	// another instance (e.g. the one who tracks some info about particular user or smth like that),
	// and that instance would see request and return either nil or actual error - which would be then
	// sent to requester (if he provided reply subject - i.e. used RequestNATS, or skipped otherwise).
	//
	// So summarizing: we need this error as part of mechanism to acknowledge request processing from async handlers.
	//
	// Using 500 status code here since this error should never actually go via NATS or furthermore
	// HTTP APIs - so if it got there somehow it's definitely an internal error (developer mistake to be more precise).
	REQUEST_SKIPPED Error = New(CODE_REQUEST_SKIPPED, http.StatusInternalServerError)
)

type Error interface {
	error

	Code() ErrorCode
	HttpCode() int
	Data() json.RawMessage
	WithData(json.RawMessage) Error
	WithField(Field) Error
	WithCause(error) Error
	WithCausef(string, ...interface{}) Error
	Cause() string
	Equal(error) bool
}

type MSError struct {
	CodeValue     ErrorCode       `json:"code,omitempty"`
	HttpCodeValue int             `json:"http_code,omitempty"`
	RawData       json.RawMessage `json:"data,omitempty"`
	Fields        []Field         `json:"fields,omitempty"`
	RawError      error           `json:"-"`
}

type Field struct {
	Name    string `json:"field_name"`
	Message string `json:"message"`
}

var NoError = MSError{}

func New(code ErrorCode, httpCode int) Error {
	return MSError{CodeValue: code, HttpCodeValue: httpCode}
}

func (e MSError) Error() string {
	s := fmt.Sprintf("%s http(%d)", e.CodeValue, e.HttpCodeValue)
	if e.RawError != nil {
		s += fmt.Sprintf(" [%s]", e.RawError)
	}
	return s
}

func (e MSError) Code() ErrorCode {
	return e.CodeValue
}

func (e MSError) Cause() string {
	if e.RawError != nil {
		return e.RawError.Error()
	}

	return ""
}

func (e MSError) HttpCode() int {
	return e.HttpCodeValue
}

func (e MSError) Data() json.RawMessage {
	return e.RawData
}

func (e MSError) WithData(data json.RawMessage) Error {
	e.RawData = data
	return e
}

func (e MSError) WithField(f Field) Error {
	e.Fields = append(e.Fields, f)
	return e
}

func (e MSError) WithCause(err error) Error {
	e.RawError = err
	return e
}

func (e MSError) WithCausef(msg string, a ...interface{}) Error {
	return e.WithCause(fmt.Errorf(msg, a...))
}

func (e MSError) Equal(other error) bool {
	return Equal(other, e)
}

// Equal validates an error against codes from passed core errors
func Equal(err error, any ...Error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	for _, ee := range any {
		if ee.Code() == e.Code() {
			return true
		}
	}
	return false
}

// Cause returns a reason of a business error occuried. The result might be nil if no reason existed.
// In other words, the function returns RawError from MSError.
func Cause(err Error) error {
	var e error = err
	for e != nil {
		cause, ok := e.(MSError)
		if !ok {
			break
		}
		e = cause.RawError
	}
	return errors.Cause(e)
}

// TODO: Add more smart error processing/validating
// MarshalJSON - encoding/json Marshaler interface implementation
// UnmarshalJSON - encoding/json Marshaler interface implementation

func Expected(err Error) bool {
	switch err.HttpCode() {
	case http.StatusUnauthorized, http.StatusForbidden:
		return true
	case http.StatusBadRequest:
		// return true for business errors
		if err.Code() != CODE_BAD_REQUEST {
			return true
		}
		// return false for request validation errors
	}

	return false
}
