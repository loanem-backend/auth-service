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
	accessToken, refreshToken, refreshExpHour, err := s.serv.Login(ctx, req.GetPhone(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return mapper.StringsToLoginResponse(accessToken, refreshToken, refreshExpHour), nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pbauth.ValidateTokenRequest) (*pbauth.ValidateTokenResponse, error) {
	claimsData, err := s.serv.ValidateToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return mapper.ClaimsToValidateTokenResponse(claimsData), nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *pbauth.RefreshTokenRequest) (*pbauth.RefreshTokenResponse, error) {
	accessToken, refreshToken, refreshTokenDur, err := s.serv.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &pbauth.RefreshTokenResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshExpirationHour: refreshTokenDur,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *pbauth.LogoutRequest) (*pbauth.LogoutResponse, error) {
	if err := s.serv.Logout(ctx, req.GetAccessToken(), req.GetRefreshToken()); err != nil {
		return nil, err
	}

	return &pbauth.LogoutResponse{}, nil
}
