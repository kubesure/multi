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

func (db *sqlite) saveBatch(b batch, jobs []job) (id string, err *multi.Error) {
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

	_, err1 := stmt.Exec(batchId, b.ttype, currentDateTime(), currentDateTime())
	defer stmt.Close()
	if err1 != nil {
		log.LogInternalError(err1.Error())
		return "", err
	}

	insertSQL := "INSERT INTO job(id,batch_id,payload,endpoint,status,error_msg,max_response,retry_interval,retry_count,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	jobStmt, joberr := tx.Prepare(insertSQL)
	if joberr != nil {
		return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	for index, job := range jobs {
		_, joberr = jobStmt.Exec(index+1, batchId, job.payload, job.endPoint, CREATED, job.errorMsg, job.maxResponse, job.retryInterval, job.retryCount, currentDateTime(), currentDateTime())
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

func (db *sqlite) getBatch(id string) (*batch, *multi.Error) {
	log := multi.NewLogger()
	row, qerr := db.sqlite3.Query("select id,type, created_datetime, updated_datetime from batch where id=?", id)

	if qerr != nil {
		log.LogInternalError(qerr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	defer row.Close()
	var b *batch
	for row.Next() {
		b = &batch{}
		var created_datetime, updated_datetime string
		row.Scan(&b.id, &b.ttype, &created_datetime, &updated_datetime)
		b.createdDateTime = *parseDateTime(created_datetime)
		b.updatedDateTime = *parseDateTime(updated_datetime)
		jobs, jerr := db.getJobs(id)
		if jerr != nil {
			return nil, jerr
		}
		b.jobs = jobs
	}

	return b, nil
}

func (db *sqlite) getJobs(batchID string) ([]job, *multi.Error) {
	log := multi.NewLogger()
	q := "select id,batch_id,payload,endpoint,status,created_datetime,updated_datetime,max_response,retry_interval,error_msg,retry_count from job where batch_id=?"
	jobs, err := db.getjobs(q, batchID)
	if err != nil {
		log.LogInternalError(err.Inner.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	return jobs, nil
}

func (db *sqlite) getJob(jobID, batchID string) (*job, *multi.Error) {
	q := "select id,batch_id,payload,endpoint,status,created_datetime,updated_datetime,max_response,retry_interval,error_msg,retry_count from job where batch_id=? and id = ?"
	log := multi.NewLogger()
	jobs, err := db.getjobs(q, batchID, jobID)

	if err != nil {
		log.LogInternalError(err.Inner.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	if len(jobs) == 1 {
		return &jobs[0], nil
	}
	return nil, nil
}

func (db *sqlite) getjobs(query string, id ...string) ([]job, *multi.Error) {
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
	jobs := make([]job, 0)
	for rows.Next() {
		j := job{}
		var created_datetime, updated_datetime string
		var errormsg sql.NullString
		var retrycount sql.NullInt32
		rows.Scan(&j.id, &j.batchId, &j.payload, &j.endPoint, &j.status, &created_datetime,
			&updated_datetime, &j.maxResponse, &j.retryInterval, &errormsg, &retrycount)
		if errormsg.Valid {
			j.errorMsg = errormsg.String
		} else {
			j.errorMsg = ""
		}

		if retrycount.Valid {
			j.retryCount = uint(retrycount.Int32)
		} else {
			j.retryCount = 0
		}
		j.createdDateTime = *parseDateTime(created_datetime)
		j.updatedDateTime = *parseDateTime(updated_datetime)
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (db *sqlite) saveJob(batchId string, j job) (err *multi.Error) {
	log := multi.NewLogger()
	insertSQL := "INSERT INTO job(id,batch_id,payload,endpoint,status,error_msg,max_response,retry_interval,retry_count,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	jobStmt, joberr := db.sqlite3.Prepare(insertSQL)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	_, joberr = jobStmt.Exec(j.id, batchId, j.payload, j.endPoint, j.status, j.errorMsg, j.maxResponse, j.retryInterval, j.retryCount, currentDateTime(), currentDateTime())
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	return nil
}

func (db *sqlite) updateJob(j job) (err *multi.Error) {
	log := multi.NewLogger()
	insertSQL := "UPDATE job SET status = ?, error_msg = ?,retry_count = ? ,updated_datetime = ? where id = ? and batch_id= ?"
	jobStmt, joberr := db.sqlite3.Prepare(insertSQL)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	_, joberr = jobStmt.Exec(j.status, j.errorMsg, j.retryCount, currentDateTime(), j.id, j.batchId)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	return nil
}

func (db *sqlite) getSchedule(id string) (*schedule, *multi.Error) {
	return nil, nil
}

func (db *sqlite) saveSchedule(s schedule) (id string, err *multi.Error) {
	return "", nil
}

func (db *sqlite) close() *multi.Error {
	log := multi.NewLogger()
	err := db.sqlite3.Close()
	if err != nil {
		log.LogInternalError(err.Error())
	}
	return &multi.Error{Code: multi.InternalError, Message: "Error while closing db connection"}
}

func newDBConn() (database, *multi.Error) {
	db, err := sql.Open("sqlite3", "../db/schedular.db")
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
