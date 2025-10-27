package main

import (
	"fmt"
	"sync"
	"time"
)

// Simula una tarea que toma un tiempo y devuelve un resultado
type Tarea struct {
	ID int
}

// Función del worker: Consume tareas, las procesa y envía el resultado
func worker(id int, tareas <-chan Tarea, resultados chan<- string) {
	// Procesa cada tarea recibida
	for tarea := range tareas {
		// Simula procesamiento
		fmt.Printf("  Worker %d: Iniciando tarea %d\n", id, tarea.ID)
		time.Sleep(time.Duration(tarea.ID%3) * 500 * time.Millisecond) // Simula tiempo variable
		resultado := fmt.Sprintf("Tarea %d completada por Worker %d", tarea.ID, id)
		// Simula finalización
		fmt.Printf("  Worker %d: Finalizando tarea %d\n", id, tarea.ID)
		resultados <- resultado
	}
}

func recibirResultados(numTareas int, resultados chan string, wg *sync.WaitGroup) {
		for i := 0; i < numTareas; i++ {
			fmt.Printf("  Main: Recibido: %s\n", <-resultados)
			wg.Done() // Decrementa el contador por cada resultado recibido
		}
		close(resultados)
	}
func main() {
	// Configuración del Worker Pool
	numTareas := 9
	numWorkers := 3
	// Canales para tareas y resultados
	tareas := make(chan Tarea, numTareas) 
	resultados := make(chan string, numTareas) 
	var wg sync.WaitGroup

	// 1. Iniciar los Workers 
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tareas, resultados)
	}

	// 2. Enviar las Tareas a la cola
	for i := 1; i <= numTareas; i++ {
		tareas <- Tarea{ID: i}
		wg.Add(1) // Contamos la tarea enviada
	}
	close(tareas) // Cerrar el canal de tareas

	// 3. Recibir los Resultados
	go recibirResultados(numTareas, resultados, &wg)
	// Esperar a que todas las tareas se completen
	wg.Wait()
	fmt.Println("Worker Pool: Todas las tareas procesadas y resultados recogidos.")
}