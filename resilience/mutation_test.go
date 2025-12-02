package resilience

import (
	"errors"
	"testing"
)

// Transaction representa el dato a procesar
type Transaction struct {
	ID     string
	Amount int
}

// ProcessTransaction es el Sistema Bajo Prueba (SUT)
func ProcessTransaction(t Transaction) error {
	if t.Amount < 0 {
		return errors.New("monto inválido: no puede ser negativo")
	}
	// Lógica adicional...
	return nil
}

func TestResilienceWithMutation(t *testing.T) {
	// 1. Generador: Crea un registro VÁLIDO base
	validData := Transaction{ID: "TX-001", Amount: 100}

	// 2. Mutator: Introduce un defecto específico
	mutatedData := validData
	mutatedData.Amount = -100 // Mutación dañina

	// 3. Quality Gate: Verificamos que el sistema sea robusto y rechace el dato
	err := ProcessTransaction(mutatedData)

	if err == nil {
		t.Errorf("FALLO DE RESILIENCIA: El sistema aceptó una transacción con monto negativo.")
	} else {
		t.Logf("ÉXITO: El sistema rechazó correctamente el dato mutado. Error: %v", err)
	}
}
