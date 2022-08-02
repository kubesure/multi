package dispatcher

import (
	"github.com/kubesure/multi"
	"github.com/kubesure/multi/internal"
)

func SaveJob(j *internal.Job) *multi.Error {
	db, err := internal.NewDBConn(internal.DISPATCHER, internal.SQLITE)

	defer db.Close()
	if err != nil {
		return err
	}

	serr := db.SaveJob(j)
	if serr != nil {
		return serr
	}
	return nil
}

func UpdateJob(j internal.Job) *multi.Error {
	db, err := internal.NewDBConn(internal.DISPATCHER, internal.SQLITE)

	defer db.Close()
	if err != nil {
		return err
	}

	serr := db.UpdateJob(&j)
	if serr != nil {
		return serr
	}
	return nil
}
