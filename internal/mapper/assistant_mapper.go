package mapper

import (
	"github.com/loanem-backend/auth-service/internal/entity"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateAssistantRequestToAssistant(req *pbauth.CreateAssistantRequest) *entity.Assistant {
	return &entity.Assistant{
		Name:   req.GetName(),
		Phone:  req.GetPhone(),
		Period: int(req.GetPeriod()),
		Email:  req.GetEmail(),
	}
}

func AssistantToPbAssistant(assistant *entity.Assistant) *pbauth.Assistant {
	return &pbauth.Assistant{
		Id:        int32(assistant.ID),
		Name:      assistant.Name,
		Phone:     assistant.Phone,
		Period:    int32(assistant.Period),
		Active:    assistant.Active,
		Email:     assistant.Email,
		CreatedAt: timestamppb.New(assistant.CreatedAt),
		UpdatedAt: timestamppb.New(assistant.UpdatedAt),
	}
}

func AssistantsToGetActiveAssistantsResponse(assistants []*entity.Assistant) *pbauth.GetActiveAssistantsResponse {
	pbAssistants := make([]*pbauth.Assistant, len(assistants))

	for i, a := range assistants {
		pbAssistants[i] = AssistantToPbAssistant(a)
	}

	return &pbauth.GetActiveAssistantsResponse{
		Assistants: pbAssistants,
	}
}
