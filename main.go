package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id        string    `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Owner     string    `json:"owner"`
}

var todos []Todo
var wg sync.WaitGroup
var id int32 = 2

func nextId() string {
	defer wg.Done()
	atomic.AddInt32(&id, 1)
	return strconv.Itoa(int(id))
}

func getTodos(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(todos)
}

func getTodo(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for _, item := range todos {
		if item.Id == params["id"] {
			json.NewEncoder(rw).Encode(item)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func deleteTodo(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range todos {
		if item.Id == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func createTodo(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var todo Todo

	err := json.NewDecoder(req.Body).Decode(&todo)

	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	todo.CreatedAt = time.Now()
	wg.Add(1)
	go func() { todo.Id = nextId() }()
	wg.Wait()
	todos = append(todos, todo)
	json.NewEncoder(rw).Encode(todo)
}

func updateTodo(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range todos {
		if item.Id == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			var todo Todo
			err := json.NewDecoder(req.Body).Decode(&todo)
			if err != nil {
				fmt.Println(err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			todo.Id = params["id"]
			if todo.Text == "" {
				todo.Text = item.Text
			}
			if todo.Owner == "" {
				todo.Owner = item.Owner
			}
			todo.CreatedAt = item.CreatedAt
			todo.UpdatedAt = time.Now()
			todos = append(todos, todo)
			json.NewEncoder(rw).Encode(todo)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func main() {
	todos = append(todos,
		Todo{
			Id:        "1",
			Text:      "Learn Go",
			Done:      false,
			CreatedAt: time.Now(),
			Owner:     "Robson",
		},
		Todo{
			Id:        "2",
			Text:      "Learn Python",
			Done:      false,
			CreatedAt: time.Now(),
			Owner:     "Robson",
		},
	)
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")

	http.Handle("/", r)

	http.ListenAndServe(":8000", nil)
}
