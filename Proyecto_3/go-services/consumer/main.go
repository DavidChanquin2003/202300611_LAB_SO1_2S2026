package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-svc:5672/")
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error abriendo canal: %v", err)
	}
	defer ch.Close()

	// Declarar la misma cola para asegurar que exista
	q, err := ch.QueueDeclare("war_reports", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}

	// Consumir los mensajes
	msgs, err := ch.Consume(
		q.Name, // nombre de la cola
		"",     // consumer
		true,   // auto-ack (confirma que ya se leyó y lo borra de la cola)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Error registrando consumidor: %v", err)
	}

	// Canal para mantener el contenedor vivo
	forever := make(chan os.Signal, 1)
	signal.Notify(forever, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range msgs {
			log.Printf("Nuevo reporte consumido desde RabbitMQ: %s", d.Body)
		}
	}()

	log.Println("Consumidor Go iniciado. Esperando reportes...")
	<-forever
}
