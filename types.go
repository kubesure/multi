package multi

import "github.com/sirupsen/logrus"

// LogEvent stores log message
type LogEvent struct {
	id      EventCode
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

type Error struct {
	Code       EventCode
	Inner      error
	Message    ErrorMessage
	StackTrace string
	Misc       map[string]interface{}
}

type EventCode int

const (
	InternalError EventCode = iota + 1
	HTTPError
	BatchNotFoundError
	JobNotFoundError
)

type ErrorMessage string

const (
	DBError            ErrorMessage = "db transaction error"
	HTTPRequestError   ErrorMessage = "http request invalid"
	HTTPResponseError  ErrorMessage = "http response error"
	MultiInternalError ErrorMessage = "internal error"
)

const (
	BatchNotFound ErrorMessage = "batch not found"
	JobNotFound   ErrorMessage = "job not found"
	InputInvalid  ErrorMessage = "input invalid"
)

type ResponseErrors struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Code    EventCode    `json:"code"`
	Message ErrorMessage `json:"message"`
}
