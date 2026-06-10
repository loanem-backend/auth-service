package service

import (
	"context"
	"time"

	"github.com/loanem-backend/auth-service/internal/repository"
	"github.com/loanem-backend/auth-service/pkg/bcryptx"
	"github.com/loanem-backend/auth-service/pkg/jwtx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	// Login returns an access token, a refresh token, refresh token expiration hour, and an error respectively.
	Login(ctx context.Context, phone, password string) (string, string, int32, error)

	ValidateToken(ctx context.Context, token string) (*jwtx.Claims, error)

	// RefreshToken returns a new access token and an error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	Logout(ctx context.Context, accessToken, refreshToken string) error
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

func (s *authService) Login(ctx context.Context, phone, password string) (string, string, int32, error) {
	phoneClean, err := cleanPhone(phone)
	if err != nil {
		return "", "", 0, status.Error(codes.InvalidArgument, err.Error())
	}

	assistant, err := s.assistantRepo.FindByPhone(ctx, phoneClean)
	if err != nil {
		return "", "", 0, status.Error(codes.Internal, err.Error())
	}

	if err := bcryptx.Validate(password, assistant.HashPassword); err != nil {
		return "", "", 0, status.Error(codes.PermissionDenied, err.Error())
	}

	accessToken, err := jwtx.GenerateAccessToken(assistant)
	if err != nil {
		return "", "", 0, status.Error(codes.Internal, err.Error())
	}

	refreshTokenDur, refreshToken, err := jwtx.GenerateRefreshToken(assistant.ID)
	if err != nil {
		return "", "", 0, status.Error(codes.Internal, err.Error())
	}

	if err := s.redisRepo.Store(ctx, prefixRedisRefreshToken+refreshToken, assistant.ID, refreshTokenDur); err != nil {
		// log error
	}

	return accessToken, refreshToken, int32(refreshTokenDur / time.Hour), nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*jwtx.Claims, error) {
	claims, err := jwtx.DecodeToken(token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check blacklist

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

func (s *authService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// remove refresh token from Redis
	if err := s.redisRepo.Delete(ctx, prefixRedisRefreshToken+refreshToken); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	claims, err := jwtx.DecodeToken(accessToken)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	blacklistDur := 30 * time.Minute
	if claims.ExpiresAt != nil {
		blacklistDur = time.Until(claims.ExpiresAt.Time)
	}

	// insert access token to Redis blacklist
	if err := s.redisRepo.Store(ctx, prefixRedisBlacklistToken+claims.JwtID, true, blacklistDur); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
