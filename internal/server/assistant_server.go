package server

import (
	"context"

	"github.com/loanem-backend/auth-service/internal/service"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

type AssistantServer struct {
	pbauth.UnimplementedAssistantServiceServer
	serv service.AssistantService
}

func NewAssistantServer(as service.AssistantService) *AssistantServer {
	return &AssistantServer{
		serv: as,
	}
}

func (s *AssistantServer) SetPassword(ctx context.Context, req *pbauth.SetAssistantPasswordRequest) (*pbauth.SetAssistantPasswordResponse, error) {
	if err := s.serv.SetPassword(ctx, req.GetOldPassword(), req.GetNewPassword(), req.GetConfirmPassword()); err != nil {
		return nil, err
	}

	return &pbauth.SetAssistantPasswordResponse{}, nil
}
