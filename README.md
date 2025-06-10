# logistics-system

A scalable logistics & operations management system designed for multi-role (Admin, Driver, Customer) interactions. Built using Go for backend APIs, C# for the frontend with ASP.NET Razor Pages, and PostgreSQL as the database.

---

## 🚀 Features

- 🔐 Role-based access control (Admin, Driver, Customer)
- 📦 Order management: create, assign, and track orders
- 🚚 Delivery assignment logic
- 🧑‍💼 User authentication and authorization (token-based)
- 📄 API documentation via Swagger / OpenAPI (in progress)
- 🐳 Dockerized PostgreSQL database setup
- 📊 Planned: dashboard analytics, reporting, and API gateway integration

---

## 🛠️ Tech Stack

| Layer       | Technologies                              |
|-------------|-------------------------------------------|
| Backend     | Go (Chi router, Clean Architecture)       |
| Frontend    | ASP.NET Core Razor Pages (C#)             |
| Database    | PostgreSQL                                |
| DevOps      | Docker, Docker Compose                    |
| Auth        | Token-based (currently mocked)            |
| Docs        | Swagger / OpenAPI (in progress)           |

---

## 🗂️ Architecture

```
apps/
├── backend/     # Go APIs (modular by domain)
├── frontend/    # Razor Pages (C#)
└── migrations/  # SQL migration scripts
```

Follows **Clean Architecture** principles:
- `usecase/`, `repository/`, and `delivery/` layers for separation of concerns.
- Decoupled request/response DTOs for validation and transformation.

---

## ⚙️ Getting Started

```bash
# Clone the repo
git clone https://github.com/kibecodes/logistics-system.git
cd logistics-system

# Start PostgreSQL DB using Docker Compose
docker-compose up -d

# Run the backend (Go)
cd apps/backend
go run main.go

# Run the frontend (ASP.NET)
cd apps/frontend
dotnet run
```

Swagger documentation (WIP):  
`http://localhost:<port>/swagger`

---

## 📈 Project Status

This project is **under active development**. Current focus areas:

- ✅ Robust error handling
- 🔄 Order workflow implementation
- 📄 API documentation and validation
- 🔌 Gateway integration planning (Kong, WSO2)

---

## 🤝 Contributing

Contributions and feedback are welcome! You can help by:

- Improving APIs or request handling
- Enhancing UI/UX in Razor Pages
- Suggesting or adding gateway/auth integrations

---

## 📝 License

This project is licensed under the **MIT License**. See [`LICENSE`](LICENSE) for more information.