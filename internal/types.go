package internal

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
	batchID, createdDate, updatedDate string
	ttype                             BatchType
}

type job struct {
	id, maxResponse, retryInterval, retryCount                           uint
	batchId, payload, status, errorMsg, createdDateTime, updatedDateTime string
}

type database interface {
	getSchedule(id string) *schedule
	saveSchedule(s schedule) (id string)
	saveBatch(b batch, jobs []job) (id string, err error)
	getBatch(id string) *batch
	//saveJob(j job) (id string)
	getJob(id string) *job
	getJobs(batchID string) []job
	close() error
}

type BatchType int

const (
	CustomerSearchType BatchType = iota
)
