package multi

import (
	"os"

	"github.com/kubesure/multi"
	"github.com/kubesure/multi/internal"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func SaveBatch(batchType internal.BatchType, c *internal.CustomerSearch) (id string, err *multi.Error) {
	db, err := internal.NewDBConn(internal.SCHEDULAR, internal.SQLITE)
	defer db.Close()
	if err != nil {
		return "", err
	}

	b := internal.Batch{Type: batchType}
	jobs := []internal.Job{}
	id, err1 := db.SaveBatch(b, jobs)
	if err1 != nil {
		return "", err1
	}
	return "", nil
}

func GetBatch(batchId string) (*internal.Batch, *multi.Error) {
	db, err := internal.NewDBConn(internal.SCHEDULAR, internal.SQLITE)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	batch, dberr := db.GetBatch(batchId)
	if dberr != nil {
		return nil, dberr
	}
	return batch, dberr
}

func UpdateJob(j *internal.Job) *multi.Error {
	db, err := internal.NewDBConn(internal.SCHEDULAR, internal.SQLITE)
	defer db.Close()
	if err != nil {
		return err
	}
	dberr := db.UpdateJob(j)
	if dberr != nil {
		return dberr
	}
	return dberr
}
