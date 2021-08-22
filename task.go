package main

import (
	"database/sql"
)

type task struct {
	Task_ID     int    `json:"task_id"`
	Category_ID string `json:"category_id"`
	Task        string `json:"task"`
	Seq         int    `json:"seq"`
	Complete    bool   `json:"complete"`
}

func (t *task) createTask(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO tasks(category_id, task, complete) VALUES ($1, $2, $3) RETURNING task_id",
		t.Category_ID, t.Task, t.Complete,
	).Scan(&t.Task_ID)
	return err
}

func (t *task) getTask(db *sql.DB) error {
	return db.QueryRow(
		"SELECT task_id, category_id, task, seq, complete FROM tasks WHERE task_id=$1",
		t.Task_ID,
	).Scan(&t.Task_ID, &t.Category_ID, &t.Task, &t.Seq, &t.Complete)
}

func (t *task) updateTask(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE tasks SET task=$1, seq=$2, complete=$3 WHERE task_id=$4",
		t.Task, t.Seq, t.Complete, t.Task_ID,
	)
	return err
}

func (t *task) deleteTask(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM tasks WHERE task_id=$1", t.Task_ID)
	return err
}

func (c *category) getTasks(db *sql.DB) ([]task, error) {
	rows, err := db.Query("SELECT task_id, category_id, task, seq, complete FROM tasks WHERE category_id=$1", c.Category_ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []task{}
	for rows.Next() {
		var tsk task
		err := rows.Scan(&tsk.Task_ID, &tsk.Category_ID, &tsk.Task, &tsk.Seq, &tsk.Complete)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, tsk)
	}
	return tasks, nil
}
