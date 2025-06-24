package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	todos  = make(map[int]Todo)
	nextID = 1
	lock   sync.Mutex
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	r.Use(loggerMiddleware)

	r.HandleFunc("/todos", listTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func listTodos(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var todoList []Todo
	for _, todo := range todos {
		todoList = append(todoList, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newTodo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	newTodo.ID = nextID
	todos[nextID] = newTodo
	nextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	todo, exists := todos[id]
	if !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	existingTodo, exists := todos[id]
	if !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if updatedTodo.Title != "" {
		existingTodo.Title = updatedTodo.Title
	}
	existingTodo.Completed = updatedTodo.Completed

	todos[id] = existingTodo

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	_, exists := todos[id]
	if !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	delete(todos, id)
	w.WriteHeader(http.StatusNoContent)
}
