package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Order struct {
	ID        string  `json:"id"`
	Item      string  `json:"item"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Processed bool    `json:"processed"`
}

func main() {
	// 1. Connection String (Docker Credentials)
	connStr := "postgres://admin:password123@localhost:5432/ecommerce?sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ FATAL: Could not open DB connection: %v", err)
	}

	// --- FIX 1: DATABASE CONNECTION POOLING ---
	// Crucial for Load Testing: Limit open connections to avoid killing Postgres.
	// Postgres default max_connections is usually 100. We stay safely under that.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ FATAL: DB unreachable: %v", err)
	}
	fmt.Println("✅ Connected to PostgreSQL successfully (Pool configured).")

	initDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders/", handleOrderRequest)

	serverPort := ":8080"
	fmt.Printf("\n🚀 REAL E-COMMERCE API STARTED on %s\n", serverPort)

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
		log.Fatalf("❌ Failed to init schema: %v", err)
	}
}

func handleOrderRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/orders/"):]

	amount := float64(rand.Intn(1000)) + 50.0
	tax := amount * 0.21
	total := amount + tax
	item := fmt.Sprintf("Product-%d", rand.Intn(999))

	// Insert into DB
	insertQuery := `INSERT INTO orders (id, item, amount, status) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(insertQuery, id, item, total, "CONFIRMED")

	if err != nil {
		// --- FIX 2: LOGGING ---
		// Print the ACTUAL error to the server console to verify the issue
		log.Printf("⚠️ DB Error on Order %s: %v", id, err)

		http.Error(w, "Database Write Failed", http.StatusInternalServerError)
		return
	}

	response := Order{ID: id, Item: item, Amount: total, Status: "CONFIRMED", Processed: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
