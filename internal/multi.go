package internal

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func SaveBatch(batchType BatchType, c CustomerSearch) (id string, err error) {
	db, err := newDBConn()
	defer db.close()
	if err != nil {
		return "", err
	}

	b := batch{ttype: batchType}
	jobs := []job{}
	id, err1 := db.saveBatch(b, jobs)
	if err1 != nil {
		return "", err1
	}
	return "", nil
}
