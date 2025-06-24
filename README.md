# TODO_MUX
A sample app for testing keploy modifications

# Run the app
```bash
go run main.go
go run .
go build . && ./todo-ts
```

# Create
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Todo","completed":false}'

# List all
curl http://localhost:8080/todos

# Get single
curl http://localhost:8080/todos/1

# Update
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","completed":true}'

# Delete
curl -X DELETE http://localhost:8080/todos/1