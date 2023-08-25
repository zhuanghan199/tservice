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
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
)

// 查询用户名是否在db，并验证pwd
func (s AuthService) verifyLogin(req LoginRequest) (*User, *ErrorInfo) {
	if db, err := db.OpenDatabase(); err != nil {
		return nil, &ErrorInfo{ERR_CODE_DB_EXCEPTION, "数据库异常"}
	} else {
		defer db.Close()
		var user User
		if err := db.DB().Where("name = ?", req.Name).First(&user).Error; err != nil {
			return nil, &ErrorInfo{ERR_CODE_ACCOUNT_NOT_EXIST, "帐号不存在"}
		} else {
			// 验证密码是否相同
			password := []byte(req.Pwd)
			hashedPassword := []byte(user.Pwd)
			if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
				return nil, &ErrorInfo{ERR_CODE_PWD_ERROR, "密码错误"}
			} else {
				return &user, nil
			}
		}
	}
}

// 生成token
func (s AuthService) genToken(usr *User) (*LoginResponse, *ErrorInfo) {
	claim := UserJWTClaim{
		UID:   usr.UID,
		Name:  usr.Name,
		Email: usr.Email,
	}
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
				UID:          usr.UID,
				Name:         usr.Name,
				Email:        usr.Email,
			}
			paramsStr, _ := json.Marshal(redisValue)
			if err := redisConn.ExpiresSave(usr.UID, paramsStr, refreshExpireDuration); err != nil {
				return nil, &ErrorInfo{ERR_CODE_REGISTER_FAILED, "redis保存失败"}
			} else {
				loginResp := LoginResponse{
					UID:                    usr.UID,
					Email:                  usr.Email,
					Phone:                  usr.Phone,
					RefreshToken:           refreshToken,
					RefreshExpireDuration:  refreshExpireDuration,
					AccessToken:            token,
					AccessExpireInDuration: int64(config.Cnf.JWTDuration),
				}
				return &loginResp, nil
			}
		}

	}

}
