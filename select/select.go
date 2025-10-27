package main

import (
	"context"
	"fmt"
	"time"
)

// longOperation: Simula una tarea que tarda 3 segundos
func longOperation(ctx context.Context, data chan string) {
	time.Sleep(3 * time.Second) 
	// Verificamos el contexto después de la operación
	select {
	case <-ctx.Done():
		// La operación fue cancelada
		fmt.Println("  Operación Larga: ¡Cancelada antes de enviar el resultado!")
		return
	default:
		// Enviar el resultado si no fue cancelada
		data <- "Resultado de Operación Larga"
	}
}

func main() {
	// 1. Caso Timeout
	fmt.Println("\nCaso 1: Timeout 1s ")
	chTimeout := make(chan string)
	go func() {
		time.Sleep(2 * time.Second) // Tarea de 2 segundos
		chTimeout <- "Resultado Rápido  2s"
	}()

	select {
	case res := <-chTimeout:
		// Recibido antes del timeout
		fmt.Printf("  Recibido: %s\n", res)
	case <-time.After(1 * time.Second): // Espera solo 1 segundo
		fmt.Println("El resultado no llegó en 1 segundo.")
	}


	// 2. Caso Cancelación por Context
	fmt.Println("\nCaso 2: Cancelación Jerárquica por Context ")
	// Crear un contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	// Aseguramos la cancelación al final
	defer cancel() 
	// Canal para recibir el resultado
	dataChannel := make(chan string, 1)
	// Iniciar la operación larga
	go longOperation(ctx, dataChannel) 
	// Esperar el resultado o cancelar después de 1 segundo
	select {
	case res := <-dataChannel:
		// Recibido antes del timeout
		fmt.Printf("  Recibido: %s\n", res)
	case <-time.After(1 * time.Second):
		// Timeout alcanzado, cancelar la operación
		fmt.Println(" Timeout preliminar 1s, cancelando operación...")
		// Cancelar el contexto
		cancel() 
		// Esperar un poco para demostrar la cancelación
		time.Sleep(50 * time.Millisecond) 
	}
	// Final del programa
	fmt.Println("Cancelación por Context completada.")
}