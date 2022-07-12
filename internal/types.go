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
	getSchedule(id string) (*schedule, *multi.Error)
	saveSchedule(s schedule) (string, *multi.Error)
	saveBatch(b batch, jobs []job) (id string, err *multi.Error)
	getBatch(id string) (*batch, *multi.Error)
	updateJob(j job) (err *multi.Error)
	saveJob(batchID string, j job) (err *multi.Error)
	getJob(jobID, batchID string) (*job, *multi.Error)
	getJobs(batchID string) ([]job, *multi.Error)
	close() *multi.Error
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
