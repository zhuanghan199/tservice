package entity

import (
	"github.com/dgrijalva/jwt-go"
)

type LoginRequest struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

type LoginRedisValue struct {
	RefreshToken string `json:"refresh_token"`
	UID          string `json:"uid"`
	Name         string `json:"Name"`
	Email        string `json:"Email"`
}

type LoginResponse struct {
	UID                    string `json:"uid"`
	Email                  string `json:"email"`
	Phone                  string `json:"phone"`
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	RefreshExpireDuration  int64  `json:"refresh_expire_duration"`
	AccessExpireInDuration int64  `json:"access_expire_duration"`
}

type UserJWTClaim struct {
	UID   string `json:"UID"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
	jwt.StandardClaims
}
