package service

import (
	"time"
	"tservice/common/logger"
	"tservice/config"
	"tservice/db"
	. "tservice/models"
	. "tservice/services/auth/entity"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"encoding/json"
)

// 解析accesstoken的jwtclaim
func (s AuthService) decodeJWT(req RefreshTokenRequest) (*UserJWTClaim, *ErrorInfo) {
	jwtClaim := UserJWTClaim{}
	_, err := jwt.ParseWithClaims(req.AccessToken, &jwtClaim, func(toke *jwt.Token) (i interface{}, e error) {
		return config.Cnf.JWTKey, nil
	})
	logger.Debugln(jwtClaim)
	if len(jwtClaim.UID) > 0 {
		return &jwtClaim, nil
	} else {
		logger.Errorln(err.Error())
		return nil, &ErrorInfo{ERR_CODE_REFRESH_TOKEN_FAILED, "无效的access token"}
	}
}

// 查询redis，比较数据
func (s AuthService) verifyToken(claim *UserJWTClaim, req RefreshTokenRequest) *ErrorInfo {
	if redisConn, err := db.NewRedisConn(); err != nil {
		return &ErrorInfo{ERR_CODE_REFRESH_TOKEN_FAILED, "redis打开失败"}
	} else {
		defer redisConn.Close()
		logger.Debugln(claim)
		if b, ok := (redisConn.Get(claim.UID)).([]byte); ok{
			var value LoginRedisValue
			if err := json.Unmarshal(b, &value); err == nil{
				logger.Debugln(value)
				if claim.Email == value.Email &&
					claim.Name == value.Name &&
					req.RefreshToken == value.RefreshToken {
					return nil
				}
			}else{
				logger.Debugln(err.Error())
			}
		}
		return &ErrorInfo{ERR_CODE_REFRESH_TOKEN_FAILED, "token校验失败"}
	}
}

// 重新生成refreshToken和accessToken，存入redis
func (s AuthService) genNewToken(claim *UserJWTClaim) (*RefreshTokenResponse, *ErrorInfo) {
	claim.ExpiresAt = time.Now().Add(time.Second * time.Duration(config.Cnf.JWTDuration)).Unix()
	logger.Debugln(config.Cnf.JWTKey)
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtKeyByte := []byte(config.Cnf.JWTKey)
	if token, err := tokenObj.SignedString(jwtKeyByte); err != nil {
		return nil, &ErrorInfo{ERR_CODE_REQUEST_FAILED, "生成token失败"}
	} else {
		refreshToken := uuid.New().String()
		refreshExpireDuration := int64(config.Cnf.JWTRefreshDuration)
		// refreshToken存redis
		if redisConn, err := db.NewRedisConn(); err != nil {
			return nil, &ErrorInfo{ERR_CODE_REGISTER_FAILED, "redis打开失败"}
		} else {
			defer redisConn.Close()
			redisValue := LoginRedisValue{
				RefreshToken: refreshToken,
				UID:          claim.UID,
				Name:         claim.Name,
				Email:        claim.Email,
			}
			if err := redisConn.ExpiresSave(claim.UID, redisValue, refreshExpireDuration); err != nil {
				return nil, &ErrorInfo{ERR_CODE_REGISTER_FAILED, "redis保存失败"}
			} else {
				resp := RefreshTokenResponse{
					RefreshToken:           refreshToken,
					RefreshExpireDuration:  refreshExpireDuration,
					AccessToken:            token,
					AccessExpireInDuration: int64(config.Cnf.JWTDuration),
				}
				return &resp, nil
			}
		}

	}
}
