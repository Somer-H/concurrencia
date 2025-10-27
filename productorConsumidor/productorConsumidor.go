package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Pedido representa un pedido de cliente
type Pedido struct {
    ID            int
    Cliente       string
    Producto      string
    Cantidad      int
    TiempoLlegada time.Time
}

// Productor: simula pedidos llegando al sistema
func recibirPedidos(pedidos chan<- Pedido, numPedidos int) {
    clientes := []string{"Schmetterling", "David", "Ali", "Joseph", "Somer"}
    productos := []string{"lap", "mouse", "teclado", "monitor", "audífonos"}

    for i := 1; i <= numPedidos; i++ {
        pedido := Pedido{
            ID:            i,
            Cliente:       clientes[rand.Intn(len(clientes))],
            Producto:      productos[rand.Intn(len(productos))],
            Cantidad:      rand.Intn(5) + 1,
            TiempoLlegada: time.Now(),
        }
        fmt.Printf("%d: %s ordenó x%d %s\n", pedido.ID, pedido.Cliente, pedido.Cantidad, pedido.Producto)
        pedidos <- pedido
        time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+200))
    }

    close(pedidos)
    fmt.Println("No hay más pedidos entrantes")
}

// Consumidor: procesa los pedidos y los guarda en lista protegida por mutex
func procesarPedidos(id int, pedidos <-chan Pedido, wg *sync.WaitGroup, mtx *sync.Mutex, procesados *[]Pedido) {
    defer wg.Done()
    for pedido := range pedidos {
        fmt.Printf("  Procesador %d: Iniciando pedido #%d de %s\n", id, pedido.ID, pedido.Cliente)
        tiempoProceso := time.Millisecond * time.Duration(rand.Intn(800)+500)
        time.Sleep(tiempoProceso)
        tiempoTotal := time.Since(pedido.TiempoLlegada)
        fmt.Printf("  Procesador %d: Pedido #%d completado en %v\n", id, pedido.ID, tiempoTotal.Round(time.Millisecond))

        // Guardar en lista compartida
        mtx.Lock()
        *procesados = append(*procesados, pedido)
        mtx.Unlock()
    }
    fmt.Printf("consumidor %d terminado\n", id)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    pedidos := make(chan Pedido, 5)
    var wg sync.WaitGroup
    var mtx sync.Mutex
    var procesados []Pedido

    fmt.Println("Sistema de Pedidos Iniciado")

    // Iniciar consumidores
    wg.Add(3)
    for i := 1; i <= 3; i++ {
        go procesarPedidos(i, pedidos, &wg, &mtx, &procesados)
    }

    // Iniciar productor
    go recibirPedidos(pedidos, 15)

    // Esperar a que los consumidores terminen
    wg.Wait()

    fmt.Println("Sistema de pedidos finalizado")
    fmt.Printf("Total procesados: %d\n", len(procesados))
}