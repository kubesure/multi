package internal

import (
	"os"

	"github.com/kubesure/multi"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func SaveBatch(jobs []Job) (batch *Batch, err *multi.Error) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	b, err1 := db.SaveBatch(jobs)
	if err1 != nil {
		return nil, err1
	}
	return b, nil
}

func GetBatch(batchId string) (*Batch, *multi.Error) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	batch, dberr := db.GetBatch(batchId)
	if dberr != nil {
		return nil, dberr
	}
	return batch, nil
}

func UpdateJob(j *Job) *multi.Error {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	defer db.Close()
	if err != nil {
		return err
	}
	dberr := db.UpdateJob(j)
	if dberr != nil {
		return dberr
	}
	return nil
}

func GetJob(jobID string, batchID string) (*Job, *multi.Error) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	jobs, dberr := db.GetJob(jobID, batchID)
	if dberr != nil {
		return nil, dberr
	}
	return jobs, nil
}

func GetJobs(batchID string) ([]Job, *multi.Error) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	jobs, dberr := db.GetJobs(batchID)
	if dberr != nil {
		return nil, dberr
	}
	return jobs, nil
}
