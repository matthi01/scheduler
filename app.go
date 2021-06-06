package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, dbname string) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	// categories
	a.Router.HandleFunc("/categories", a.getCategories).Methods("GET")
	a.Router.HandleFunc("/category", a.createCategory).Methods("POST")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}", a.getCategory).Methods("GET")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}", a.updateCategory).Methods("PUT")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}", a.deleteCategory).Methods("DELETE")

	// tasks
	a.Router.HandleFunc("/category/{category_id:[0-9]+}/tasks", a.getTasks).Methods("GET")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}/task", a.createTask).Methods("POST")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}/task/{task_id:[0-9]+}", a.getTask).Methods("GET")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}/task/{task_id:[0-9]+}", a.updateTask).Methods("PUT")
	a.Router.HandleFunc("/category/{category_id:[0-9]+}/task/{task_id:[0-9]+}", a.deleteTask).Methods("DELETE")
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

// categories
func (a *App) getCategory(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	c := category{Category_ID: id}
	if err := c.getCategory(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Category not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) getCategories(w http.ResponseWriter, req *http.Request) {
	categories, err := getCategories(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, categories)
}

func (a *App) updateCategory(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var c category
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()
	c.Category_ID = id

	if err := c.updateCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) createCategory(w http.ResponseWriter, req *http.Request) {
	var c category
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	if err := c.createCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

// to do: when deleting a category you need to also delete all tasks under that category first
func (a *App) deleteCategory(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	c := category{Category_ID: id}
	if err = c.deleteCategoryTasks(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err = c.deleteCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// tasks
func (a *App) createTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	t := task{Category_ID: id}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	if err := t.createTask(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, t)
}

func (a *App) getTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	categoryId, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	taskId, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	t := task{Category_ID: categoryId, Task_ID: taskId}
	if err := t.getTask(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Task not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) getTasks(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	categoryId, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	c := category{Category_ID: categoryId}
	tasks, err := c.getTasks(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (a *App) updateTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskId, err := strconv.Atoi(vars["task_id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var t task
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer req.Body.Close()

	t.Task_ID = taskId
	if err := t.updateTask(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) deleteTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskId, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	t := task{Task_ID: taskId}
	if err := t.deleteTask(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
