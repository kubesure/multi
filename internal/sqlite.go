package internal

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kubesure/multi"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type sqlite struct {
	sqlite3 *sql.DB
}

func (db *sqlite) SaveBatch(b Batch, jobs []Job) (id string, err *multi.Error) {
	log := multi.NewLogger()
	tx, txerr := db.sqlite3.Begin()
	if txerr != nil {
		log.LogInternalError(txerr.Error())
		return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	stmt, serr := tx.Prepare("INSERT INTO batch(id,type,created_datetime,updated_datetime) VALUES (?,?,?,?)")
	if serr != nil {
		log.LogInternalError(serr.Error())
		return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	batchId := uuid.New().String()

	_, err1 := stmt.Exec(batchId, b.Type, currentDateTime(), currentDateTime())
	defer stmt.Close()
	if err1 != nil {
		log.LogInternalError(err1.Error())
		return "", err
	}

	insertSQL := "INSERT INTO job(id,batch_id,req_payload,endpoint,status,error_msg,max_response,retry_interval,retry_count,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	jobStmt, joberr := tx.Prepare(insertSQL)
	if joberr != nil {
		return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	for index, job := range jobs {
		_, joberr = jobStmt.Exec(index+1, batchId, job.Payload, job.EndPoint, CREATED, job.ErrorMsg, job.MaxResponse, job.RetryInterval, job.RetryCount, currentDateTime(), currentDateTime())
		if joberr != nil {
			tx.Rollback()
			log.LogInternalError(joberr.Error())
			return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
		}
	}

	txerr = tx.Commit()
	if txerr != nil {
		log.LogInternalError(txerr.Error())
		return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	return batchId, nil
}

func (db *sqlite) GetBatch(id string) (*Batch, *multi.Error) {
	log := multi.NewLogger()
	row, qerr := db.sqlite3.Query("select id,type, created_datetime, updated_datetime from batch where id=?", id)

	if qerr != nil {
		log.LogInternalError(qerr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	defer row.Close()
	var b *Batch
	for row.Next() {
		b = &Batch{}
		var created_datetime, updated_datetime string
		row.Scan(&b.id, &b.Type, &created_datetime, &updated_datetime)
		b.createdDateTime = *parseDateTime(created_datetime)
		b.updatedDateTime = *parseDateTime(updated_datetime)
		jobs, jerr := db.GetJobs(id)
		if jerr != nil {
			return nil, jerr
		}
		b.jobs = jobs
	}

	return b, nil
}

func (db *sqlite) GetJobs(batchID string) ([]Job, *multi.Error) {
	log := multi.NewLogger()
	q := "select id,batch_id,payload,result,endpoint,status,created_datetime,updated_datetime,max_response,retry_interval,error_msg,retry_count from job where batch_id=?"
	jobs, err := db.Getjobs(q, batchID)
	if err != nil {
		log.LogInternalError(err.Inner.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	return jobs, nil
}

func (db *sqlite) GetJob(jobID, batchID string) (*Job, *multi.Error) {
	q := "select id,batch_id,payload,result,endpoint,status,created_datetime,updated_datetime,max_response,retry_interval,error_msg,retry_count from job where batch_id=? and id = ?"
	log := multi.NewLogger()
	jobs, err := db.Getjobs(q, batchID, jobID)

	if err != nil {
		log.LogInternalError(err.Inner.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	if len(jobs) == 1 {
		return &jobs[0], nil
	}
	return nil, nil
}

func (db *sqlite) Getjobs(query string, id ...string) ([]Job, *multi.Error) {
	log := multi.NewLogger()
	var qerr error
	var rows *sql.Rows
	if len(id) > 1 {
		rows, qerr = db.sqlite3.Query(query, id[0], id[1])
	} else {
		rows, qerr = db.sqlite3.Query(query, id[0])
	}

	if qerr != nil {
		log.LogInternalError(qerr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError, Inner: qerr}
	}

	defer rows.Close()
	jobs := make([]Job, 0)
	for rows.Next() {
		j := Job{}
		var created_datetime, updated_datetime string
		var errormsg sql.NullString
		var retrycount sql.NullInt32
		rows.Scan(&j.Id, &j.BatchId, &j.Payload, &j.Result, &j.EndPoint, &j.Status, &created_datetime,
			&updated_datetime, &j.MaxResponse, &j.RetryInterval, &errormsg, &retrycount)
		if errormsg.Valid {
			j.ErrorMsg = &errormsg.String
		} else {
			val := ""
			j.ErrorMsg = &val
		}

		if retrycount.Valid {
			count := uint(retrycount.Int32)
			j.RetryCount = &count
		} else {
			var count uint = 0
			j.RetryCount = &count
		}
		j.CreatedDateTime = *parseDateTime(created_datetime)
		j.UpdatedDateTime = *parseDateTime(updated_datetime)
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (db *sqlite) SaveJob(j *Job) (err *multi.Error) {
	log := multi.NewLogger()
	insertSQL := "INSERT INTO job(id,batch_id,payload,result,endpoint,status,error_msg,max_response,retry_interval,retry_count,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	jobStmt, joberr := db.sqlite3.Prepare(insertSQL)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	_, joberr = jobStmt.Exec(j.Id, j.BatchId, j.Payload, j.Result, j.EndPoint, j.Status, j.ErrorMsg, j.MaxResponse, j.RetryInterval, j.RetryCount, currentDateTime(), currentDateTime())
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	return nil
}

func (db *sqlite) UpdateJob(j Job) (err *multi.Error) {
	log := multi.NewLogger()
	insertSQL := "UPDATE job SET status = coalesce(?,status), error_msg = coalesce(?,error_msg),retry_count = coalesce(?,retry_count) " +
		",updated_datetime = ?, result = coalesce(?,result) where id = ? and batch_id= ?"
	jobStmt, joberr := db.sqlite3.Prepare(insertSQL)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	_, joberr = jobStmt.Exec(j.Status, j.ErrorMsg, j.RetryCount, currentDateTime(), j.Result, j.Id, j.BatchId)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	return nil
}

func (db *sqlite) GetSchedule(id string) (*schedule, *multi.Error) {
	return nil, nil
}

func (db *sqlite) SaveSchedule(s schedule) (id string, err *multi.Error) {
	return "", nil
}

func (db *sqlite) Close() *multi.Error {
	log := multi.NewLogger()
	err := db.sqlite3.Close()
	if err != nil {
		log.LogInternalError(err.Error())
	}
	return &multi.Error{Code: multi.InternalError, Message: "Error while closing db connection"}
}

func NewDBConn(name dbname, dbtype dbtype) (database, *multi.Error) {
	db, err := sql.Open(string(dbtype), string(name))
	if err != nil {
		log.Println(err)
		return nil, &multi.Error{Code: multi.InternalError, Message: "Error opening schedular database"}
	}
	return &sqlite{sqlite3: db}, nil
}

func currentDateTime() string {
	return time.Now().Format(time.RFC3339)
}

func parseDateTime(dt string) *time.Time {
	t, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return nil
	}
	return &t
}
