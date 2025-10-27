package main

import (
	"fmt"
	"sync"
	"time"
)

// Estructura del Broker Pub/Sub simple
type Broker struct {
	subscribers map[int]chan string 
	publish     chan string         
	subscribe   chan chan string    
	mu          sync.Mutex          // Mutex para proteger el mapa
}

func NewBroker() *Broker {
	// Inicializa el Broker
	b := &Broker{
		subscribers: make(map[int]chan string),
		publish:     make(chan string),
		subscribe:   make(chan chan string),
	}
	go b.run() // Inicia el goroutine del broker
	return b
}

// run: Goroutine del Broker que maneja las suscripciones y la publicación
func (b *Broker) run() {
	// ID para cada suscriptor
	subscriberID := 0
	for {
		// Espera eventos de publicación o suscripción
		select {
		case msg := <-b.publish:
			// Distribuir el mensaje a todos los suscriptores activos
			b.mu.Lock()
			for _, subCh := range b.subscribers {
				// select: Intenta enviar, si el canal está lleno, salta
				select { 
				case subCh <- msg:
				default:
					// Opcional: manejar suscriptores lentos
				}
			}
			b.mu.Unlock()

		case subCh := <-b.subscribe:
			// Nuevo suscriptor, añadir al mapa
			b.mu.Lock()
			// Incrementa el ID y registra el canal
			subscriberID++
			// Canal con buffer para el nuevo suscriptor
			b.subscribers[subscriberID] = subCh
			// Desbloquea el mutex
			b.mu.Unlock()
			// Log de registro
			fmt.Printf("Nuevo Suscriptor %d registrado.\n", subscriberID)
		}
	}
}

// Publish: Método para que un publisher envíe un mensaje
func (b *Broker) Publish(msg string) {
	b.publish <- msg
}

// Subscribe: Método para que un cliente se suscriba
func (b *Broker) Subscribe() <-chan string {
	ch := make(chan string, 1) // Canal con buffer para el suscriptor
	b.subscribe <- ch
	return ch
}

// consumer: Simula un suscriptor escuchando un canal
func consumer(id int, ch <-chan string) {
	for msg := range ch {
		fmt.Printf("    Suscriptor %d: Recibido '%s'\n", id, msg)
	}
}

func main() {
	
	broker := NewBroker() 

	// 1. Suscriptores (consumidores) se unen
	sub1Ch := broker.Subscribe()
	go consumer(1, sub1Ch) 

	sub2Ch := broker.Subscribe()
	go consumer(2, sub2Ch)

	time.Sleep(100 * time.Millisecond) 

	// 2. Publishers envían mensajes
	fmt.Println("Publicando mensaje 'Evento 1'")
	broker.Publish("Evento 1: Servidor Iniciado")
	// Pequeña pausa para permitir procesamiento
	time.Sleep(50 * time.Millisecond)
	// Otro mensaje
	fmt.Println("Publicando mensaje 'Evento 2'")
	broker.Publish("Evento 2: Tarea Finalizada")

	time.Sleep(500 * time.Millisecond)
}