package internal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kubesure/multi"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

//TODO Implement batch_id as binary type in all crud statements

type sqlite struct {
	sqlite3 *sql.DB
}

func (db *sqlite) SaveBatch(jobs []Job) (batch *Batch, err *multi.Error) {
	log := multi.NewLogger()
	var b *Batch = &Batch{}
	tx, txerr := db.sqlite3.Begin()
	if txerr != nil {
		log.LogInternalError(txerr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	stmt, serr := tx.Prepare("INSERT INTO batch(id,created_datetime,updated_datetime) VALUES (?,?,?)")
	if serr != nil {
		log.LogInternalError(serr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	b.id = uuid.New().String()

	_, err1 := stmt.Exec(b.id, currentDateTime(), currentDateTime())
	defer stmt.Close()
	if err1 != nil {
		log.LogInternalError(err1.Error())
		return nil, err
	}

	insertEPSQL := "INSERT INTO endpoint (job_id,uri,method,auth_type,auth_srvcert,auth_uname,auth_pass,headers)" +
		"VALUES (?,?,?,?,?,?,?,?)"
	ePStmt, ePErr := tx.Prepare(insertEPSQL)
	if ePErr != nil {
		log.LogInternalError(ePErr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	insertSQL := "INSERT INTO job(id,batch_id,payload,compress_dispatch,status,error_msg,max_response,retry_interval,retry_count,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	jobStmt, joberr := tx.Prepare(insertSQL)
	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	for _, job := range jobs {
		_, joberr = jobStmt.Exec(job.Id, b.id, job.Payload.Data, job.Payload.CompressedDispatch, CREATED, job.ErrorMsg, job.MaxResponseSeconds, job.RetryIntervalSeconds, job.MaxRetry, currentDateTime(), currentDateTime())
		headers := makeHeaders(job.EndPoint.Headers)
		_, ePErr = ePStmt.Exec(job.Id, job.EndPoint.Uri, job.EndPoint.Method, job.EndPoint.Auth.Type, job.EndPoint.Auth.ServerCertificate, job.EndPoint.Auth.UserName, job.EndPoint.Auth.Password, headers)
		if joberr != nil || ePErr != nil {
			tx.Rollback()
			if joberr != nil {
				log.LogInternalError(joberr.Error())
			}
			if ePErr != nil {
				log.LogInternalError(ePErr.Error())
			}
			return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
		}
	}

	txerr = tx.Commit()
	if txerr != nil {
		log.LogInternalError(txerr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	return b, nil
}

func makeHeaders(headers []Header) string {
	var h string
	for _, v := range headers {
		h += fmt.Sprintf("%v:%v-", v.Key, v.Value)
	}
	return h
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
	q := "select id,batch_id,payload,compress_dispatch,result,status,created_datetime,updated_datetime,max_response,retry_interval,error_msg,retry_count from job where batch_id=?"
	jobs, err := db.Getjobs(q, batchID)
	if err != nil {
		log.LogInternalError(err.Inner.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}

	return jobs, nil
}

func (db *sqlite) GetJob(jobID, batchID string) (*Job, *multi.Error) {
	q := "select id,batch_id,payload,compress_dispatch,result,status,created_datetime,updated_datetime,max_response,retry_interval,error_msg,retry_count from job where batch_id=? and id = ?"
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
		j.Payload = &Payload{}
		var created_datetime, updated_datetime string
		var errormsg sql.NullString
		var retrycount sql.NullInt32
		var result sql.NullString

		serr := rows.Scan(&j.Id, &j.BatchId, &j.Payload.Data, &j.Payload.CompressedDispatch, &result, &j.Status, &created_datetime,
			&updated_datetime, &j.MaxResponseSeconds, &j.RetryIntervalSeconds, &errormsg, &retrycount)

		if serr != nil {
			log.LogInternalError(serr.Error())
			return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError, Inner: serr}
		}

		if errormsg.Valid {
			j.ErrorMsg = &errormsg.String
		} else {
			val := ""
			j.ErrorMsg = &val
		}

		if retrycount.Valid {
			count := uint(retrycount.Int32)
			j.MaxRetry = count
		} else {
			var count uint = 0
			j.MaxRetry = count
		}

		if result.Valid {
			j.Result = &result.String
		} else {
			val := ""
			j.Result = &val
		}

		j.CreatedDateTime = *parseDateTime(created_datetime)
		j.UpdatedDateTime = *parseDateTime(updated_datetime)
		jobs = append(jobs, j)
	}

	rerr := rows.Err()
	if rerr != nil {
		log.LogInternalError(qerr.Error())
		return nil, &multi.Error{Code: multi.InternalError, Message: multi.DBError, Inner: rerr}
	}

	return jobs, nil
}

func (db *sqlite) SaveJob(j *Job) (err *multi.Error) {
	log := multi.NewLogger()
	insertSQL := "INSERT INTO job(id,batch_id,payload,compress_dispatch,status,error_msg,max_response,retry_interval,retry_count,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	jobStmt, joberr := db.sqlite3.Prepare(insertSQL)
	if joberr != nil {
		log.LogInternalError("error inserting job " + joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	_, joberr = jobStmt.Exec(j.Id, j.BatchId, j.Payload.Data, j.Payload.CompressedDispatch, CREATED, j.ErrorMsg, j.MaxResponseSeconds, j.RetryIntervalSeconds, j.MaxRetry, currentDateTime(), currentDateTime())

	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	return nil
}

func (db *sqlite) UpdateJob(j *Job) (err *multi.Error) {
	log := multi.NewLogger()
	updateSQL := "UPDATE job SET status = coalesce(?,status), error_msg = coalesce(?,error_msg),retry_count = coalesce(?,retry_count)" +
		",result = coalesce(?,result), updated_datetime = ? where id = ?"

	//	and batch_id = ?
	jobStmt, joberr := db.sqlite3.Prepare(updateSQL)

	if joberr != nil {
		log.LogInternalError(joberr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError, Inner: joberr}
	}

	result, serr := jobStmt.Exec(j.Status, j.ErrorMsg, j.MaxRetry, currentDateTime(), j.Result, j.Id)

	if err != nil {
		log.LogInternalError(serr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError, Inner: serr}
	}

	count, rserr := result.RowsAffected()
	if rserr != nil {
		log.LogInternalError(rserr.Error())
		return &multi.Error{Code: multi.InternalError, Message: multi.DBError, Inner: rserr}
	}

	if count == 0 {
		log.LogInternalError("no job updated")
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
