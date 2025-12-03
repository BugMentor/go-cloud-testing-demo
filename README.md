# Cloud Testing with Go: SDET Demos

This repository implements the three pillars of a modern **Cloud SDET** architecture: **Isolation, Scale, and Resilience**. It demonstrates how to transition from fragile, slow automation to robust engineering using Go.

Based on the concepts from "Cloud Testing with Go".

## 📋 Table of Contents

- [Cloud Testing with Go: SDET Demos](#cloud-testing-with-go-sdet-demos)
  - [📋 Table of Contents](#-table-of-contents)
  - [🛠 Prerequisites](#-prerequisites)
  - [📦 Installation](#-installation)
  - [🏛 The Three Pillars](#-the-three-pillars)
    - [1. Isolation: Mocking Dependencies](#1-isolation-mocking-dependencies)
    - [2. Scale: Full-Stack Load Testing](#2-scale-full-stack-load-testing)
    - [3. Resilience: Mutation Testing](#3-resilience-mutation-testing)
  - [✅ Data Verification](#-data-verification)
  - [⚙️ CI/CD Pipeline](#️-cicd-pipeline)
  - [🚀 How to Run](#-how-to-run)
    - [Step 1: Generate Mocks](#step-1-generate-mocks)
    - [Step 2: Run Unit Tests](#step-2-run-unit-tests)
    - [Step 3: Run Full-Stack Load Test](#step-3-run-full-stack-load-test)
  - [📂 Project Structure](#-project-structure)

---

## 🛠 Prerequisites

Before running the demos, ensure you have the following installed:

* **Go** (version 1.21 or higher)
* **Docker & Docker Compose** (for the real database integration)
* **Mockgen** (GoMock tool for mock generation)

To install `mockgen`, run:
```bash
go install github.com/golang/mock/mockgen@latest
````

-----

## 📦 Installation

1.  Initialize the module and download dependencies:
    ```bash
    go mod tidy
    ```

-----

## 🏛 The Three Pillars

### 1\. Isolation: Mocking Dependencies

**Problem:** Direct dependencies on infrastructure (like databases) make tests slow and fragile.
**Solution:** We use **Interfaces** and **`gomock`** to isolate business logic. This allows us to simulate edge cases (e.g., DB errors) in milliseconds without a real database.

  * **Key Concept:** Dependency Injection.
  * **Tool:** `github.com/golang/mock/gomock`.

### 2\. Scale: Full-Stack Load Testing

**Problem:** Testing microservices often requires generating massive amounts of synthetic data, which is slow with traditional loops.
**Solution:** We use the **Worker Pool pattern** to leverage Go's concurrency. This demo generates thousands of records in parallel using Goroutines and Channels.

**Real-World Implementation:**
Instead of a simulation, this project now runs a **Full-Stack Load Test**:

1.  **Load Client (`scale/`):** A high-performance generator using 100+ concurrent workers.
2.  **Real API Server (`server/`):** An HTTP server handling validation and business logic.
3.  **Database (`postgres`):** A Dockerized PostgreSQL instance storing real data.

<!-- end list -->

  * **Key Concept:** Concurrency (Fan-out / Fan-in) & Connection Pooling.
  * **Performance:** Validates system throughput (RPS) and database write capacity.

### 3\. Resilience: Mutation Testing

**Problem:** Passing tests don't always mean the system is robust against bad data.
**Solution:** We use **Mutation Testing**. We take valid data, "mutate" it (introduce a defect), and ensure the **Quality Gate** (our system) rejects it.

  * **Flow:** Generator (Valid) -\> Mutator (Invalid) -\> Assertion (Expect Error).

-----

## ✅ Data Verification

After running the load test, it is critical to verify data integrity. We implemented a **Verification Tool** (`verify/`) that:

1.  Connects to the PostgreSQL database.
2.  Counts the total records created.
3.  Fetches a sample of the last inserted rows for auditing.

-----

## ⚙️ CI/CD Pipeline

The repository includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that automatically:

1.  **Spins up a Service Container:** A real PostgreSQL instance for testing.
2.  **Runs Unit Tests:** Checks Isolation and Resilience logic.
3.  **Builds Binaries:** Compiles the Server and Load Client.
4.  **Smoke Test:** Runs a quick integration test to ensure the system starts correctly.

-----

## 🚀 How to Run

Follow these steps to execute the complete Cloud SDET workflow.

### Step 1: Generate Mocks

Generate the mock implementations for the Isolation demo:

```bash
go generate ./isolation/...
```

### Step 2: Run Unit Tests

Execute the unit tests (Isolation & Resilience):

```bash
go test ./isolation/... ./resilience/... -v
```

### Step 3: Run Full-Stack Load Test

This demo requires 3 terminals to simulate a real microservice environment.

**1. Start Infrastructure (Docker):**
Start the PostgreSQL database:

```bash
docker-compose up -d
```

**2. Start the API Server (Terminal 1):**
This connects to the DB and listens for traffic.

```bash
go run server/main.go
```

*Wait for: `🚀 REAL E-COMMERCE API STARTED`*

**3. Run the Load Generator (Terminal 2):**
This simulates 10,000 concurrent users.

```bash
go run scale/main.go
```

**4. Verify Results (Terminal 3):**
Audit the database to confirm data integrity.

```bash
go run verify/main.go
```

-----

## 📂 Project Structure

```text
go-cloud-testing-demo/
├── .github/                # CI/CD Pipelines
│   └── workflows/ci.yml    # GitHub Actions workflow
├── docker-compose.yml      # Database Infrastructure
├── go.mod                  # Module definition
├── isolation/              # PILLAR 1: Isolation
│   ├── mocks/              # Auto-generated mocks
│   ├── user.go             # Service & Interface definition
│   └── user_test.go        # Tests using Gomock
├── scale/                  # PILLAR 2: Scale (Load Client)
│   └── main.go             # High-concurrency load generator
├── server/                 # REAL API SERVER
│   └── main.go             # HTTP Server + Postgres Connection
├── verify/                 # DATA AUDIT
│   └── main.go             # Verification tool
└── resilience/             # PILLAR 3: Resilience
    └── mutation_test.go    # Mutation testing logic
```

-----

**Author:** Ing. Matías J. Magni (CEO @ BugMentor).
