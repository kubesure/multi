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

type EventCode int32

const (
	InternalError EventCode = iota
)

type ErrorMessage string

const (
	DBError ErrorMessage = "DB Transaction Error"
)
