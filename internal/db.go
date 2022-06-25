package internal

type database interface {
	getSchedule(id string) *schedule
	saveSchedule(s schedule) (id string)
}

type sqllite struct{}

func (db *sqllite) getSchedule(id string) *schedule {
	return nil
}

func (db *sqllite) saveSchedule(s schedule) (id string) {
	return ""
}
