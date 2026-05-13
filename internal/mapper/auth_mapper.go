package mapper

import (
	"github.com/loanem-backend/auth-service/pkg/jwtx"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

func StringToLoginResponse(s string) *pbauth.LoginResponse {
	return &pbauth.LoginResponse{
		Token: s,
	}
}

func ClaimsToValidateTokenResponse(c *jwtx.Claims) *pbauth.ValidateTokenResponse {
	return &pbauth.ValidateTokenResponse{
		IsValid: true,
		UserId:  int32(c.ID),
		Name:    c.Name,
	}
}
