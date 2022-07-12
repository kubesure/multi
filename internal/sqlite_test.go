package internal

import (
	"testing"
)

func TestSaveBatch(t *testing.T) {
	db, err := newDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	id, errsave := db.saveBatch(batchs(), jobs())
	if errsave != nil {
		t.Errorf("should have saved batch ")
	}

	if len(id) == 0 {
		t.Errorf("did not get id from save op")
	}
}

func TestGetBatchFound(t *testing.T) {
	db, err := newDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, err := db.getBatch("45a515a1-9f8b-45ca-aad1-a81e11108a68")

	if err != nil {
		t.Errorf("Should have reterived batch")
	}

	if b.id != "45a515a1-9f8b-45ca-aad1-a81e11108a68" {
		t.Errorf("cannot reterive id")
	}
}

func TestGetBatchNotFound(t *testing.T) {
	db, err := newDBConn(SCHEDULAR, SQLITE)
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
	db, err := newDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	jobs, err := db.getJobs("45a515a1-9f8b-45ca-aad1-a81e11108a68")

	if err != nil {
		t.Errorf("Should have reterived jobs")
	}

	if len(jobs) != 2 {
		t.Errorf("should have reterived two jobs")
	}
}

func TestGetJob(t *testing.T) {
	db, err := newDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	job, err := db.getJob("1", "45a515a1-9f8b-45ca-aad1-a81e11108a68")

	if err != nil {
		t.Errorf("Should have reterived job")
	}

	if job.Id != 1 {
		t.Errorf("should have reterived job id 1")
	}

}

func TestSaveJob(t *testing.T) {
	db, err := newDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}

	errsave := db.saveJob("45a515a1-9f8b-45ca-aad1-a81e11108a68", jobs()[0])
	if errsave != nil {
		t.Errorf("should have saved job")
	}
}

func TestUpdateJob(t *testing.T) {
	db, err := newDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}

	j := jobs()[0]
	j.Status = string(COMPLETED)
	j.BatchId = "45a515a1-9f8b-45ca-aad1-a81e11108a68"
	errupdate := db.updateJob(j)
	if errupdate != nil {
		t.Errorf("should have updated job")
	}
}

func jobs() []job {
	j1 := job{}
	j1.Id = 4
	j1.Payload = "payload"
	j1.EndPoint = "http://localhost/customer/search"
	j1.MaxResponse = 5
	j1.RetryInterval = 3
	j1.ErrorMsg = "error msg"
	j1.RetryCount = 10
	j1.Status = string(CREATED)

	j2 := job{}
	j2.Payload = "payload"
	j2.EndPoint = "http://localhost/customer/search"
	j2.MaxResponse = 5
	j2.RetryInterval = 3
	j2.ErrorMsg = "error msg"
	j2.RetryCount = 10
	j2.Status = string(CREATED)

	jobs := make([]job, 0)
	jobs = append(jobs, j1)
	jobs = append(jobs, j2)
	return jobs
}

func batchs() batch {
	return batch{ttype: CustomerSearchType}
}
