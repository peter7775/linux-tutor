package repository

import "database/sql"

type ProgressRepo struct{ DB *sql.DB }

func (r ProgressRepo) Load() (int, int, error)                                          { return 0, 0, nil }
func (r ProgressRepo) Save(c, w int) error                                              { return nil }
func (r ProgressRepo) SaveAttempt(topic, prompt, answer, notes string, delta int) error { return nil }
