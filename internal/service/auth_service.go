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

	// RefreshToken returns a new access token and an error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type authService struct {
	assistantRepo repository.AssistantRepository
	redisRepo     repository.RedisRepository
}

func NewAuthService(ar repository.AssistantRepository, rr repository.RedisRepository) AuthService {
	return &authService{
		assistantRepo: ar,
		redisRepo:     rr,
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

	refreshTokenDur, refreshToken, err := jwtx.GenerateRefreshToken(assistant.ID)
	if err != nil {
		return "", "", status.Error(codes.Internal, err.Error())
	}

	if err := s.redisRepo.Store(ctx, prefixRedisRefreshToken+refreshToken, assistant.ID, refreshTokenDur); err != nil {
		// log error
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

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	assistantID, err := s.redisRepo.GetInt(ctx, prefixRedisRefreshToken+refreshToken)
	if err != nil {
		return "", status.Error(codes.DeadlineExceeded, "refresh token has been expired")
	}

	assistant, err := s.assistantRepo.FindByID(ctx, assistantID)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	accessToken, err := jwtx.GenerateAccessToken(assistant)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return accessToken, nil
}
