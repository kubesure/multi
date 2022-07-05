package internal

import (
	"log"
	"testing"
)

func TestSaveBatch(t *testing.T) {
	db, err := newDBConn()
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	id, errsave := db.saveBatch(batchs(), jobs())
	if errsave != nil {
		log.Println(errsave)
		t.Errorf("should have saved")
	}

	if len(id) == 0 {
		t.Errorf("did not get id from save op")
	}
}

func jobs() []job {
	j1 := job{}
	j1.payload = "payload"
	j1.maxResponse = 5
	j1.retryInterval = 3

	j2 := job{}
	j2.payload = "payload"
	j2.maxResponse = 5
	j2.retryInterval = 3

	jobs := make([]job, 0)
	jobs = append(jobs, j1)
	jobs = append(jobs, j2)
	return jobs
}

func batchs() batch {
	return batch{ttype: CustomerSearchType}
}
