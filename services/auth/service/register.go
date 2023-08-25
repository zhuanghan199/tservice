package service

import (
	"time"
	"tservice/db"
	. "tservice/models"
	. "tservice/services/auth/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// 查表，检测是否重复注册
func (s AuthService) verifyRegiser(req RegisterRequest) (*ErrorInfo) {
	if db, err := db.OpenDatabase(); err != nil {
		return &ErrorInfo{ERR_CODE_DB_EXCEPTION, "数据库异常"}
	} else {
		defer db.Close()
		if err := db.DB().Where("name = ?", req.Name).First(&User{}).Error; err == nil {
			return &ErrorInfo{ERR_CODE_REGISTER_FAILED, "帐号已存在"}
		} else {
			if err := db.DB().Where("email = ?", req.Email).First(&User{}).Error; err == nil {
				return &ErrorInfo{ERR_CODE_REGISTER_FAILED, "邮箱已被注册过"}
			} else{
				if err := db.DB().Where("phone = ?", req.Phone).First(&User{}).Error; err == nil {
					return &ErrorInfo{ERR_CODE_REGISTER_FAILED, "手机号已被注册过"}
				} 
			}
			return nil
		}
	}
}

// ????应当严重手机号/邮箱/验证码等的验证

// 写表，生成各注册字段
func (s AuthService) insertUser(req RegisterRequest) (*ErrorInfo) {
	var user User
	if db, err := db.OpenDatabase(); err != nil {
		return &ErrorInfo{ERR_CODE_DB_EXCEPTION, "数据库异常"}
	} else {
		defer db.Close()
		// 密码做hash
		password := []byte(req.Pwd)
		if hashPwd, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost); err != nil {
			return &ErrorInfo{ERR_CODE_REQUEST_FAILED, err.Error()}
		} else {
			user.UID = uuid.New().String()
			user.Pwd = string(hashPwd)
			user.Name = req.Name
			user.Email = req.Email
			user.Admin = 0
			user.Phone = req.Phone
			user.LoginAt = time.Now().Unix()
			user.CreateAt = time.Now().Unix()
			if err := db.DB().Create(&user).Error; err != nil {
				return &ErrorInfo{ERR_CODE_REQUEST_FAILED, err.Error()}
			} else {
				return nil
			}
		}
	}
}
