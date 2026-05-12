package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/loanem-backend/auth-service/config"
	"github.com/loanem-backend/auth-service/internal/server"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	s := grpc.NewServer()

	server.Start(s, db)

	lis := config.InitListener()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed serving grpc: %v", err)
	}
}
