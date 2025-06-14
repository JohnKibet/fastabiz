
# 🚚 logistics-system

A scalable logistics & operations management system built to support Admins, Drivers, and Customers with clean APIs, secure role-based access, and smooth delivery workflows.

---

## 🚀 Features

- 🔐 Role-based access control (Admin, Driver, Customer)
- 📦 Order management: create, assign, and track orders
- 🚚 Delivery and Driver management logic
- 💳 Payment and feedback handling
- 🧑‍💼 Token-based authentication (mocked for now)
- 📄 Swagger/OpenAPI API documentation
- 🐳 Dockerized backend and PostgreSQL setup
- 🧪 Postman/Newman API testing
- 📊 Planned: dashboard analytics, API gateway (Kong), CI/CD integration

---

## 🛠️ Tech Stack

| Layer       | Technologies                                |
|-------------|---------------------------------------------|
| Backend     | Go (Chi router, Clean Architecture)         |
| Frontend    | ASP.NET Core Razor Pages (C#)               |
| Database    | PostgreSQL                                  |
| DevOps      | Docker, Docker Compose                      |
| Auth        | Token-based (mocked; JWT planned)           |
| Docs        | Swagger / OpenAPI                           |

---

## 🗂️ Project Structure

```
logistics-system/
├── apps/
│   ├── logistics-backend/  # Go API backend (modular by domain)
│   └── frontend/           # Razor Pages frontend (C#)
├── migrations/             # SQL migration scripts
├── postman/                # Postman collections & environments
├── docker-compose.yml      # Multi-service orchestration
└── README.md
```

Backend follows **Clean Architecture**:
- Modular usecase/repository layering
- Decoupled DTOs for input/output validation
- Unit-testable services and handlers

---

## ⚙️ Getting Started

### 🔧 Prerequisites

- [Docker](https://www.docker.com/)
- [.NET SDK](https://dotnet.microsoft.com/en-us/download) (for frontend)
- Go 1.22+ (only if running backend manually)

---

### 🐳 Run Everything with Docker

```bash
git clone https://github.com/kibecodes/logistics-system.git
cd logistics-system

# Start backend + DB via Docker
docker-compose up --build
```

Backend available at: `http://localhost:8080`  
PostgreSQL: `localhost:5432` (user: `admin`, password: `secret`, db: `logistics_db`)

---

### 🧪 Run Backend Manually (Optional)

```bash
cd apps/logistics-backend
go run main.go
```

Ensure `.env` contains:

```env
DATABASE_URL=postgres://admin:secret@localhost:5432/logistics_db?sslmode=disable
API_BASE_URL=http://localhost:8080
```

---

### 🌐 Run Frontend (Optional)

```bash
cd apps/frontend
dotnet run
```

Available at: `https://localhost:<frontend-port>`

---

## 📄 API Documentation

Swagger UI (if enabled):  
`http://localhost:8080/swagger/index.html`

---

## 🧪 API Testing with Postman & Newman

- ✅ Full Postman collection available for all major API endpoints
- ✅ Environment-based variables used for base URL flexibility
- ✅ Newman CLI runs automated tests via CI or locally

### 🔄 Switching Base URLs

Ensure your Postman environment is set to match your active setup:

| Environment | Base URL                        |
|-------------|----------------------------------|
| Docker      | `http://localhost:8080`         |
| VM / LAN    | `http://192.168.100.11:8080`    |
| Local Dev   | `http://localhost:8080` (same)  |

> Use Postman's "Environments" feature with `base_url` variable.

### ▶️ Run Newman tests locally

```bash
newman run postman/collection.json -e postman/environment.json
```

---

## 📈 Project Status

This project is under **active development**.

✅ Dockerized backend + PostgreSQL  
✅ Working Postman + Newman tests  
🔄 Full CRUD in progress  
🔌 Planning Kong API Gateway for onboarding  
📊 Upcoming dashboards and analytics

---

## 📦 Planned Enhancements

- Kong gateway config + onboarding docs
- GitHub Actions for CI
- Newman tests on PRs
- Real JWT authentication
- Driver route optimization
- Admin reports dashboard

---

## 🤝 Contributing

All contributions are welcome 🙌  
You can help by:

- Improving APIs or adding endpoints
- Fixing issues or bugs
- Enhancing UI/UX in Razor Pages
- Writing docs (README, onboarding, Postman collection)

---

## 📝 License

This project is licensed under the **MIT License**.  
See [`LICENSE`](LICENSE) for more information.
