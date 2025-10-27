package main

import (
    "fmt"
    "strings"
    "sync"
)

// Etapa 1: recibe nombres
func etapaExtract(nombres []string, out chan<- string) {
	// Envía cada nombre al canal
    for _, nombre := range nombres {
        out <- nombre
    }
	// Cierra el canal al finalizar
    close(out)
}

// Etapa 2: transforma nombres a mayúsculas
func etapaTransform(in <-chan string, out chan<- string) {
	// Transforma cada nombre a mayúsculas y lo envía al siguiente canal
    for nombre := range in {
		// Convierte el nombre a mayúsculas
        out <- strings.ToUpper(nombre)
    }
	// Cierra el canal al finalizar
    close(out)
}

// Etapa 3: guarda nombres en lista final
func etapaLoad(in <-chan string, resultado *[]string, wg *sync.WaitGroup) {
	// Asegura que el WaitGroup se complete al final
    defer wg.Done()
    for nombre := range in {
		// Guarda el nombre en la lista final
        fmt.Println("Guardando:", nombre)
		// Agrega el nombre a la lista
        *resultado = append(*resultado, nombre)
    }
}

func main() {
	// Lista de nombres a procesar
    nombres := []string{"somer", "joseph", "schmetterling", "carolinne", "charly"}
	// Canales para la comunicación entre etapas
    ch1 := make(chan string)
    ch2 := make(chan string)
	// Lista para almacenar el resultado final
    var resultado []string
	// WaitGroup para esperar a que todas las etapas terminen
    var wg sync.WaitGroup
	// Añade una rutina al WaitGroup para la etapa de carga
    wg.Add(1)
	// Inicia las etapas en goroutines
    go etapaExtract(nombres, ch1)
    go etapaTransform(ch1, ch2)
    go etapaLoad(ch2, &resultado, &wg)
	// Espera a que la etapa de carga termine
    wg.Wait()
    fmt.Println("Resultado final:", resultado)
}