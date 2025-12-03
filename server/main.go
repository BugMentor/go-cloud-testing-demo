package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/lib/pq" // Postgres Driver
)

// Global DB Connection Pool
var db *sql.DB

// Order represents the data structure for an e-commerce order
type Order struct {
	ID        string  `json:"id"`
	Item      string  `json:"item"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Processed bool    `json:"processed"`
}

func main() {
	// 1. Connect to Real Database (PostgreSQL running in Docker)
	// Credentials matches docker-compose.yml: user=admin, password=password123
	var err error
	connStr := "postgres://admin:password123@localhost:5432/ecommerce?sslmode=disable"

	// IF USING LOCAL WINDOWS POSTGRES (Uncomment this line instead):
	// connStr := "postgres://postgres:password123@localhost:5432/ecommerce?sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ FATAL: Could not open DB connection: %v", err)
	}

	// Validate connection with a Ping
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ FATAL: Database is unreachable. Is Docker running? Error: %v", err)
	}
	fmt.Println("✅ Connection to PostgreSQL established successfully.")

	// 2. Initialize Schema (Auto-create table)
	initDB()

	// 3. Start API Server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders/", handleOrderRequest)

	serverPort := ":8080"
	fmt.Printf("\n🚀 REAL E-COMMERCE API STARTED\n")
	fmt.Printf("📡 Listening on port %s\n", serverPort)
	fmt.Printf("📝 Waiting for traffic...\n\n")

	if err := http.ListenAndServe(serverPort, mux); err != nil {
		log.Fatal("Server crashed:", err)
	}
}

func initDB() {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		item TEXT,
		amount DECIMAL,
		status TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database schema: %v", err)
	}
	fmt.Println("✅ Database schema initialized (Table 'orders' ready).")
}

func handleOrderRequest(w http.ResponseWriter, r *http.Request) {
	// Extract Order ID from URL
	id := r.URL.Path[len("/api/orders/"):]

	// Business Logic: Generate random order details
	amount := float64(rand.Intn(1000)) + 50.0
	tax := amount * 0.21
	total := amount + tax
	item := fmt.Sprintf("Product-%d", rand.Intn(999))

	// --- REAL DATABASE I/O ---
	// We perform a synchronous INSERT. High concurrency here tests the DB connection pool.
	insertQuery := `INSERT INTO orders (id, item, amount, status) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(insertQuery, id, item, total, "CONFIRMED")

	if err != nil {
		// Log the specific DB error to console (e.g., connection limit reached)
		// log.Printf("⚠️ DB Error inserting order %s: %v", id, err)
		http.Error(w, "Database Write Failed", http.StatusInternalServerError)
		return
	}

	// Prepare JSON Response
	response := Order{
		ID:        id,
		Item:      item,
		Amount:    total,
		Status:    "CONFIRMED",
		Processed: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
