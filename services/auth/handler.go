package auth
import(
	. "tservice/services/auth/service"
	. "tservice/services/auth/endpoint"
	. "tservice/services/auth/transport"

	httpt "github.com/go-kit/kit/transport/http"
)
func AuthHandler() *httpt.Server{
	authService := AuthService{}
	authEndpoint := GenAuthEndPoint(authService)
	authHandler := httpt.NewServer(authEndpoint, DecodeAuthRequest, EncodeAuthResponse)
	return authHandler
}
