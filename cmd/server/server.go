package server

import (
	"auth/config"
	"auth/generated/auth"
	"auth/logs"
	"auth/service"
	"auth/storage"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer(storage storage.IStorage) {
	logs.InitLogger()

	cfg := config.Load()

	listener, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		logs.Logger.Error("listener error: ", err)
		log.Fatal(err)
	}

	s := grpc.NewServer()

	server := service.AuthService{
		Log:  logs.Logger,
		User: storage,
	}

	auth.RegisterAuthServiceServer(s, &server)

	logs.Logger.Info("Starting gRPC server on port: ","PORT", cfg.GRPC_PORT)

	log.Print("server is running ", "PORT", cfg.GRPC_PORT)
	err = s.Serve(listener)

	if err != nil {
		logs.Logger.Error("server error: ", err)
		log.Fatal(err)
	}

}
