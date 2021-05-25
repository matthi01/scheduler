package main

import (
	"database/sql"
)

type item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (i *item) createItem(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO items(name, description) VALUES ($1, $2) RETURNING id",
		i.Name, i.Description,
	).Scan(&i.ID)
	return err
}

func (i *item) getItem(db *sql.DB) error {
	return db.QueryRow(
		"SELECT name, description FROM items WHERE id=$1",
		i.ID,
	).Scan(&i.Name, &i.Description)
}

func (i *item) updateItem(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE items SET name=$1, description=$2 WHERE id=$3",
		i.Name, i.Description, i.ID,
	)
	return err
}

func (i *item) deleteItem(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM items WHERE id=$1", i.ID)
	return err
}

func getItems(db *sql.DB) ([]item, error) {
	rows, err := db.Query("SELECT name, description FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []item{}
	for rows.Next() {
		var i item
		err := rows.Scan(&i.ID, &i.Name, &i.Description)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
