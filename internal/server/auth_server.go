package server

import (
	"context"

	"github.com/loanem-backend/auth-service/internal/mapper"
	"github.com/loanem-backend/auth-service/internal/service"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

type AuthServer struct {
	pbauth.UnimplementedAuthServiceServer
	authServ service.AuthService
}

func NewAuthServer(as service.AuthService) *AuthServer {
	return &AuthServer{
		authServ: as,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	tokenData, err := s.authServ.Login(ctx, req.GetPhone(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return mapper.StringToLoginResponse(tokenData), nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pbauth.ValidateTokenRequest) (*pbauth.ValidateTokenResponse, error) {
	claimsData, err := s.authServ.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return mapper.ClaimsToValidateTokenResponse(claimsData), nil
}
