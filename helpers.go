package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

const categoriesTableCreation = `CREATE TABLE IF NOT EXISTS categories
(
    category_id SERIAL,
    name TEXT NOT NULL,
	description TEXT NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (category_id)
)`

const tasksTableCreation = `CREATE TABLE IF NOT EXISTS tasks
(
    task_id SERIAL PRIMARY KEY,
    category_id SERIAL references categories,
    task TEXT NOT NULL,
	seq SERIAL NOT NULL,
	complete BOOL NOT NULL
)`

var a App

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ensureTablesExists() {
	if _, err := a.DB.Exec(categoriesTableCreation); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(tasksTableCreation); err != nil {
		log.Fatal(err)
	}
}

func clearCategoriesTable() {
	a.DB.Exec("DELETE FROM categories")
	a.DB.Exec("ALTER SEQUENCE categories_category_id_seq RESTART WITH 1")
}

func clearTasksTable() {
	a.DB.Exec("DELETE FROM tasks")
	a.DB.Exec("ALTER SEQUENCE tasks_task_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE tasks_seq_seq RESTART WITH 1")
}

func clearTables() {
	clearTasksTable()
	clearCategoriesTable()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

func addCategory() int {
	var c category
	a.DB.QueryRow("INSERT INTO categories(name, description) VALUES($1, $2) RETURNING category_id", "Test Category", "Test Category Description").Scan(&c.Category_ID)
	return c.Category_ID
}

func addCategories(count int) []int {
	if count < 1 {
		count = 1
	}
	var categoryIds []int
	for i := 0; i < count; i++ {
		var c category
		name := "Category " + strconv.Itoa(i)
		desc := "Description " + strconv.Itoa(i)
		a.DB.QueryRow("INSERT INTO categories(name, description) VALUES($1, $2) RETURNING category_id", name, desc).Scan(&c.Category_ID)
		categoryIds = append(categoryIds, c.Category_ID)
	}
	return categoryIds
}

func addTaskToCategory(categoryId int) int {
	var t task
	a.DB.QueryRow("INSERT INTO tasks(category_id, task, complete) VALUES($1, $2, $3) RETURNING task_id", categoryId, "Test Task", false).Scan(&t.Task_ID)
	return t.Task_ID
}

func addTasksToCategory(categoryId, count int) []int {
	if count < 1 {
		count = 1
	}
	var taskIds []int
	for i := 0; i < count; i++ {
		var t task
		task := "Task " + strconv.Itoa(i)
		a.DB.QueryRow("INSERT INTO tasks(category_id, task, complete) VALUES($1, $2, $3) RETURNING task_id", categoryId, task, false).Scan(&t.Task_ID)
		taskIds = append(taskIds, t.Task_ID)
	}
	return taskIds
}
