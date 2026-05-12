package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loanem-backend/auth-service/infra/database/sqlc"
	"github.com/loanem-backend/auth-service/internal/repository"
	"github.com/loanem-backend/auth-service/internal/service"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"google.golang.org/grpc"
)

func Start(s *grpc.Server, p *pgxpool.Pool) {
	registerServers(s, p)
}

func registerServers(s *grpc.Server, p *pgxpool.Pool) {
	queries := sqlc.New(p)

	var (
		assistantRepo = repository.NewAssistantRepository(queries)
	)

	var (
		assistantServ = service.NewAuthService(assistantRepo)
	)

	pbauth.RegisterAuthServiceServer(s, NewAuthServer(assistantServ))
}
