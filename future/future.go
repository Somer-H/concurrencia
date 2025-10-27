package main

import (
	"fmt"
	"time"
)

// Future representa el resultado que estará disponible en el futuro
type Future <-chan int

// calcularAsincrono: Inicia el cálculo en una goroutine y retorna un Future (canal)
func calcularAsincrono(a, b int) Future {
	resultado := make(chan int, 1) // Canal con buffer 1 para no bloquear el envío
	
	go func() {
		// Simula una tarea larga
		fmt.Printf("Iniciando cálculo %d + %d\n", a, b)
		// Simula tiempo de cálculo		
		time.Sleep(2 * time.Second)
		// Realiza el cálculo 
		suma := a + b
		// Muestra el resultado
		fmt.Printf("Cálculo finalizado. Resultado: %d\n", suma)
		resultado <- suma // Envía el resultado al canal
	}()

	return resultado
}

func main() {
	// Iniciar la tarea asíncrona
	future1 := calcularAsincrono(5, 7)
	
	// El programa principal puede seguir haciendo otras cosas
	fmt.Println("Siguiendo con otras tareas mientras se calcula...")
	// Simula otra tarea
	time.Sleep(500 * time.Millisecond)
	// Tarea intermedia
	fmt.Println("Tarea intermedia completada.")
	
	// Esperar por el resultado del Future (esto bloquea hasta que el resultado esté en el canal)
	fmt.Println("Esperando el resultado de Future...")
	// Recibe el resultado del Future
	resultado1 := <-future1 
	// Muestra el resultado recibido
	fmt.Printf("Resultado recibido de Future: %d\n", resultado1)
}