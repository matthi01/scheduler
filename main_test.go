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

	req, _ := http.NewRequest("GET", "/items", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentItem(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/item/13", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Item not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Item not found'. Got %s", m["error"])
	}
}

func TestCreateItem(t *testing.T) {
	clearTable()

	jsonString := []byte(`{"name":"Test item","description": "Test description"}`)

	req, _ := http.NewRequest("POST", "/item", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Test item" {
		t.Errorf("Expected name of item to be 'Test item'. Got %s", m["name"])
	}
	if m["description"] != "Test description" {
		t.Errorf("Expected name of item to be 'Test description'. Got %s", m["description"])
	}
}

func TestFetchItem(t *testing.T) {
	clearTable()
	addItems(1)

	req, _ := http.NewRequest("GET", "/item/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateItem(t *testing.T) {
	clearTable()
	addItems(1)

	req, _ := http.NewRequest("GET", "/item/1", nil)
	response := executeRequest(req)

	var originalItem map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalItem)

	jsonString := []byte(`{"name": "Updated name", "description": "Updated description"}`)
	req, _ = http.NewRequest("PUT", "/item/1", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedItem map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &updatedItem)

	if updatedItem["id"] != originalItem["id"] {
		t.Errorf("Expected ID to be %v. Got %v", originalItem["id"], updatedItem["id"])
	}
	if updatedItem["name"] == originalItem["name"] {
		t.Errorf(
			"Expected name to change from %v to %v. Got %v",
			originalItem["name"],
			"Updated name",
			updatedItem["name"],
		)
	}
	if updatedItem["description"] == originalItem["description"] {
		t.Errorf(
			"Expected description to change from %v to %v. Got %v",
			originalItem["description"],
			"Updated description",
			updatedItem["description"],
		)
	}
}

func TestDeleteItem(t *testing.T) {
	clearTable()
	addItems(1)

	req, _ := http.NewRequest("DELETE", "/item/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/item/1", nil)
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
	a.DB.Exec("DELETE FROM items")
	a.DB.Exec("ALTER SEQUENCE items_id_seq RESTART WITH 1")
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

func addItems(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO items(name, description) VALUES($1, $2)", "Item "+strconv.Itoa(i), "Description "+strconv.Itoa(i))
	}
}

// constants
const tableCreationQuery = `CREATE TABLE IF NOT EXISTS items
(
    id SERIAL,
    name TEXT NOT NULL,
	description TEXT NOT NULL,
    CONSTRAINT items_pkey PRIMARY KEY (id)
)`
