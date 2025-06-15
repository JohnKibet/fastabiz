# 🚚 logistics-system

A scalable logistics & operations management system with clean APIs, Dockerized microservices, Kong API Gateway, JWT security, and automated API tests.

---

## 🚀 Features

- 🔐 Role-based access control (Admin, Driver, Customer)
- 📦 Full CRUD for orders, deliveries, drivers, payments, feedbacks, notifications
- 🐳 Dockerized backend + Postgres + Kong Gateway
- 🌐 Swagger documentation, served via Kong proxy
- 🧪 Postman + Newman API tests in GitHub Actions CI
- 🔑 Kong plugins: rate limiting & JWT auth
- 🔜 Planned: gRPC/Kafka communication, frontend, production deployment

---

## 🛠️ Tech Stack

| Layer       | Technologies                                 |
|-------------|----------------------------------------------|
| Backend     | Go (Chi, Clean Architecture, Swagger)        |
| Gateway     | Kong (JWT auth + rate limiting)              |
| Database    | PostgreSQL                                   |
| CI/CD       | GitHub Actions + Newman (Docker mode)        |
| Containerization | Docker, Docker Compose                 |

---

## 📁 Repository Structure

```
logistics-system/
├── apps/
│   └── logistics-backend/       # Go APIs
├── kong/
│   └── kong.yml                 # Kong declarative config
├── postman/
│   ├── collection.json          # API test collection
│   └── environment.json         # API test environment
├── .github/
│   └── workflows/
│       └── api-tests.yml       # CI config
├── .env.docker                  # Docker environment variables
├── Dockerfile                   # Backend Dockerfile
├── docker-compose.yml          # Compose services
└── README.md
```

---

## ⚙️ Getting Started

### Prerequisites

- Docker & Docker Compose
- Git

---

### 🚀 Running Locally with Docker

```bash
git clone https://github.com/kibecodes/logistics-system.git
cd logistics-system

# Start all services: DB, backend, Kong
docker compose up --build
```

- **API & Swagger:** `http://localhost:8000/api/swagger/index.html`
- **Backend logs:** `docker compose logs -f backend`
- **Kong Admin:** `http://localhost:8001`

---

### 🧪 Running Postman Tests Locally

```bash
docker run --rm   -v "${PWD}/postman:/etc/newman"   postman/newman:alpine run collection.json   --environment=environment.json   --reporters cli
```

---

## 🧩 Environment Configuration

**.env.docker**

```env
PUBLIC_API_BASE_URL=http://localhost:8000/api
INTERNAL_API_BASE_URL=http://backend:8080
DATABASE_URL=postgres://admin:secret@db:5432/logistics_db?sslmode=disable
PORT=8080
JWT_SECRET=<your-secret>
```

Kong connects to the backend on `http://backend:8080` internally, while clients use `localhost:8000`.

---

## 🛡️ Kong Configuration

- `/api/swagger` route: public + rate-limiting
- `/api/*` route: JWT-protected + rate-limiting
- Consumer `test-user` with JWT secret in `kong.yml`

---

## 🔐 JWT & Swagger

Swagger UI uses `@securityDefinitions.apikey JWT`, allowing you to authorize with a valid token (issue via `/api/users/login`) to test protected endpoints interactively.

---

## 📈 CI with GitHub Actions

The CI workflow (`.github/workflows/api-tests.yml`) uses:

- Docker-based Newman to run Postman tests
- Environment variables defined in `postman/environment.json`
- Outputs test artifacts in JSON and HTML formats

---

## ⏭️ Next Steps

- Implement business logic: orders, drivers, routes
- Integrate gRPC / Kafka for inter-service communication
- Add frontend (ASP.NET or other)
- Enable production CI/CD, monitoring, and deployment

---

## 🤝 Contributing

Your contributions are welcome! Suggested areas:

- Completing business logic and clean architecture layers
- Adding frontend user interfaces or dashboards
- Production-grade logging, monitoring, and gateway enhancements
- Message bus integrations (Kafka / RabbitMQ)

---

## 📝 License

MIT License – see [LICENSE](LICENSE)

---