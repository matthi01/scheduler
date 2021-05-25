package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	loadEnv()
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/categories", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentCategory(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/category/13", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Category not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Category not found'. Got %s", m["error"])
	}
}

func TestCreateCategory(t *testing.T) {
	clearTable()

	jsonString := []byte(`{"name":"Test category","description": "Test description"}`)

	req, _ := http.NewRequest("POST", "/category", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Test category" {
		t.Errorf("Expected name of category to be 'Test category'. Got %s", m["name"])
	}
	if m["description"] != "Test description" {
		t.Errorf("Expected name of category to be 'Test description'. Got %s", m["description"])
	}
}

func TestFetchCategory(t *testing.T) {
	clearTable()
	addCategories(1)

	req, _ := http.NewRequest("GET", "/category/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateCategory(t *testing.T) {
	clearTable()
	addCategories(1)

	req, _ := http.NewRequest("GET", "/category/1", nil)
	response := executeRequest(req)

	var originalCategory map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCategory)

	jsonString := []byte(`{"name": "Updated name", "description": "Updated description"}`)
	req, _ = http.NewRequest("PUT", "/category/1", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedCategory map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &updatedCategory)

	if updatedCategory["category_id"] != originalCategory["category_id"] {
		t.Errorf("Expected ID to be %v. Got %v", originalCategory["category_id"], updatedCategory["category_id"])
	}
	if updatedCategory["name"] == originalCategory["name"] {
		t.Errorf(
			"Expected name to change from %v to %v. Got %v",
			originalCategory["name"],
			"Updated name",
			updatedCategory["name"],
		)
	}
	if updatedCategory["description"] == originalCategory["description"] {
		t.Errorf(
			"Expected description to change from %v to %v. Got %v",
			originalCategory["description"],
			"Updated description",
			updatedCategory["description"],
		)
	}
}

func TestDeleteCategory(t *testing.T) {
	clearTable()
	addCategories(1)

	req, _ := http.NewRequest("DELETE", "/category/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/category/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// helpers
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM categories")
	a.DB.Exec("ALTER SEQUENCE categories_category_id_seq RESTART WITH 1")
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

func addCategories(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO categories(name, description) VALUES($1, $2)", "Category "+strconv.Itoa(i), "Description "+strconv.Itoa(i))
	}
}

// constants
const tableCreationQuery = `CREATE TABLE IF NOT EXISTS categories
(
    category_id SERIAL,
    name TEXT NOT NULL,
	description TEXT NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (category_id)
)`
