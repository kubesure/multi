package internal

import (
	"time"

	"github.com/kubesure/multi"
)

type CustomerSearch struct {
	MaxResponseTimeSeconds uint       `json:"maxResponseTimeSeconds"`
	Customers              []Customer `json:"customers"`
}

type Customer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ScheduleResult struct {
	searchId string
}

type schedule struct {
	id, startDateTime, endDateTime string
	batch                          batch
}

type batch struct {
	id                               string
	createdDateTime, updatedDateTime time.Time
	ttype                            BatchType
	jobs                             []job
}

type job struct {
	id, maxResponse, retryInterval   uint
	batchId, payload, status         string
	errorMsg                         string
	retryCount                       uint
	createdDateTime, updatedDateTime time.Time
}

type database interface {
	getSchedule(id string) (*schedule, *multi.Error)
	saveSchedule(s schedule) (string, *multi.Error)
	saveBatch(b batch, jobs []job) (id string, err *multi.Error)
	getBatch(id string) (*batch, *multi.Error)
	saveJob(j job) (id string, err *multi.Error)
	getJob(id string) (*job, *multi.Error)
	getJobs(batchID string) ([]job, *multi.Error)
	close() *multi.Error
}

type BatchType int

const (
	CustomerSearchType BatchType = iota
)
