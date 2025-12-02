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
    - [2. Scale: Worker Pools](#2-scale-worker-pools)
    - [3. Resilience: Mutation Testing](#3-resilience-mutation-testing)
  - [🚀 How to Run](#-how-to-run)
    - [Step 1: Generate Mocks](#step-1-generate-mocks)
    - [Step 2: Run Tests (Isolation \& Resilience)](#step-2-run-tests-isolation--resilience)
    - [Step 3: Run the Scale Demo](#step-3-run-the-scale-demo)
  - [📂 Project Structure](#-project-structure)

---

## 🛠 Prerequisites

Before running the demos, ensure you have the following installed:

* **Go** (version 1.16 or higher)
* **Mockgen** (GoMock tool for mock generation)

To install `mockgen`, run:
```bash
go install [github.com/golang/mock/mockgen@latest](https://github.com/golang/mock/mockgen@latest)
````

-----

## 📦 Installation

1.  Clone the repository:

    ```bash
    git clone [https://github.com/BugMentor/go-cloud-testing-demo.git](https://github.com/BugMentor/go-cloud-testing-demo.git)
    cd go-cloud-testing-demo
    ```

2.  Download dependencies:

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

### 2\. Scale: Worker Pools

**Problem:** Testing microservices often requires generating massive amounts of synthetic data, which is slow with traditional loops.
**Solution:** We use the **Worker Pool pattern** to leverage Go's concurrency. This demo generates thousands of records in parallel using Goroutines and Channels.

  * **Key Concept:** Concurrency (Fan-out / Fan-in).
  * **Performance:** Drastically reduces execution time compared to sequential processing.

### 3\. Resilience: Mutation Testing

**Problem:** Passing tests don't always mean the system is robust against bad data.
**Solution:** We use **Mutation Testing**. We take valid data, "mutate" it (introduce a defect), and ensure the **Quality Gate** (our system) rejects it.

  * **Flow:** Generator (Valid) -\> Mutator (Invalid) -\> Assertion (Expect Error).

-----

## 🚀 How to Run

Follow these steps to execute all demonstrations.

### Step 1: Generate Mocks

First, generate the mock implementations for the Isolation demo using the `go generate` directive found in `isolation/user.go`:

```bash
go generate ./isolation/...
```

### Step 2: Run Tests (Isolation & Resilience)

Execute the unit tests to verify the Isolation logic (mocks) and the Resilience logic (mutation):

```bash
go test ./... -v
```

  * **Expected Output:** You should see `PASS` for `TestGetUserName` (Isolation) and `TestResilienceWithMutation` (Resilience).

### Step 3: Run the Scale Demo

Run the `main.go` file to see the Worker Pool in action generating synthetic data:

```bash
go run scale/main.go
```

  * **Expected Output:** A log indicating the processing of 5,000 jobs and the total execution time.

-----

## 📂 Project Structure

```text
go-cloud-testing-demo/
├── go.mod                  # Module definition
├── isolation/              # PILLAR 1: Isolation
│   ├── mocks/              # Auto-generated mocks (do not edit)
│   ├── user.go             # Service & Interface definition
│   └── user_test.go        # Tests using Gomock
├── scale/                  # PILLAR 2: Scale
│   └── main.go             # Worker Pool implementation for data generation
└── resilience/             # PILLAR 3: Resilience
    └── mutation_test.go    # Mutation testing logic
```

-----

**Author:** Ing. Matías J. Magni (CEO @ BugMentor).
