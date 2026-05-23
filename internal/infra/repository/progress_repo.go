package repository

import "database/sql"

type ProgressRepo struct { DB *sql.DB }

func (r ProgressRepo) Load() (correct, wrong int, err error) {
	row := r.DB.QueryRow(`SELECT correct, wrong FROM progress WHERE id = 1`)
	err = row.Scan(&correct, &wrong)
	return
}

func (r ProgressRepo) Save(correct, wrong int) error {
	_, err := r.DB.Exec(`UPDATE progress SET correct = ?, wrong = ? WHERE id = 1`, correct, wrong)
	return err
}
