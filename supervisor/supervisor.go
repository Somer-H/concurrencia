package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// workerFailable: Una goroutine que falla aleatoriamente
func workerFailable(ctx context.Context, id int) {
	// Indica que el worker ha iniciado
	fmt.Printf("  Worker %d: Iniciado.\n", id)
	// Simula trabajo con posibilidad de fallo
	for {
		// Verificar si el contexto ha sido cancelado
		select {
			// Si el contexto fue cancelado, terminar la goroutine
		case <-ctx.Done():
			// Propagar el error de contexto
			fmt.Printf("  Worker %d: Cancelado por el Supervisor. Terminando.\n", id)
			// Terminar la goroutine
			return
		default:
			// Simulación de trabajo con fallo aleatorio
			time.Sleep(500 * time.Millisecond)
			
			if rand.Intn(10) < 3 { // Falla 30% de las veces
				fmt.Printf("  Worker %d: Terminando inesperadamente.\n", id)
				return // La goroutine muere aquí
			}
			// Trabajo exitoso
			fmt.Printf("  Worker %d: Trabajo.\n", id)
		}
	}
}

// Supervisor: Monitorea y relanza el Worker
func supervisor(ctx context.Context, id int) {
	for {
		// Indica que el supervisor está lanzando un worker
		fmt.Printf("[Supervisor %d]: Lanzando Worker %d...\n", id, id)
		// Crear un WaitGroup para esperar al worker
		var wg sync.WaitGroup
		// Incrementar el contador del WaitGroup
		wg.Add(1)
		// Lanzar el worker en una goroutine
		go func() {
			defer wg.Done()
			workerFailable(ctx, id)
		}()
		
		// Esperar a que el worker termine
		wg.Wait()
		
		// Verificar el contexto para saber si fue un fallo o una cancelación externa
		select {
			// Si el contexto fue cancelado, terminar el supervisor
		case <-ctx.Done():
			// Propagar el error de contexto
			fmt.Printf("Supervisor %d: Recibió señal de cancelación externa.\n", id)
			// Terminar el supervisor
			return
		default:
			// Si el worker terminó sin cancelación del contexto, fue un fallo. ¡Relanzar!
			fmt.Printf("[Supervisor %d]: Worker %d falló.\n", id, id)
			// Esperar un momento antes de relanzar
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// Crear un contexto con cancelación para el supervisor
	ctx, cancel := context.WithCancel(context.Background())
	// Asegurar la cancelación al final
	defer cancel()
	// Lanzar el supervisor
	go supervisor(ctx, 1)

	// Dejar que corra por un tiempo para ver los fallos y reinicios
	time.Sleep(5 * time.Second) 
	
	// Cancelar el Supervisor para una parada elegante
	fmt.Println("Main: Enviando señal de cancelación al Supervisor.")
}