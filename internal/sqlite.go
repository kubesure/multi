package internal

type sqllite struct{}

func (db *sqllite) getSchedule(id string) *schedule {
	return nil
}

func (db *sqllite) saveSchedule(s schedule) (id string) {
	return ""
}

func (db *sqllite) saveBatch(b batch) (id string) {
	return id
}

func (db *sqllite) getBatch(id string) *batch {
	return nil
}

func (db *sqllite) saveJob(j job) (id string) {
	return ""
}
func (db *sqllite) getJob(id string) *job {
	return nil
}

func (db *sqllite) getJobs(batchID string) []job {
	return nil
}
