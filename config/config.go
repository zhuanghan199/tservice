/*
*

	全局的配置中心

*
*/
package config

import (
	"tservice/models"
)

var Cnf models.Config

func init() {
	Cnf = models.Config{
		MysqlAddr:     "localhost:3306",
		MysqlPWD:      "",
		MysqlUserName: "root",
		MysqlDBName:   "tservice",
		RedisAddr:     "localhost:6379",

		JWTKey:			"tservice1314",
		JWTDuration:    3600,
		JWTRefreshDuration: 2592000,
	}
}
