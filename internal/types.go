package internal

import (
	"time"

	"github.com/kubesure/multi"
)

type schedule struct {
	id, startDateTime, endDateTime string
	batch                          Batch
}

type Batch struct {
	Id              string    `json:"id" binding:"required"`
	CreatedDateTime time.Time `json:"createdDateTime"`
	UpdatedDateTime time.Time `json:"updatedDateTime"`
	ttype           BatchType `json:"type"`
	jobs            []Job     `json:"jobs" binding:"required"`
}

type Job struct {
	Id                   string    `json:"id" binding:"required"`
	MaxResponseSeconds   uint      `json:"maxResponseTimeSeconds" binding:"required"`
	RetryIntervalSeconds uint      `json:"retryIntervalSeconds" binding:"required"`
	MaxRetry             uint      `json:"maxRetry" binding:"required"`
	BatchId              string    `json:"batchId" binding:"required"`
	Payload              *Payload  `json:"payload" binding:"required"`
	Result               *string   `json:"result"`
	Status               *string   `json:"status"`
	ErrorMsg             *string   `json:"errorMessage"`
	EndPoint             *Endpoint `json:"endPoint" binding:"required"`
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

type CreateBatchReq struct {
	Jobs []Job
}

type SaveJobsReq struct {
	BatchId string `json:"batchId"`
	Jobs    []Job
}

type Endpoint struct {
	BatchId string      `json:"batchId"`
	Uri     string      `json:"uri" binding:"required"`
	Method  string      `json:"method" binding:"required"`
	Auth    Auth        `json:"auth" binding:"required"`
	Headers []Header    `json:"headers" binding:"required"`
	Type    RequestType `json:"type"`
}

type Auth struct {
	Type              AuthType `json:"type"`
	UserName          string   `json:"userName"`
	Password          string   `json:"password"`
	ServerCertificate string   `json:"serverCertificate"`
}

type Header struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type database interface {
	GetSchedule(id string) (*schedule, *multi.Error)
	SaveSchedule(s schedule) (string, *multi.Error)
	SaveJobs(jobs []Job) (*Batch, *multi.Error)
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
