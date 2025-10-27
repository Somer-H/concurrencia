package main

import (
	"fmt"
	"sync"
	"time"
)

// Función de trabajo en etapas
func etapa1(id int, wg *sync.WaitGroup) {
	
	fmt.Printf("Routine %d: Iniciando Etapa ", id)
	time.Sleep(time.Duration(id) * 100 * time.Millisecond) // Simula tiempo variable
	fmt.Printf("Routine %d: Etapa 1 Finalizada", id)
	wg.Done() // Señaliza la finalización de la Etapa 1
}

func etapa2(id int, barrera2 *sync.WaitGroup) {
			defer barrera2.Done()
			fmt.Printf("Routine %d: Ejecutando Etapa 2", id)
			time.Sleep(200 * time.Millisecond)
		}

func main() {
	numRoutines := 4
	var barrera1 sync.WaitGroup // Barrera para la Etapa 1

	// 1. Configurar la barrera y lanzar trabajadores
	barrera1.Add(numRoutines)
	// Lanzar las goroutines para la Etapa 1
	for i := 1; i <= numRoutines; i++ {
		go etapa1(i, &barrera1)
	}

	// 2. La goroutine principal espera en la barrera
	fmt.Println("Esperando a que todos los trabajadores completen la Etapa 1")
	barrera1.Wait() // Bloquea hasta que el contador llegue a cero

	// 3. Continuar a la siguiente etapa
	fmt.Println("Iniciando Etapa 2 de forma concurrente.")

	// Etapa 2:
	var barrera2 sync.WaitGroup
	// Configurar la barrera para la Etapa 2
	barrera2.Add(numRoutines)
	// Lanzar las goroutines para la Etapa 2
	for i := 1; i <= numRoutines; i++ {
		go etapa2(i, &barrera2)
	}
	// Esperar a que todos completen la Etapa 2
	barrera2.Wait()
	// Finalización
	fmt.Println("Simulación con Barrera Completada.")
}