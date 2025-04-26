# Todo List API (Golang + Docker)
A simple RESTful API for managing a todo list, built with Go and stores data in a JSON file.

## Features
- View all todos
- Add a new todo
- Save todos to a local JSON file
- Run easily with Docker

## Requirements
- Go 1.24+
- Docker

## Setup Instructions

1. Clone the project
git clone [https://github.com/boraprae/go-lang-practice.git](https://github.com/boraprae/go-lang-practice.git)
cd todo-list-practice

2. (Optional) Generate go.sum
If not already present, run:

--bash--
`go mod tidy`

This will generate the go.sum file.

## Running with Docker
1. Build the Docker image
--bash--
`docker build -t todo-app .`

2. Create an empty JSON file (if not existing)
--bash--
`echo [] > todos.json`

3. Run the container
Bind mount the todos.json into the container to persist data:

--bash--
`docker run -p 8080:8080 -v ${PWD}/todos.json:/app/todos.json todo-app`

(If you are using PowerShell on Windows, replace ${PWD} with $PWD.)

## API Endpoints
### Get all todos
--http--
`GET /todos`

Response:

-json--
[
  {
    "id": 1,
    "task": "Learn Go",
    "done": false
  }
]

### Add a new todo
--http--
`POST /todos`
`Content-Type: application/json`

Body:
{
  "task": "Write Dockerfile",
  "done": false
}
Response:
json
{
  "id": 2,
  "task": "Write Dockerfile",
  "done": false
}

## Notes
All todos are stored inside todos.json.

When the server starts, it loads existing todos from the file.

When you add a new todo, it updates the todos.json automatically.

This project is for learning purposes. No database is required.
