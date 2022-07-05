package internal

import (
	"bytes"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
)

const tmpl = `{
    "query": {
        "bool": {
            "must": [
                {
                    "wildcard": {
                        "firstName": {{.FirstName}}
                    }
                }
            ]
        }
    }
}`

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func Schedule(customers *[]Customer) ScheduleResult {
	id := save(schedule{})
	return ScheduleResult{searchId: id}
}

func save(sch schedule) string {
	//db := newDBConn()
	//return db.saveSchedule(schedule{})
	return ""
}

func makeQuery(c Customer) {
	t, err := template.New("search-request").Parse(tmpl)
	//t, err := template.New("search-request").ParseFiles("search-request.txt")
	if err != nil {
		log.Errorf("error %v", err)
	}

	buff := new(bytes.Buffer)

	err1 := t.Execute(buff, c)
	if err1 != nil {
		log.Errorf("error %v", err)
	}
	log.Println(buff.String())
}
