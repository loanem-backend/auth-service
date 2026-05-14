package service

import (
	"context"
	"strconv"

	"github.com/loanem-backend/auth-service/internal/repository"
	"github.com/loanem-backend/auth-service/pkg/bcryptx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AssistantService interface {
	SetPassword(ctx context.Context, oldPw, newPw, confirmPw string) error
}

type assistantService struct {
	assistantRepo repository.AssistantRepository
}

func NewAssistantService(ar repository.AssistantRepository) AssistantService {
	return &assistantService{
		assistantRepo: ar,
	}
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
	if err := s.assistantRepo.UpdatePassword(ctx, assistant); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
