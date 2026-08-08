package service

import (
	"context"
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
	SetPassword(ctx context.Context, oldPw, newPw, confirmPw string) error

	GetActiveAssistants(ctx context.Context) ([]*entity.Assistant, error)
}

type assistantService struct {
	assistantRepo repository.AssistantRepository
}

func NewAssistantService(ar repository.AssistantRepository) AssistantService {
	return &assistantService{
		assistantRepo: ar,
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

func (s *assistantService) SetPassword(ctx context.Context, oldPw, newPw, confirmPw string) error {
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

	assistant.HashPassword = hashedNewPw
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
