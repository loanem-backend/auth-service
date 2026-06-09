package mapper

import (
	"github.com/loanem-backend/auth-service/pkg/jwtx"
	pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"
)

// StringsToLoginResponse accepts an access token and a refresh token in order.
func StringsToLoginResponse(at, rt string) *pbauth.LoginResponse {
	return &pbauth.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}
}

func ClaimsToValidateTokenResponse(c *jwtx.Claims) *pbauth.ValidateTokenResponse {
	return &pbauth.ValidateTokenResponse{
		IsValid: true,
		UserId:  int32(c.ID),
		Name:    c.Name,
	}
}
