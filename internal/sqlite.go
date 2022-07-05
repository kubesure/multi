package internal

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqllite struct {
	sqlite3 *sql.DB
}

func (db *sqllite) getSchedule(id string) *schedule {
	return nil
}

func (db *sqllite) saveSchedule(s schedule) (id string) {
	return ""
}

func (db *sqllite) saveBatch(b batch, jobs []job) (id string, err error) {
	tx, txerr := db.sqlite3.Begin()
	if txerr != nil {
		return "", txerr
	}

	stmt, serr := tx.Prepare("INSERT INTO batch(id,type,created_datetime,updated_datetime) VALUES (?,?,?,?)")
	if serr != nil {
		return "", serr
	}

	batchId := uuid.New().String()

	_, err1 := stmt.Exec(batchId, b.ttype, currentDateTime(), currentDateTime())
	defer stmt.Close()
	if err1 != nil {
		return "", err
	}

	insertSQL := "INSERT INTO job(id,batch_id,payload,status,max_response,retry_interval,created_datetime,updated_datetime)" +
		"VALUES (?,?,?,?,?,?,?,?)"
	jobStmt, joberr := tx.Prepare(insertSQL)
	if joberr != nil {
		return "", joberr
	}
	for index, job := range jobs {
		_, joberr = jobStmt.Exec(index+1, batchId, job.payload, "created", job.maxResponse, job.retryInterval, currentDateTime(), currentDateTime())
		if joberr != nil {
			tx.Rollback()
			return "", joberr
		}
	}

	txerr = tx.Commit()
	if txerr != nil {
		return "", txerr
	}

	return batchId, nil
}

func (db *sqllite) getBatch(id string) *batch {
	return nil
}

/*func (db *sqllite) saveJob(j job) (id string) {
	return ""
}*/

func (db *sqllite) getJob(id string) *job {
	return nil
}

func (db *sqllite) getJobs(batchID string) []job {
	return nil
}

func (db *sqllite) close() error {
	return nil
}

func newDBConn() (database, error) {
	db, err := sql.Open("sqlite3", "../db/schedular.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &sqllite{sqlite3: db}, nil
}

func currentDateTime() string {
	return time.Now().Format(time.RFC3339)
}
