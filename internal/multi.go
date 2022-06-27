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

func SaveBatch(batchType BatchType, c CustomerSearch) (id string) {
	return ""
}

func newDBConn() database {
	return &sqllite{}
}
