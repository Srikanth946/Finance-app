# Finance-app

A comprehensive personal finance management application built with Go, utilizing a clean architecture to manage transactions and provide dashboard summaries.

## 🚀 Features
- **AI Finance Agent**:
  - Intelligent conversational interface using Google Gemini
  - Custom Tooling for real-time financial data retrieval
  - Configurable Human-in-the-Loop (HITL) confirmation flow
- **Transaction Management**:
  - Create new transactions
  - Retrieve all transactions or a specific transaction by ID (mobile number)
  - Mark transactions as paid
  - Extend loans with new interest rates and compound durations
- **Interest Calculation**:
  - General interest calculation based on amount, rate, and duration
  - Transaction-specific interest calculation
- **Financial Dashboard**:
  - Get a summary overview of finances
- **Technical Stack**:
  - High-performance routing with `Gin Gonic`
  - Structured logging with `zerolog`
  - Lightweight data storage with `SQLite`
  - AI Framework: `google.golang.org/adk/v2`


## 📂 Project Structure
```
Finance-app/
├── cmd/
│   └── finance/
│       └── main.go               # Entry point of the application
├── internal/
│   ├── controller/               # Request handling and validation
│   │   ├── dashboard_controller.go
│   │   ├── interest_controller.go
│   │   └── transaction_controller.go
│   ├── models/                   # Data structures and entities
│   │   └── transaction.go
│   ├── repository/               # Database access layer
│   │   └── transaction_repository.go
│   ├── router/                   # API route definitions
│   │   └── router.go
│   └── service/                  # Business logic layer
│       ├── dashboard_service.go
│       ├── intrest_service.go
│       └── transaction_service.go
├── my_agent/                     # AI Agent implementation
│   ├── agent.go                  # Agent entry point & configuration
│   ├── tools/                    # Tool definitions and implementations
│   │   ├── tool_impl.go          # Custom tool implementations
│   │   └── tools.go              # Tool helper functions
│   ├── utils/                    # Agent utilities and configuration
│   │   └── config.go             # Tool and agent configuration
│   └── models/                   # Agent-specific models
│       └── gemini.go
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

### 🤖 Running the AI Agent
1. Ensure the main finance application is running.
2. Navigate to the agent directory:
   ```bash
   cd my_agent
   ```
3. Run the agent in console mode:
   ```bash
   go run . console
   ```