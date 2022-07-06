package internal

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kubesure/multi"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type sqllite struct {
	sqlite3 *sql.DB
}

func (db *sqllite) saveBatch(b batch, jobs []job) (id string, err *multi.Error) {
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

	insertSQL := "INSERT INTO job(id,batch_id,payload,status,max_response,retry_interval,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?)"
	jobStmt, joberr := tx.Prepare(insertSQL)
	if joberr != nil {
		return "", &multi.Error{Code: multi.InternalError, Message: multi.DBError}
	}
	for index, job := range jobs {
		_, joberr = jobStmt.Exec(index+1, batchId, job.payload, "created", job.maxResponse, job.retryInterval, currentDateTime(), currentDateTime())
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

func (db *sqllite) getBatch(id string) (*batch, *multi.Error) {
	log := multi.NewLogger()
	row, qerr := db.sqlite3.Query("select id,type, created_datetime, updated_datetime from batch where id='" + id + "'")

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
		b.createdDate = *parseDateTime(created_datetime)
		b.updatedDate = *parseDateTime(updated_datetime)
	}

	return b, nil
}

func (db *sqllite) saveJob(j job) (id string, err *multi.Error) {
	return "", nil
}

func (db *sqllite) getJob(id string) (*job, *multi.Error) {
	return nil, nil
}

func (db *sqllite) getJobs(batchID string) ([]job, *multi.Error) {
	return nil, nil
}

func (db *sqllite) getSchedule(id string) (*schedule, *multi.Error) {
	return nil, nil
}

func (db *sqllite) saveSchedule(s schedule) (id string, err *multi.Error) {
	return "", nil
}

func (db *sqllite) close() *multi.Error {
	return nil
}

func newDBConn() (database, *multi.Error) {
	db, err := sql.Open("sqlite3", "../db/schedular.db")
	if err != nil {
		log.Println(err)
		return nil, &multi.Error{Code: multi.InternalError, Message: "DB Transaction Error"}
	}
	return &sqllite{sqlite3: db}, nil
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
