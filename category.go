package main

import (
	"database/sql"
)

type category struct {
	Category_ID int    `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *category) createCategory(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO categories(name, description) VALUES ($1, $2) RETURNING category_id",
		c.Name, c.Description,
	).Scan(&c.Category_ID)
	return err
}

func (c *category) getCategory(db *sql.DB) error {
	return db.QueryRow(
		"SELECT name, description FROM categories WHERE category_id=$1",
		c.Category_ID,
	).Scan(&c.Name, &c.Description)
}

func (c *category) updateCategory(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE categories SET name=$1, description=$2 WHERE category_id=$3",
		c.Name, c.Description, c.Category_ID,
	)
	return err
}

func (c *category) deleteCategory(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM categories WHERE category_id=$1", c.Category_ID)
	return err
}

func (c *category) deleteCategoryTasks(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM tasks WHERE category_id=$1", c.Category_ID)
	return err
}

func getCategories(db *sql.DB) ([]category, error) {
	rows, err := db.Query("SELECT category_id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []category{}
	for rows.Next() {
		var c category
		err := rows.Scan(&c.Category_ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
