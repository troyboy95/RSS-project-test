# RSS‑project‑test

A simple self‑hosted RSS feed aggregator written in Go.  
Users can manage accounts, subscribe to feeds, and fetch posts from followed feeds.

---

## 🧱 Overview

- **Language**: Go  
- **Database**: SQL (PostgreSQL  via `sqlc`)  
- **Features**:
  - Add, list & follow RSS feeds  
  - Fetch the latest posts for followed feeds  
  - User authentication middleware  
  - Health/readiness endpoints  
  - Automated Background scraping via concurrency

---

## ⚙️ Prerequisites

- Go **1.20+**
- PostgreSQL or MySQL
- `sqlc` installed for code generation (see [sqlc.dev](https://sqlc.dev))
- Database migration tool `goose`

---

## 🚀 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/troyboy95/RSS-project-test.git
   cd RSS-project-test

2. **Install dependencies & generate SQL bindings**
    go mod tidy
    sqlc generate

3. **Apply database migrations**
    goose up -dir sql/schema <"Your-connection-string">

4. **Build & run the server**
    go build; ./RSS-project-test --for windows
    go build & ./RSS-project-test --for mac/linux
