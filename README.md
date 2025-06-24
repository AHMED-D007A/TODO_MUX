# TODO_MUX
A sample app for testing keploy modifications

# Run the app
```bash
go run main.go
go run .
go build . && ./todo-ts
```

# Create
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Todo","completed":false}'
```

# List all
```bash
curl http://localhost:8080/todos
```

# Get single
```bash
curl http://localhost:8080/todos/1
```

# Update
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","completed":true}'

# Delete
```bash
curl -X DELETE http://localhost:8080/todos/1
```