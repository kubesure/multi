package internal

import "github.com/kubesure/multi"

func SaveJob(j job) *multi.Error {
	db, err := newDBConn(DISPATCHER, SQLITE)

	defer db.close()
	if err != nil {
		return err
	}
	return nil
}
