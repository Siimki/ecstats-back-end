# 🛠️ Estonian Cycling Stats – Backend (Go)

This is the backend for the Estonian Cycling Stats project. It serves data via REST API. The backend is written in Go and uses PostgreSQL as the database.

---

## 📁 Project Structure

```
/config/           → Configuration loader (from config.yaml)
/db/               → Database logic and queries
/handlers/         → HTTP handlers
/models/           → Structs and DB models
/utils/            → Helpers and async utilities
main.go            → Entry point
```

---

## Features

- API endpoints
- Rankings, stats, and homepage routes
- YAML-based config
- Modular structure (db, models, handlers)

---

## ⚙️ Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL running locally (port 5433)

### Config

Edit `config/config.yaml` for DB connection and other values.

### Run the server

```bash
go run main.go
```

The API will be available at `http://localhost:1337`

---

## 🧩 Connected to Frontend

The front-end (Next.js) fetches from this backend at `localhost:1337`. Make sure both are running for full functionality.

