# 📝 Go To-Do List App (API + CLI)

This project is a beginner-friendly **To-Do List** app written in Go. It allows users to:

✅ Add tasks  
📋 View tasks  
✏️ Update task description or status  
❌ Delete tasks  

Tasks are stored in a local `tasks.json` file. The app supports both:

- A **Command-Line Interface (CLI)**
- A **REST API (HTTP server)**

---

## 📁 Task Format

Each task includes:

- `description` (string)
- `status` (one of):
  - `"Not started"`
  - `"Started"`
  - `"Completed"`

---

## ▶️ Run the App

Start the HTTP server on `http://localhost:8080`:

```bash
go run main.go

```


# 🌐 Using the API (All Commands with `curl`)

Below are all available API endpoints with example commands you can run from **Git Bash**, **WSL**, or **macOS/Linux terminal**.

---

### ✅ CREATE a New Task  
**POST /create**  
Adds a new task with a description and status:

```bash
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -d '{"description": "Buy oat milk", "status": "Not started"}'
```

---

### 📋 GET All Tasks  
**GET /get**  
Returns all tasks stored in `tasks.json`:

```bash
curl http://localhost:8080/get
```

curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"index": 0, "status": "In Progress"}'

  curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"index": 0, "status": "Completed"}'

Example response:

```json
[
  { "description": "Buy oat milk", "status": "Not started" },
  { "description": "Walk dog", "status": "Completed" }
]
```

---

### ✏️ UPDATE a Task  
**POST /update**  
You can update **description**, **status**, or **both**. Task index starts at **0**.

Update description only:
```bash
curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"index": 0, "description": "Walk the dog"}'
```

Update status only:
```bash
curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"index": 1, "status": "Completed"}'
```

Update both description and status:
```bash
curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"index": 2, "description": "Wash car", "status": "Started"}'
```

---

### ❌ DELETE a Task  
**POST /delete**  
Deletes task at the provided index:

```bash
curl -X POST http://localhost:8080/delete \
  -H "Content-Type: application/json" \
  -d '{"index": 1}'
```

Expected result: HTTP 204 No Content (task deleted)

---

### 📌 Notes

- Task indexes start from **0**
- Valid statuses are:
  - `"Not started"`
  - `"Started"`
  - `"Completed"`
- All data is saved to `tasks.json` in the project directory
- Each endpoint automatically loads, modifies, and saves the task list

---

### 🧪 Running Unit Tests

Use this command to run all your Go tests in the `store` package:

```bash
go test ./store
```
