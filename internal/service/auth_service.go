package service

import (
	"context"

	"github.com/loanem-backend/auth-service/internal/repository"
	"github.com/loanem-backend/auth-service/pkg/bcryptx"
	"github.com/loanem-backend/auth-service/pkg/jwtx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	// Login returns an access token, a refresh token, and an error in order.
	Login(ctx context.Context, phone, password string) (string, string, error)

	ValidateToken(ctx context.Context, token string) (*jwtx.Claims, error)
}

type authService struct {
	assistantRepo repository.AssistantRepository
}

func NewAuthService(ar repository.AssistantRepository) AuthService {
	return &authService{
		assistantRepo: ar,
	}
}

func (s *authService) Login(ctx context.Context, phone, password string) (string, string, error) {
	phoneClean, err := cleanPhone(phone)
	if err != nil {
		return "", "", status.Error(codes.InvalidArgument, err.Error())
	}

	assistant, err := s.assistantRepo.FindByPhone(ctx, phoneClean)
	if err != nil {
		return "", "", status.Error(codes.Internal, err.Error())
	}

	if err := bcryptx.Validate(password, assistant.HashPassword); err != nil {
		return "", "", status.Error(codes.PermissionDenied, err.Error())
	}

	accessToken, err := jwtx.GenerateAccessToken(assistant)
	if err != nil {
		return "", "", status.Error(codes.Internal, err.Error())
	}

	refreshToken, err := jwtx.GenerateRefreshToken(assistant.ID)
	if err != nil {
		return "", "", status.Error(codes.Internal, err.Error())
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*jwtx.Claims, error) {
	claims, err := jwtx.DecodeToken(token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return claims, nil
}
