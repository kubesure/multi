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
	batch                          Batch
}

type Batch struct {
	id                               string
	createdDateTime, updatedDateTime time.Time
	Type                             BatchType
	jobs                             []Job
}

type Job struct {
	Id              uint      `json:"id"`
	MaxResponse     uint      `json:"maxResponseTimeSeconds"`
	RetryInterval   uint      `json:"retryInterval"`
	BatchId         string    `json:"batchId"`
	Payload         string    `json:"payload"`
	Status          string    `json:"status"`
	ErrorMsg        string    `json:"errorMessage"`
	EndPoint        string    `json:"endPoimt"`
	RetryCount      uint      `json:"retryCount"`
	CreatedDateTime time.Time `json:"createdDateTime"`
	UpdatedDateTime time.Time `json:"updatedDateTime"`
}

type database interface {
	GetSchedule(id string) (*schedule, *multi.Error)
	SaveSchedule(s schedule) (string, *multi.Error)
	SaveBatch(b Batch, jobs []Job) (id string, err *multi.Error)
	GetBatch(id string) (*Batch, *multi.Error)
	UpdateJob(j Job) (err *multi.Error)
	SaveJob(j Job) (err *multi.Error)
	GetJob(jobID, batchID string) (*Job, *multi.Error)
	GetJobs(batchID string) ([]Job, *multi.Error)
	Close() *multi.Error
}

type BatchType int

const (
	CustomerSearchType BatchType = iota
)

type jobstatus string

const (
	CREATED   jobstatus = "CREATED"
	COMPLETED jobstatus = "COMPLETED"
)

type dbname string

const (
	SCHEDULAR  dbname = "../db/schedular.db"
	DISPATCHER dbname = "../db/dispatcher.db"
)

type dbtype string

const (
	SQLITE dbtype = "sqlite3"
)
