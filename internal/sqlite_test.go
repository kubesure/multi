package internal

import (
	"testing"
)

func TestSaveBatch(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	id, errsave := db.SaveBatch(batchs(), jobs())
	if errsave != nil {
		t.Errorf("should have saved batch ")
	}

	if len(id) == 0 {
		t.Errorf("did not get id from save op")
	}
}

func TestGetBatchFound(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, err := db.GetBatch("45a515a1-9f8b-45ca-aad1-a81e11108a68")

	if err != nil {
		t.Errorf("Should have reterived batch")
	}

	if b.id != "45a515a1-9f8b-45ca-aad1-a81e11108a68" {
		t.Errorf("cannot reterive id")
	}
}

func TestGetBatchNotFound(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, err := db.GetBatch("30c42b49-98d3-47b0-969c-65754f9052a")

	if err != nil {
		t.Errorf("Should not be an error it should be nil batch")
	}

	if b != nil {
		t.Errorf("Record should not have been found")
	}
}

func TestGetJobs(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	jobs, err := db.GetJobs("45a515a1-9f8b-45ca-aad1-a81e11108a68")

	if err != nil {
		t.Errorf("Should have reterived jobs")
	}

	if len(jobs) != 2 {
		t.Errorf("should have reterived two jobs")
	}
}

func TestGetJob(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	job, err := db.GetJob("5", "45a515a1-9f8b-45ca-aad1-a81e11108a68")

	if err != nil {
		t.Errorf("Should have reterived job")
	}

	if job.Id != 5 {
		t.Errorf("should have reterived job id 1")
	}

}

func TestSaveJob(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}

	errsave := db.SaveJob(&jobs()[0])
	if errsave != nil {
		t.Errorf("should have saved job")
	}
}

func TestUpdateJob(t *testing.T) {
	db, err := NewDBConn(SCHEDULAR, SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}

	j := jobUpdate()
	status := string(COMPLETED)
	j.Status = &status
	j.BatchId = "45a515a1-9f8b-45ca-aad1-a81e11108a68"
	var count uint = 17
	j.RetryCount = &count
	errmsg := "conn error"
	j.ErrorMsg = &errmsg
	res := "result"
	j.Result = &res
	errupdate := db.UpdateJob(&j)
	if errupdate != nil {
		t.Errorf("should have updated job")
	}
}

func jobUpdate() Job {
	j1 := Job{}
	j1.Id = 5
	j1.BatchId = "45a515a1-9f8b-45ca-aad1-a81e11108a68"
	j1.Payload = "payload"
	j1.EndPoint = "http://localhost/customer/search"
	j1.MaxResponse = 5
	j1.RetryInterval = 3
	msg := "error msg"
	j1.ErrorMsg = &msg
	var count uint = 10
	j1.RetryCount = &count
	//j1.Status = string(CREATED)
	return j1
}

func jobs() []Job {
	j1 := Job{}
	j1.Id = 5
	j1.BatchId = "45a515a1-9f8b-45ca-aad1-a81e11108a68"
	j1.Payload = "payload"
	j1.EndPoint = "http://localhost/customer/search"
	j1.MaxResponse = 5
	j1.RetryInterval = 3
	msg := "error msg"
	j1.ErrorMsg = &msg
	var count uint = 10
	j1.RetryCount = &count
	status := string(CREATED)
	j1.Status = &status

	j2 := Job{}
	j2.Payload = "payload"
	j2.EndPoint = "http://localhost/customer/search"
	j2.MaxResponse = 5
	j2.RetryInterval = 3
	msg = "error msg"
	j2.ErrorMsg = &msg
	j2.RetryCount = &count
	status = string(CREATED)
	j2.Status = &status

	jobs := make([]Job, 0)
	jobs = append(jobs, j1)
	jobs = append(jobs, j2)
	return jobs
}

func batchs() Batch {
	return Batch{Type: CustomerSearchType}
}
