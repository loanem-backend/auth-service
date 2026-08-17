package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/loanem-backend/auth-service/internal/entity"
	"github.com/loanem-backend/auth-service/internal/repository"
	"github.com/loanem-backend/auth-service/pkg/bcryptx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AssistantService interface {
	Create(ctx context.Context, a *entity.Assistant) (int, error)
	SendPasswordChangeConfirmation(ctx context.Context, oldPw, newPw, confirmPw string) error
	SetPassword(ctx context.Context, pwChangeToken string) error

	GetActiveAssistants(ctx context.Context) ([]*entity.Assistant, error)
	DeleteAssistant(ctx context.Context, assistantID int) error
}

type assistantService struct {
	assistantRepo repository.AssistantRepository
	redisRepo     repository.RedisRepository
	emailServ     EmailService
}

func NewAssistantService(ar repository.AssistantRepository, rr repository.RedisRepository, es EmailService) AssistantService {
	return &assistantService{
		assistantRepo: ar,
		redisRepo:     rr,
		emailServ:     es,
	}
}

func (s *assistantService) Create(ctx context.Context, a *entity.Assistant) (int, error) {
	hashedPassword, err := bcryptx.Hash(defaultAssistantPassword)
	if err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}
	a.HashPassword = hashedPassword

	a.Phone, err = cleanPhone(a.Phone)
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.assistantRepo.FindByPhone(ctx, a.Phone)
	if err != nil && err != repository.ErrFindByPhoneNotFound {
		return 0, status.Error(codes.Internal, err.Error())
	}
	if err == nil {
		return 0, status.Error(codes.AlreadyExists, "phone already registered")
	}

	assistantID, err := s.assistantRepo.Insert(ctx, a)
	if err != nil {
		return 0, status.Error(codes.Internal, "failed inserting assistant row")
	}

	return int(assistantID), nil
}

func (s *assistantService) SendPasswordChangeConfirmation(ctx context.Context, oldPw, newPw, confirmPw string) error {
	if confirmPw != newPw {
		return status.Error(codes.InvalidArgument, "password confirmation mismatches")
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.DataLoss, "missing metadata")
	}

	ids := meta.Get("id")

	id, err := strconv.ParseInt(ids[0], 10, 0)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	assistant, err := s.assistantRepo.FindByID(ctx, int(id))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err := bcryptx.Validate(oldPw, assistant.HashPassword); err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	hashedNewPw, err := bcryptx.Hash(newPw)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	token, err := generatePasswordChangeToken()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err := s.redisRepo.StorePasswordChange(
		ctx,
		prefixRedisPasswordChange+token,
		repository.RedisPasswordChange{
			HashedNewPassword: hashedNewPw,
			AssistantID:       assistant.ID,
		},
		5*time.Minute,
	); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err := s.emailServ.Send(ctx, emailParam{
		to:   assistant.Email,
		text: emailPasswordChange,
		args: []any{token},
	}); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func generatePasswordChangeToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed generating token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *assistantService) SetPassword(ctx context.Context, pwChangeToken string) error {
	data, err := s.redisRepo.GetPasswordChange(ctx, prefixRedisPasswordChange+pwChangeToken)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	assistant, err := s.assistantRepo.FindByID(ctx, data.AssistantID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	assistant.HashPassword = data.HashedNewPassword
	assistant.UpdatedAt = time.Now()
	if err := s.assistantRepo.UpdatePassword(ctx, assistant); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *assistantService) GetActiveAssistants(ctx context.Context) ([]*entity.Assistant, error) {
	assistants, err := s.assistantRepo.FindActiveAssistants(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return assistants, nil
}

func (s *assistantService) DeleteAssistant(ctx context.Context, assistantID int) error {
	if err := s.assistantRepo.DeleteAssistantByID(ctx, int16(assistantID)); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
