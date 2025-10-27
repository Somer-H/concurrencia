package main

import (
    "fmt"
    "math"
    "sync"
)

// Verifica si un número es primo
func esPrimo(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

// analiza si los números son primos
func routine(id int, nums <-chan int, resultados chan<- int, wg *sync.WaitGroup) {
	// Asegura que el WaitGroup se complete al final
    defer wg.Done()
    for n := range nums {
		// Verifica si el número es primo
        if esPrimo(n) {
            fmt.Printf("Routine %d encontró primo: %d\n", id, n)
            resultados <- n
        }
    }
}

// distribuye los datos a las routines
func enviarNumeros(numeros []int, nums chan<- int) {
	// Envía los números al canal
    for _, n := range numeros {
		// Asegura que cada número se envíe al canal
        nums <- n
    }
	// Cierra el canal al finalizar
    close(nums)
}

// espera a que los workers terminen y cierra resultados
func recolectarResultados(wg *sync.WaitGroup, resultados chan int) {
	// Espera a que todos los workers terminen
    wg.Wait()
	// Cierra el canal de resultados
    close(resultados)
}

func main() {
	// Lista de números a analizar
    numeros := []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
    numWorkers := 3
	// Canales para fan-out y fan-i
    nums := make(chan int)
    resultados := make(chan int)
    var wg sync.WaitGroup

    // Fan-Out: lanzar routines
    for i := 1; i <= numWorkers; i++ {
		// Asegura que cada routine se registre en el WaitGroup
        wg.Add(1)
		// Inicia cada routine
        go routine(i, nums, resultados, &wg)
    }

    // Enviar datos
    go enviarNumeros(numeros, nums)

    // Fan-In: cerrar resultados cuando los workers terminen
    go recolectarResultados(&wg, resultados)

    // Consolidar resultados
    fmt.Println("Primos encontrados:")
	// Recibe los resultados y los imprime
    for primo := range resultados {
        fmt.Println(primo)
    }
}