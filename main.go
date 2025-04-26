package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var (
	todos []Todo
	mu    sync.Mutex
)

const filePath = "todos.json"

func loadTodos() {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		todos = []Todo{}
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &todos)
}

func saveTodos() {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		log.Println("Error marshalling todos:", err)
		return
	}
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Println("Error writing file:", err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input", 400)
		return
	}

	// กำหนด id แบบอัตโนมัติ
	if len(todos) == 0 {
		t.ID = 1
	} else {
		t.ID = todos[len(todos)-1].ID + 1
	}
	todos = append(todos, t)

	saveTodos()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func main() {
	loadTodos()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodos(w, r)
		} else if r.Method == http.MethodPost {
			addTodo(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	log.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
