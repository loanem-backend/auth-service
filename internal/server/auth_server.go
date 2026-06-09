package server

import (
	"context"

	"github.com/loanem-backend/auth-service/internal/mapper"
	"github.com/loanem-backend/auth-service/internal/service"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

type AuthServer struct {
	pbauth.UnimplementedAuthServiceServer
	serv service.AuthService
}

func NewAuthServer(as service.AuthService) *AuthServer {
	return &AuthServer{
		serv: as,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	accessToken, refreshToken, err := s.serv.Login(ctx, req.GetPhone(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return mapper.StringsToLoginResponse(accessToken, refreshToken), nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pbauth.ValidateTokenRequest) (*pbauth.ValidateTokenResponse, error) {
	claimsData, err := s.serv.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return mapper.ClaimsToValidateTokenResponse(claimsData), nil
}
