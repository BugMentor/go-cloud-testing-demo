package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	_ "github.com/lib/pq" // Postgres Driver
)

// Order represents the DB record structure
type Order struct {
	ID        string
	Item      string
	Amount    float64
	Status    string
	CreatedAt string
}

func main() {
	// 1. Connect to Database
	// Ensure these credentials match your docker-compose or local setup
	connStr := "postgres://admin:password123@localhost:5432/ecommerce?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Could not connect to DB: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Database unreachable: %v", err)
	}

	fmt.Println("\n🔎 STARTING DATA VERIFICATION...")
	fmt.Println("------------------------------------------------")

	// 2. Get Total Count
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	if err != nil {
		log.Fatalf("❌ Failed to count records: %v", err)
	}

	fmt.Printf("📊 TOTAL RECORDS IN DB: %d\n", count)
	fmt.Println("------------------------------------------------")

	// 3. Fetch Last 10 Records to verify data integrity
	rows, err := db.Query("SELECT id, item, amount, status, created_at FROM orders ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		log.Fatalf("❌ Failed to query records: %v", err)
	}
	defer rows.Close()

	// Use tabwriter for pretty printing the table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID\t ITEM\t AMOUNT\t STATUS\t CREATED_AT")
	fmt.Fprintln(w, "--\t ----\t ------\t ------\t ----------")

	for rows.Next() {
		var o Order
		// Note: created_at comes as a string from the driver by default for simplicity
		if err := rows.Scan(&o.ID, &o.Item, &o.Amount, &o.Status, &o.CreatedAt); err != nil {
			log.Printf("⚠️ Error scanning row: %v", err)
			continue
		}
		fmt.Fprintf(w, "%s\t %s\t $%.2f\t %s\t %s\n", o.ID, o.Item, o.Amount, o.Status, o.CreatedAt)
	}
	w.Flush()
	fmt.Println("\n✅ Verification Complete.")
}
