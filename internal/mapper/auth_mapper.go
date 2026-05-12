package mapper

import pbauth "github.com/loanem-backend/protos/pb/proto/services/auth/v1"

func StringToLoginResponse(s string) *pbauth.LoginResponse {
	return &pbauth.LoginResponse{
		Token: s,
	}
}
