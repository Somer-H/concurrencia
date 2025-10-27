package main

import (
	"context"
	"fmt"
	"time"
)

// workerConContext: Una goroutine que respeta la señal de cancelación del Context
func workerConContext(ctx context.Context, id int) {
	fmt.Printf("  Worker %d: trabajando...", id)
	// Loop de trabajo que verifica el contexto
	for {
		// Verifica si el contexto ha sido cancelado
		select {
		case <-ctx.Done(): // El canal Done se cierra si el Context es cancelado
			err := ctx.Err()
			fmt.Printf("  Worker %d: Finalizando graciosamente. Razón: %v\n", id, err)
			return // La goroutine termina
		default:
			// Simula trabajo continuo
			fmt.Printf("  Worker %d: Ejecutando tarea...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// 1. Context con Timeout (para Deadline)
	ctxDeadline, cancelDeadline := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	// Aseguramos la cancelación al final
	defer cancelDeadline()
	// Iniciar un worker que respete el contexto
	go workerConContext(ctxDeadline, 1)
	// Esperar a que el contexto expire
	select {
	// Caso de éxito
	case <-ctxDeadline.Done():
		// Context finalizado por timeout
		fmt.Println("Main: Context 1 finalizado. Razón:", ctxDeadline.Err())
		// Caso de espera prolongada
	case <-time.After(2 * time.Second):
	}

	// 2. Context con Cancelación Jerárquica
	ctxPadre, cancelPadre := context.WithCancel(context.Background())
	// Aseguramos la cancelación al final
	defer cancelPadre()
	// Crear un contexto hijo
	ctxHijo, _ := context.WithCancel(ctxPadre)
	// Iniciar workers que respeten los contextos
	go workerConContext(ctxPadre, 2)
	go workerConContext(ctxHijo, 3)
	// Dejar que trabajen un poco
	time.Sleep(1 * time.Second)
	fmt.Println("Cancelando Contexto Padre")
	cancelPadre()
	// Esperar un poco para ver la finalización
	time.Sleep(1 * time.Second)
	fmt.Println("Control de Concurrencia por Context finalizado.")
}
