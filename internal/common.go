package internal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kubesure/multi"
)

func PreChecks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ct := c.Request.Header.Get("Content-Type")
		if len(ct) == 0 || ct != "application/json" {
			c.AbortWithError(http.StatusBadRequest, errors.New("content type invalid"))
		}
	}
}

func BeforeResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		c.Writer.Header().Add("Content-Type", "application/json")
	}
}

func ResponseError(err multi.ErrorResponse, errs *multi.ResponseErrors) *multi.ResponseErrors {
	if errs == nil {
		errs = &multi.ResponseErrors{}
	}
	if len(errs.Errors) == 0 {
		errs.Errors = make([]multi.ErrorResponse, 0)
	}
	errs.Errors = append(errs.Errors, err)
	return errs
}

func UnmarshalAny[T any](bytes []byte) (*T, *multi.Error) {
	log := multi.NewLogger()
	out := new(T)
	if err := json.Unmarshal(bytes, out); err != nil {
		log.LogInternalError(err.Error())
		return nil, &multi.Error{Code: multi.HTTPError, Message: multi.MultiInternalError}
	}
	return out, nil
}

func MarshalAny[T any](t T) ([]byte, *multi.Error) {
	log := multi.NewLogger()
	data, err := json.Marshal(t)
	if err != nil {
		log.LogInternalError(err.Error())
		return nil, &multi.Error{Code: multi.HTTPError, Message: multi.MultiInternalError}
	}
	return data, nil
}
