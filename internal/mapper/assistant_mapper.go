package mapper

import (
	"github.com/loanem-backend/auth-service/internal/entity"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

func CreateAssistantRequestToAssistant(req *pbauth.CreateAssistantRequest) *entity.Assistant {
	return &entity.Assistant{
		Name:   req.GetName(),
		Phone:  req.GetPhone(),
		Period: int(req.GetPeriod()),
	}
}
