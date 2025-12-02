package resilience

import (
	"errors"
	"testing"
)

// Transaction represents the data structure to be processed
type Transaction struct {
	ID     string
	Amount int
}

// ProcessTransaction is the System Under Test (SUT)
func ProcessTransaction(t Transaction) error {
	// Business Logic Validation
	if t.Amount < 0 {
		return errors.New("invalid amount: cannot be negative")
	}
	// Additional logic...
	return nil
}

func TestResilienceWithMutation(t *testing.T) {
	// 1. Generator: Create a VALID base record
	validData := Transaction{ID: "TX-001", Amount: 100}

	// 2. Mutator: Introduce a specific defect
	mutatedData := validData
	mutatedData.Amount = -100 // Harmful mutation

	// 3. Quality Gate: Verify that the system is robust and rejects the data
	err := ProcessTransaction(mutatedData)

	if err == nil {
		t.Errorf("RESILIENCE FAILURE: System accepted a transaction with negative amount.")
	} else {
		t.Logf("SUCCESS: System correctly rejected mutated data. Error: %v", err)
	}
}
