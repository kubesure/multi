package internal

import (
	"testing"
)

func TestSaveBatch(t *testing.T) {
	db, err := newDBConn()
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	id, errsave := db.saveBatch(batchs(), jobs())
	if errsave != nil {
		t.Errorf("should have saved")
	}

	if len(id) == 0 {
		t.Errorf("did not get id from save op")
	}
}

func TestGetBatchFound(t *testing.T) {
	db, err := newDBConn()
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, err := db.getBatch("30c42b49-98d3-47b0-969c-65754f9052a2")

	if err != nil {
		t.Errorf("Should have reterived batch")
	}

	if b.id != "30c42b49-98d3-47b0-969c-65754f9052a2" {
		t.Errorf("cannot reterive id")
	}
}

func TestGetBatchNotFound(t *testing.T) {
	db, err := newDBConn()
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, err := db.getBatch("30c42b49-98d3-47b0-969c-65754f9052a")

	if err != nil {
		t.Errorf("Should not be an error it should be nil batch")
	}

	if b != nil {
		t.Errorf("Record should not have been found")
	}
}

func TestGetJobs(t *testing.T) {
	db, err := newDBConn()
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	jobs, err := db.getJobs("30c42b49-98d3-47b0-969c-65754f9052a2")

	if err != nil {
		t.Errorf("Should have reterived jobs")
	}

	if len(jobs) != 2 {
		t.Errorf("should have reterived two jobs")
	}
}

func jobs() []job {
	j1 := job{}
	j1.payload = "payload"
	j1.maxResponse = 5
	j1.retryInterval = 3
	j1.errorMsg = "error msg"
	j1.retryCount = 10

	j2 := job{}
	j2.payload = "payload"
	j2.maxResponse = 5
	j2.retryInterval = 3
	j2.errorMsg = "error msg"
	j2.retryCount = 10

	jobs := make([]job, 0)
	jobs = append(jobs, j1)
	jobs = append(jobs, j2)
	return jobs
}

func batchs() batch {
	return batch{ttype: CustomerSearchType}
}
