package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loanem-backend/auth-service/infra/database/sqlc"
	"github.com/loanem-backend/auth-service/internal/repository"
	"github.com/loanem-backend/auth-service/internal/service"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"github.com/redis/go-redis/v9"
	"github.com/resend/resend-go/v3"
	"google.golang.org/grpc"
)

func Start(s *grpc.Server, p *pgxpool.Pool, rc *redis.Client, emailClient *resend.Client) {
	registerServers(s, p, rc, emailClient)
}

func registerServers(s *grpc.Server, p *pgxpool.Pool, rc *redis.Client, ec *resend.Client) {
	queries := sqlc.New(p)

	var (
		assistantRepo = repository.NewAssistantRepository(queries)
		redisRepo     = repository.NewRedisRepository(rc)
	)

	var (
		emailServ     = service.NewEmailService(ec)
		authServ      = service.NewAuthService(assistantRepo, redisRepo)
		assistantServ = service.NewAssistantService(assistantRepo, redisRepo, emailServ)
	)

	pbauth.RegisterAuthServiceServer(s, NewAuthServer(authServ))
	pbauth.RegisterAssistantServiceServer(s, NewAssistantServer(assistantServ))
}
