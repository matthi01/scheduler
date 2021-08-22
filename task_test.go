package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestEmptyTasksTable(t *testing.T) {
	clearTables()
	categoryId := addCategory()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v/tasks", categoryId), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTask(t *testing.T) {
	clearTables()
	categoryId := addCategory()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v/task/999", categoryId), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Task not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Task not found'. Got %s", m["error"])
	}
}

func TestCreateTask(t *testing.T) {
	clearTables()
	categoryId := addCategory()

	jsonString := []byte(`{"task":"Test task","complete": false}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/category/%v/task", categoryId), bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["task"] != "Test task" {
		t.Errorf("Expected name of category to be 'Test task'. Got %s", m["task"])
	}
}

func TestFetchTask(t *testing.T) {
	clearTables()
	categoryId := addCategory()

	taskCount := 2
	addTasksToCategory(categoryId, taskCount)

	for i := 0; i < taskCount; i++ {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v/task/%v", categoryId, i+1), nil)
		response := executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var m map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &m)

		if m["task"] != fmt.Sprintf("Task %v", i) {
			t.Errorf("Expected name of task to be 'Task %v'. Got %s", i+1, m["task"])
		}
	}
}

func TestFetchAllTasks(t *testing.T) {
	clearTables()
	categoryId := addCategory()

	addTasksToCategory(categoryId, 3)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v/tasks", categoryId), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var tasks []task
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&tasks); err != nil {
		t.Errorf("Could not read tasks data from request response. Error: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("Expected to receive 3 tasks under category_id %v. Got %v", categoryId, len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	clearTables()
	categoryId := addCategory()
	taskId := addTaskToCategory(categoryId)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/category/%v/task/%v", categoryId, taskId), nil)
	response := executeRequest(req)

	var originalTask map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTask)

	jsonString := []byte(`{"task":"updated task", "complete":true}`)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/category/%v/task/%v", categoryId, taskId), bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedTask map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &updatedTask)

	if updatedTask["task_id"] != originalTask["task_id"] {
		t.Errorf("Expected Task ID to be %v. Got %v", originalTask["task_id"], updatedTask["task_id"])
	}
	if updatedTask["task"] == originalTask["task"] {
		t.Errorf("Expected task to change from '%s' to '%s'. Got '%s'.", originalTask["task"], "updated task", updatedTask["task"])
	}
	if updatedTask["complete"] == originalTask["complete"] {
		t.Errorf("Expected complete to change from '%s' to '%s'. Got '%s'.", originalTask["complete"], "true", updatedTask["complete"])
	}
}

func TestDeleteTask(t *testing.T) {
	clearTables()
	categoryId := addCategory()
	taskId := addTaskToCategory(categoryId)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/category/%v/task/%v", categoryId, taskId), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/category/%v/task/%v", categoryId, taskId), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
