package main

import (
	"context"
	"log"
	"net"

	pb "server/proto" // Importamos los protos locales

	"google.golang.org/grpc"
)

// Implementamos el servicio definido en el .proto
type server struct {
	pb.UnimplementedWarReportServiceServer
}

func (s *server) SendReport(ctx context.Context, req *pb.WarReportRequest) (*pb.WarReportResponse, error) {
	log.Printf("Reporte recibido por gRPC: País=%v, Aviones=%d, Barcos=%d", req.Country, req.WarplanesInAir, req.WarshipsInWater)

	// AQUÍ IRA LA LÓGICA DE RABBITMQ EN LA FASE 4

	return &pb.WarReportResponse{
		Status:  "success",
		Message: "Reporte procesado exitosamente por el Servidor gRPC",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar en puerto 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWarReportServiceServer(s, &server{})

	log.Println("Servidor Go (gRPC) escuchando en puerto 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar gRPC: %v", err)
	}
}
