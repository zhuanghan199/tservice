package service

import (
	"tservice/common/logger"
	. "tservice/models"
	. "tservice/services/auth/entity"
)

/* 遗留： 
 * 1. 手机号/邮箱的验证 2. 需要登录验证码
 */

// / 类似于flutter的api，负责实现api
type IAuthService interface {
	Login(req LoginRequest) (*LoginResponse, *ErrorInfo)
	Register(req RegisterRequest) *ErrorInfo
	RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, *ErrorInfo)
}
type AuthService struct {
}

// 登录
// 生成token
func (s AuthService) Login(req LoginRequest) (*LoginResponse, *ErrorInfo) {
	logger.Debugln("login user:", req.Name)
	if usr, err := s.verifyLogin(req); err != nil {
		logger.Warningln("verify error: ", err.Msg)
		return nil, err
	} else {
		if loginRsp, err := s.genToken(usr); err != nil {
			logger.Warningln("genToken error: ", err.Msg)
			return nil, err
		} else {
			return loginRsp, nil
		}
	}
}

// 注册
// 存入表中
func (s AuthService) Register(req RegisterRequest) *ErrorInfo {
	logger.Debugln("Register user:", req.Name)
	if err := s.verifyRegiser(req); err != nil {
		logger.Errorln("verify error: ", err.Msg)
		return err
	} else {
		if err := s.insertUser(req); err != nil {
			logger.Errorln("insertUser error: ", err.Msg)
			return err
		} else {
			return nil
		}
	}
}

// token刷新
func (s AuthService) RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, *ErrorInfo) {
	logger.Debugln("refreshToken")
	if claim, err := s.decodeJWT(req); err != nil {
		logger.Warningln("verify error: ", err.Msg)
		return nil, err
	} else {
		if err := s.verifyToken(claim, req); err != nil {
			logger.Warningln("genToken error: ", err.Msg)
			return nil, err
		} else {
			if rsp, err := s.genNewToken(claim); err != nil {
				logger.Warningln("genToken error: ", err.Msg)
				return nil, err
			} else {
				return rsp, nil
			}
		}
	}
}