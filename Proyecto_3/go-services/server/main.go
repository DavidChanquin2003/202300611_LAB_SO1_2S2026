package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	pb "server/proto"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWarReportServiceServer
}

// Estructura para enviar a la cola
type QueueMessage struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

func (s *server) SendReport(ctx context.Context, req *pb.WarReportRequest) (*pb.WarReportResponse, error) {
	log.Printf("Reporte recibido por gRPC: País=%v, Aviones=%d, Barcos=%d", req.Country, req.WarplanesInAir, req.WarshipsInWater)

	// Conectar a RabbitMQ (usamos el nombre del Service de K8s)
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-svc:5672/")
	if err != nil {
		log.Printf("Error conectando a RabbitMQ: %v", err)
		return &pb.WarReportResponse{Status: "error", Message: "Fallo conexión RabbitMQ"}, nil
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Error abriendo canal RabbitMQ: %v", err)
		return &pb.WarReportResponse{Status: "error", Message: "Fallo canal RabbitMQ"}, nil
	}
	defer ch.Close()

	// Declarar la cola
	q, err := ch.QueueDeclare("war_reports", true, false, false, false, nil)
	if err != nil {
		log.Printf("Error declarando cola: %v", err)
		return &pb.WarReportResponse{Status: "error", Message: "Fallo cola RabbitMQ"}, nil
	}

	// Preparar mensaje
	msg := QueueMessage{
		Country:         req.Country.String(),
		WarplanesInAir:  req.WarplanesInAir,
		WarshipsInWater: req.WarshipsInWater,
		Timestamp:       req.Timestamp,
	}
	body, _ := json.Marshal(msg)

	// Publicar en la cola
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	if err != nil {
		log.Printf("Error publicando mensaje: %v", err)
		return &pb.WarReportResponse{Status: "error", Message: "Fallo publicando en RabbitMQ"}, nil
	}

	log.Println("Mensaje publicado en RabbitMQ exitosamente")
	return &pb.WarReportResponse{
		Status:  "success",
		Message: "Reporte procesado y encolado por el Servidor gRPC",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar en puerto 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWarReportServiceServer(s, &server{})

	log.Println("Servidor Go (gRPC Writer) escuchando en puerto 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar gRPC: %v", err)
	}
}
