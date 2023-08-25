package entity

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}


type RefreshTokenResponse struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	RefreshExpireDuration  int64  `json:"refresh_expire_duration"`
	AccessExpireInDuration int64  `json:"access_expire_duration"`
}