package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup" // Requiere go get golang.org/x/sync/errgroup
)

// workerControl: Goroutine de trabajo que puede fallar
func workerControl(ctx context.Context, id int, job int) (string, error) {
	fmt.Printf("  [Worker %d]: Procesando Job %d\n", id, job)
	
	// Simular un fallo crítico en el Job 5
	if job == 5 {
		// Simula tiempo antes de fallar
		time.Sleep(1 * time.Second)
		// Retorna un error crítico
		return "", errors.New(fmt.Sprintf("Job %d falló: Error crítico", job))
	}

	// Simular una tarea normal
	time.Sleep(time.Duration(job%4) * 200 * time.Millisecond) 
	
	// Verificar la cancelación del contexto (si otro job falló)
	select {
		// Si el contexto fue cancelado, retornar el error
	case <-ctx.Done():
		// Retornar el error del contexto
		return "", ctx.Err()
		// Si no fue cancelado, retornar el resultado exitoso
	default:
		// Retornar resultado exitoso
		return fmt.Sprintf("Resultado: Job %d completado", job), nil
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8} // El 5 debe fallar
	numWorkers := 3 // Control: Limitamos a 3 goroutines concurrentes
	// Canal limitador de concurrencia
	limiter := make(chan struct{}, numWorkers) 

	// Fan-out/Controller: Lanzar goroutines
	for _, job := range jobs {
		// Capturar variable para la goroutine
		job := job 
		// Control de concurrencia
		limiter <- struct{}{} // Bloquea si hay numWorkers goroutines activas
		// Lanzar la goroutine de trabajo
		g.Go(func() error {
			// Liberar el espacio en el limitador al finalizar
			defer func() { <-limiter }() 
			// Ejecutar el worker
			res, err := workerControl(ctx, job, job)
			//
			if err != nil {
				fmt.Printf("  Worker %d:Error %v\n", job, err)
				return err // Retornar el error para cancelar todo el grupo
			}
			// Éxito
			fmt.Printf("  Worker %d: Éxito. %s\n", job, res)
			return nil
		})
	}

	// Fan-in/Controller: Esperar. g.Wait() retorna el primer error y cancela el ctx
	fmt.Println("Esperando a que el Fan Controller finalice...")
	if err := g.Wait(); err != nil {
		// Hubo un error crítico
		fmt.Printf("Fan Controller finalizado por error: %v\n", err)
	} else {
		// Todos los trabajos completaron exitosamente
		fmt.Println("Todos los trabajos completados sin errores.")
	}
}