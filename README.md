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

* **Go** (version 1.18 or higher)
* **Mockgen** (GoMock tool for mock generation)

To install `mockgen`, run:
```bash
go install github.com/golang/mock/mockgen@latest
````

-----

## 📦 Installation

1.  Initialize the module (if you haven't already):
    ```bash
    go mod tidy
    ```

-----

## 🏛 The Three Pillars

### 1\. Isolation: Mocking Dependencies

[cite_start]**Problem:** Direct dependencies on infrastructure (like databases) make tests slow and fragile[cite: 45, 47].
**Solution:** We use **Interfaces** and **`gomock`** to isolate business logic. [cite_start]This allows us to simulate edge cases (e.g., DB errors) in milliseconds without a real database[cite: 54, 126].

  * **Key Concept:** Dependency Injection.
  * **Tool:** `github.com/golang/mock/gomock`.

### 2\. Scale: Worker Pools

[cite_start]**Problem:** Testing microservices often requires generating massive amounts of synthetic data, which is slow with traditional loops[cite: 81].
**Solution:** We use the **Worker Pool pattern** to leverage Go's concurrency. [cite_start]This demo generates thousands of records in parallel using Goroutines and Channels[cite: 88, 128].

  * **Key Concept:** Concurrency (Fan-out / Fan-in).
  * **Performance:** Drastically reduces execution time compared to sequential processing.

### 3\. Resilience: Mutation Testing

**Problem:** Passing tests don't always mean the system is robust against bad data.
**Solution:** We use **Mutation Testing**. [cite_start]We take valid data, "mutate" it (introduce a defect), and ensure the **Quality Gate** (our system) rejects it[cite: 114, 120].

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

  * **Expected Output:** A log indicating the processing of 5,000 jobs and the total execution time (approx. \~6.5s for 8 workers).

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
