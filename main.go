package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// Loading todos from JSON file
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

// Saving todos to JSON file
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

// Get all data
func getTodos(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Add new item
func addTodo(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var t Todo
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		// เพิ่มการแสดง error ใน log และไม่ให้แสดง "Invalid input" หากไม่มีข้อผิดพลาด
		log.Println("Error decoding input:", err) // แสดง error ที่เกิดขึ้นใน log
		http.Error(w, "Invalid input", 400)
		return
	}

	// Auto ID
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

// Update data item (toggle done)
func updateTodo(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	var updateData Todo
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid input", 400)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos[i].Task = updateData.Task
			todos[i].Done = updateData.Done
			saveTodos()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", 404)
}

// Delete todo item
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Todo not found", 404)
}

func main() {
	loadTodos()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodos(w, r)
		case http.MethodPost:
			addTodo(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			updateTodo(w, r)
		case http.MethodDelete:
			deleteTodo(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	log.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
