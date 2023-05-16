package internal

import (
	"time"

	"github.com/kubesure/multi"
)

type CustomerSearch struct {
	MaxResponseTimeSeconds uint       `json:"maxResponseTimeSeconds"`
	EndPoint               Endpoint   `json:"endPoint"`
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
	Id              string    `json:"id"`
	CreatedDateTime time.Time `json:"createdDateTime"`
	UpdatedDateTime time.Time `json:"updatedDateTime"`
	ttype           BatchType `json:"type"`
	jobs            []Job     `json:"jobs"`
}

type Job struct {
	Id                   string    `json:"id"`
	MaxResponseSeconds   uint      `json:"maxResponseTimeSeconds"`
	RetryIntervalSeconds uint      `json:"retryIntervalSeconds"`
	MaxRetry             uint      `json:"maxRetry"`
	BatchId              string    `json:"batchId"`
	Payload              *Payload  `json:"payload"`
	Result               *string   `json:"result"`
	Status               *string   `json:"status"`
	ErrorMsg             *string   `json:"errorMessage"`
	EndPoint             *Endpoint `json:"endPoint"`
	CreatedDateTime      time.Time `json:"createdDateTime"`
	UpdatedDateTime      time.Time `json:"updatedDateTime"`
}

type Payload struct {
	CompressedDispatch bool
	Data               string
}

type AuthType string

const (
	BASIC  AuthType = "BASIC"
	MTLS   AuthType = "MTLS"
	BEARER AuthType = "BEARER"
)

type RequestType string

const (
	HTTP RequestType = "HTTP"
	GRPC RequestType = "GRPC"
)

type Endpoint struct {
	BatchId string      `json:"batchId"`
	Uri     string      `json:"uri"`
	Method  string      `json:"method"`
	Auth    Auth        `json:"auth"`
	Headers []Header    `json:"headers"`
	Type    RequestType `json:"type"`
}

type Auth struct {
	Type              AuthType `json:"type"`
	UserName          string   `json:"userName"`
	Password          string   `json:"password"`
	ServerCertificate string   `json:"serverCertificate"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type database interface {
	GetSchedule(id string) (*schedule, *multi.Error)
	SaveSchedule(s schedule) (string, *multi.Error)
	SaveBatch(jobs []Job) (*Batch, *multi.Error)
	GetBatch(id string) (*Batch, *multi.Error)
	UpdateJob(j *Job) *multi.Error
	SaveJob(j *Job) *multi.Error
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
	SCHEDULAR  dbname = "../../db/schedular.db"
	DISPATCHER dbname = "../../db/dispatcher.db"
)

type dbtype string

const (
	SQLITE dbtype = "sqlite3"
)
