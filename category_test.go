package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestEmptyCategoriesTable(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/categories", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentCategory(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/category/b119178b-2fd2-4a5c-9301-190c341df180", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Category not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Category not found'. Got %s", m["error"])
	}
}

func TestCreateCategory(t *testing.T) {
	clearTables()

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
	clearTables()
	categoryId := addCategory()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v", categoryId), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestFetchAllCategories(t *testing.T) {
	clearTables()
	addCategories(3)

	req, _ := http.NewRequest("GET", "/categories", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var categories []category
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&categories); err != nil {
		t.Errorf("Could not read categories data from request response. Error: %v", err)
	}

	if len(categories) != 3 {
		t.Errorf("Expected to receive 3 categories. Got %v", len(categories))
	}
}

func TestUpdateCategory(t *testing.T) {
	clearTables()
	categoryId := addCategory()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v", categoryId), nil)
	response := executeRequest(req)

	var originalCategory map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCategory)

	jsonString := []byte(`{"name": "Updated name", "description": "Updated description"}`)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/category/%v", categoryId), bytes.NewBuffer(jsonString))
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
	clearTables()
	categoryId := addCategory()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/category/%v", categoryId), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/category/%v", categoryId), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// task related tests:
func TestDeleteCategoryWithTasks(t *testing.T) {
	clearTables()
	categoryId := addCategory()
	addTasksToCategory(categoryId, 3)

	// delete the category
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/category/%v", categoryId), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/category/%v", categoryId), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	// check that tasks were deleted
	req, _ = http.NewRequest("GET", fmt.Sprintf("/category/%v/tasks", categoryId), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var tasks []task
	decoder := json.NewDecoder(response.Body)

	if err := decoder.Decode(&tasks); err != nil {
		t.Errorf("Could not read tasks data from request response. Error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected to receive 0 tasks under category_id %v. Got %v", categoryId, len(tasks))
	}
}

// handle it on the front-end and
// func TestUpdateCategoryTaskOrdering(t *testing.T) {
// 	clearTables()
// 	categoryId := addCategory()
// 	addTasksToCategory(categoryId, 5)

// 	// get all tasks check their order
// 	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v/tasks", categoryId), nil)
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	var tasks []task
// 	decoder := json.NewDecoder(response.Body)
// 	if err := decoder.Decode(&tasks); err != nil {
// 		t.Errorf("Could not read tasks data from request response. Error: %v", err)
// 	}

// 	var originalTasksListOrder []string
// 	for _, t := range tasks {
// 		originalTasksListOrder = append(originalTasksListOrder, t.Task)
// 	}
// 	fmt.Println(originalTasksListOrder)

// 	fmt.Println(tasks)

// 	// reverse the sequence of the tasks
// 	for i := 0; i < len(tasks); i++ {
// 		j := len(tasks) - i
// 		tasks[i].Seq = j
// 	}

// 	fmt.Println(tasks)

// 	// update order
// 	jsonString, err := json.Marshal(tasks)
// 	if err != nil {
// 		t.Errorf("Could not encode tasks list back into json. Error: %v", err)
// 	}

// 	req, _ = http.NewRequest("PUT", "/category/1/tasks", bytes.NewBuffer(jsonString))
// 	executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	// get all tasks and check their order

// 	t.Error("Not Implemented.")
// }
