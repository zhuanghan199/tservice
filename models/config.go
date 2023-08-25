package models

type Config struct {
	MysqlAddr     string
	MysqlUserName string
	MysqlPWD      string
	MysqlDBName   string

	RedisAddr string

	JWTKey             string
	JWTDuration        int
	JWTRefreshDuration int
}
