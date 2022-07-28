package multi

import "github.com/sirupsen/logrus"

//LogEvent stores log message
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
	InternalError EventCode = iota
	HTTPError
)

type ErrorMessage string

const (
	DBError          ErrorMessage = "DB Transaction Error"
	HTTPRequestError ErrorMessage = "HTTP Request Invalid"
)

type Erroresponse struct {
	Code    EventCode    `json:"errorCode"`
	Message ErrorMessage `json:"errorMessage"`
}
