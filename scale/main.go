package main

import (
	"fmt"
	"sync"
	"time"
)

// Simula la generación de un dato (ej. consulta a API externa)
func generateData(id int) string {
	time.Sleep(10 * time.Millisecond) // Simula latencia
	return fmt.Sprintf("Registro_%d", id)
}

// Worker: Lee jobs y envía resultados
func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		data := generateData(j)
		results <- data
	}
}

func main() {
	const numJobs = 5000
	const numWorkers = 8

	fmt.Printf("Iniciando procesamiento de %d trabajos con %d workers...\n", numJobs, numWorkers)

	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)
	var wg sync.WaitGroup

	start := time.Now()

	// 1. Iniciar Workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Enviar trabajos (Dispatcher)
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// 3. Esperar y cerrar canal de resultados
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Recolectar resultados
	count := 0
	for range results {
		count++
	}

	elapsed := time.Since(start)
	fmt.Printf("Procesados %d registros en %v\n", count, elapsed)
}
