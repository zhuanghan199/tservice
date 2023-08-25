package endpoint

import (
	"context"

	"tservice/common/logger"
	. "tservice/models"
	. "tservice/services/auth/entity"
	. "tservice/services/auth/service"

	"encoding/json"
	"github.com/go-kit/kit/endpoint"
)

// 类似于flutter中的entity模块，请求和响应的model,
func GenAuthEndPoint(service IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		httpReq := req.(HttpRequest)
		logger.Debugln("auth endpoint: ", httpReq)
		paramsStr, err := json.Marshal(httpReq.Params); 
		if err != nil{
			return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "params errro"), nil
		}
		switch httpReq.Method {
		case "login":
			var loginReq LoginRequest
			if err := json.Unmarshal(paramsStr, &loginReq); err != nil {
				return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "params errro"), nil
			}
			logger.Debugln(loginReq)
			return login(service, loginReq)
		case "register":
			var registerReq RegisterRequest
			if err := json.Unmarshal(paramsStr, &registerReq); err != nil {
				return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "params errro"), nil
			}
			logger.Debugln(registerReq)
			return register(service, registerReq)
		case "refresh_token":
			var rtReq RefreshTokenRequest
				if err := json.Unmarshal(paramsStr, &rtReq); err != nil {
					return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "params errro"), nil
				}else{
					logger.Debugln(httpReq.Params)
					logger.Debugln(rtReq)
					return refreshToken(service, rtReq)
				}
		default:
			return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "params errro"), nil
		}
	}
}

func login(service IAuthService, req interface{}) (HttpResponse, error) {
	r := req.(LoginRequest)
	if len(r.Name) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no username"), nil
	}
	if len(r.Pwd) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no pwd"), nil
	}
	loginRsp, errInfo := service.Login(r)
	logger.Debugln(loginRsp)
	logger.Debugln(errInfo)
	if errInfo != nil {
		return NewErrorResonseWithInfo(errInfo), nil
	} else {
		return NewSuccessResonse(loginRsp), nil
	}
}

func register(service IAuthService, req interface{}) (HttpResponse, error) {
	r := req.(RegisterRequest)
	if len(r.Name) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no username"), nil
	}
	if len(r.Pwd) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no pwd"), nil
	}
	if len(r.Email) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no email"), nil
	}
	if len(r.Phone) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no phone"), nil
	}

	errInfo := service.Register(r)
	logger.Debugln(errInfo)
	if errInfo != nil {
		return NewErrorResonseWithInfo(errInfo), nil
	} else {
		return NewSuccessResonse(nil), nil
	}
}

func refreshToken(service IAuthService, req interface{}) (HttpResponse, error) {
	r := req.(RefreshTokenRequest)
	if len(r.RefreshToken) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no refresh_token"), nil
	}
	if len(r.AccessToken) == 0 {
		return NewErrorResonse(ERR_CODE_INVALID_PARAMS, "no access_token"), nil
	}
	refreshTokenRsp, errInfo := service.RefreshToken(r)
	logger.Debugln(refreshTokenRsp)
	logger.Debugln(errInfo)
	if errInfo != nil {
		return NewErrorResonseWithInfo(errInfo), nil
	} else {
		return NewSuccessResonse(refreshTokenRsp), nil
	}
}
