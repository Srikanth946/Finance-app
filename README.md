# Finance-app

A comprehensive personal finance management application built with Go, utilizing a clean architecture to manage transactions and provide dashboard summaries.

## 🚀 Features
- Transaction management (Create, Read)
- Mark transactions as paid
- Dashboard summary for financial overview
- Structured logging with `zerolog`
- High-performance routing with `Gin Gonic` (In progress)


## 📂 Project Structure
```
Finance-app/
├── cmd/
│   └── finance/
│       └── main.go               # Entry point of the application
├── internal/
│   ├── controller/               # Request handling and validation
│   │   └── transaction_controller.go
│   ├── models/                   # Data structures and entities
│   │   └── transaction.go
│   ├── repository/               # Database access layer
│   │   └── transaction_repository.go
│   ├── router/                   # API route definitions
│   │   └── router.go
│   └── service/                  # Business logic layer
│       └── transaction_service.go
├── go.mod                        # Go module definitions
├── go.sum                        # Dependency checksums
└── README.md                     # Project documentation
```

## 🏁 Getting Started

### Prerequisites
- Go 1.23+
- SQLite3

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/Srikanth946/Finance-app.git
   cd Finance-app
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the application:
   ```bash
   go run cmd/finance/main.go
   ```