package internal

import (
	"encoding/json"

	"github.com/kubesure/multi"
)

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
