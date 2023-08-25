package db

import (
	"fmt"

	"database/sql"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"

	"tservice/common/logger"
	"tservice/config"
	"tservice/models"
)

var (
	logged     = false // 是否打印日志
	dbAddress  string  // 连接地址
	dbUsername string  // 用户名
	dbPassword string  // 密码
	dbName     string  // 数据库名称
)

// 配置连接信息
func InitMysql(isLogDB bool) {
	defer logger.Debugln("InitMysql ok")
	dbAddress = config.Cnf.MysqlAddr
	dbUsername = config.Cnf.MysqlUserName
	dbPassword = config.Cnf.MysqlPWD
	dbName = config.Cnf.MysqlDBName
	logged = isLogDB
	if db, err := OpenDatabase(); err != nil {
		logger.Errorln("open db error: ", err.Error())
	} else {
		if err := db.DB().AutoMigrate(&models.User{}).Error; err != nil {
			logger.Errorln("Auto migrate table error:", err)
		}
	}
}

type Database struct {
	name string // 数据库名
	db   *gorm.DB
}

func createDatabase() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", dbUsername, dbPassword, dbAddress))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " CHARACTER SET utf8mb4;")
	if err != nil {
		panic(err)
	}
}

func OpenDatabase() (*Database, error) {
	createDatabase()
	conn, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local", dbUsername, dbPassword, dbAddress, dbName),
	)
	if err != nil {
		logger.Errorln("Open db error: ", err.Error())
		return nil, err
	}
	conn.LogMode(logged)
	return &Database{dbName, conn}, nil
}

func (db *Database) Name() string {
	return db.name
}

func (db *Database) DB() *gorm.DB {
	return db.db
}

func (db *Database) Close() {
	if db.DB() != nil {
		db.DB().Close()
		db = nil
	}
}
