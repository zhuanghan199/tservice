package models

type User struct {
	ID       int64  `gorm:"PRIMARY_KEY; AUTO_INCREMENT;"  json:"id"`
	UID      string `gorm:"column:uid; size:64; not null;" json:"uid"`
	Name     string `gorm:"column:name; size:256; not null;" json:"name"`
	Email    string `gorm:"column:email; size:256; not null;" json:"email"`
	Pwd      string `gorm:"column:pwd; size:256; not null;"`
	Phone    string `gorm:"column:phone; size:256;" json:"phone"`
	Avatar   string `gorm:"column:avatar; size:512;" json:"avatar"`
	Admin    int8   `gorm:"column:admin; default:0;" json:"admin"` // 0: 普通用户, 1: 管理员
	LoginAt  int64  `gorm:"column:login_at; default:0;" json:"login_at"`
	CreateAt int64  `gorm:"column:create_at; default:0;" json:"create_at"`
}

func (User) TableName() string {
	return "ts_users"
}
