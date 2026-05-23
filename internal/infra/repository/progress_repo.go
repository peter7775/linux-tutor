package repository

import "database/sql"

type ProgressRepo struct { DB *sql.DB }
func (r ProgressRepo) Load() (int, int, error) { var c, w int; err := r.DB.QueryRow(`SELECT correct, wrong FROM progress WHERE id=1`).Scan(&c, &w); return c, w, err }
func (r ProgressRepo) Save(c, w int) error { _, err := r.DB.Exec(`UPDATE progress SET correct=?, wrong=? WHERE id=1`, c, w); return err }
