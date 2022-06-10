## Simple Todo Rest API with Golang

This project is a simple CRUD API using gorilla/mux.

It provides an API for manager TODO Tasks. 

These are the endpoints available:

| Operation      | Endpoint    | HTTP Verb |
|----------------|-------------|-----------|
| List all Todos | /todos      | GET       |
| Create Todo    | /todos      | POST      |
| Get Todo       | /todos/{id} | GET       |
| Update Todo    | /todos/{id} | PUT       |
| Delete Todo    | /todos/{id} | DELETE    |

## Running

```go
$ go mod tidy

$ go run main.go
```