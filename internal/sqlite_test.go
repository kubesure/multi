package internal

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestSaveBatch(t *testing.T) {
	db, err := NewDBConn("../db/schedular.db", SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, errsave := db.SaveJobs(job())
	if errsave != nil {
		t.Errorf("should have saved batch ")
	}

	if b == nil {
		t.Errorf("did not get id from save op")
	}
}

func TestGetBatchFound(t *testing.T) {
	db, err := NewDBConn("../db/schedular.db", SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	b, err := db.GetBatch("0dc73091-e790-488a-98ef-28cdfdbeba3c")

	if err != nil {
		t.Errorf("Should have reterived batch")
	}

	if b.Id != "0dc73091-e790-488a-98ef-28cdfdbeba3c" {
		t.Errorf("cannot reterive id")
	}
}

func TestGetBatchNotFound(t *testing.T) {
	db, err := NewDBConn("../db/schedular.db", SQLITE)
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
	db, err := NewDBConn("../db/schedular.db", SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	jobs, err := db.GetJobs("05ddeaec-3b43-4a35-8e77-1196c1d3c61c")

	if err != nil {
		t.Errorf("Should have reterived jobs")
	}

	if len(jobs) != 1 {
		t.Errorf("should have reterived two jobs")
	}
}

func TestGetJob(t *testing.T) {
	db, err := NewDBConn("../db/schedular.db", SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}
	job, err := db.GetJob("102234456242", "05ddeaec-3b43-4a35-8e77-1196c1d3c61c")

	if err != nil {
		t.Errorf("Should have reterived job")
	} else {
		if job.Payload.Data != "{data}" {
			t.Errorf("should have reterived data %v", "{data}")
		}
	}

}

func TestSaveJob(t *testing.T) {
	db, err := NewDBConn("../db/schedular.db", SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}

	errsave := db.SaveJob(&jobs()[0])
	if errsave != nil {
		t.Errorf("should have saved job")
	}
}

func TestUpdateJob(t *testing.T) {
	db, err := NewDBConn("../db/schedular.db", SQLITE)
	if err != nil {
		t.Errorf("should have had not got a db conn error")
	}

	j := jobUpdate()
	status := string(COMPLETED)
	j.Status = &status
	//j.BatchId = "05ddeaec-3b43-4a35-8e77-1196c1d3c61c"
	var count uint = 17
	j.MaxRetry = count
	errmsg := "conn error"
	j.ErrorMsg = &errmsg
	res := "{result}"
	j.Result = &res
	errupdate := db.UpdateJob(&j)
	if errupdate != nil {
		t.Errorf("should have updated job")
	}
}

func jobUpdate() Job {
	j1 := Job{}
	j1.Id = "102234456242"
	j1.BatchId = "05ddeaec-3b43-4a35-8e77-1196c1d3c61c"
	j1.Payload = txtPayload()
	j1.EndPoint = endPoint()
	j1.MaxResponseSeconds = 5
	j1.RetryIntervalSeconds = 3
	msg := "error msg new"
	j1.ErrorMsg = &msg
	var count uint = 10
	j1.MaxRetry = count
	//j1.Status = string(CREATED)
	return j1
}

func txtPayload() *Payload {
	return &Payload{CompressedDispatch: true, Data: "{data}"}
}

func job() []Job {
	j1 := Job{}
	j1.Id = "992234456242"
	j1.BatchId = "45a515a1-9f8b-45ca-aad1-a81e99908a68"
	j1.Payload = txtPayload()
	j1.EndPoint = endPoint()
	j1.MaxResponseSeconds = 5
	j1.RetryIntervalSeconds = 3
	msg := "error msg 2"
	j1.ErrorMsg = &msg
	var count uint = 10
	j1.MaxRetry = count
	status := string(CREATED)
	j1.Status = &status
	jobs := []Job{}
	jobs = append(jobs, j1)
	return jobs
}

func jobs() []Job {
	j1 := Job{}
	j1.Id = "1023424234"
	j1.BatchId = "45a515a1-9f8b-45ca-aad1-a81e11108a68"
	j1.Payload = txtPayload()
	j1.EndPoint = endPoint()
	j1.MaxResponseSeconds = 5
	j1.RetryIntervalSeconds = 3
	msg := "error msg"
	j1.ErrorMsg = &msg
	var count uint = 10
	j1.MaxRetry = count
	status := string(CREATED)
	j1.Status = &status

	j2 := Job{}
	j2.Id = "12674633"
	j2.Payload = txtPayload()
	j2.EndPoint = endPoint()
	j2.MaxResponseSeconds = 5
	j2.RetryIntervalSeconds = 3
	msg = "error msg"
	j2.ErrorMsg = &msg
	j2.MaxRetry = count
	status = string(CREATED)
	j2.Status = &status

	jobs := make([]Job, 0)
	jobs = append(jobs, j1)
	jobs = append(jobs, j2)
	return jobs
}

func endPoint() *Endpoint {
	ep := Endpoint{}
	ep.Method = "GET"
	ep.Type = "HTTP"
	ep.Uri = "http://localhost/customer/search"
	ep.Auth.Type = BASIC
	ep.Auth.Password = "pass"
	ep.Auth.UserName = "user"
	ep.Auth.ServerCertificate = "sdlfkas;dgjads;gokads;fke0rsdl;gkafg9eriglkdbjadijgdlkg"
	headers := make([]Header, 0)
	h := Header{Key: "Bearer", Value: "sfsdfas345wrsfsfd"}
	headers = append(headers, h)
	ep.Headers = headers
	return &ep
}

func TestUUID(t *testing.T) {
	u := uuid.New()
	fmt.Printf("%v \n", u)
}
