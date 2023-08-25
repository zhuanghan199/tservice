package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"tservice/common/logger"
	"tservice/config"
	"tservice/db"

	"github.com/spf13/viper"
)

var (
	_VERSION_ = "1.0.0"

	bin            = "zservice"
	flagLogDir     = "./"
	flagLogLevel   = "i"
	flagLogConsole = false
	flagLogDB      = false
	flagDebug      = false
	flagLogHttp    = false
)

func init() {
	initFlag()
	initConfig()
	initLog()
	initDB()
}

// 初始化flag
func initFlag() {
	flagHelp := false
	flagVer := false

	flag.BoolVar(&flagHelp, "h", false, "Print this message and exit")
	flag.BoolVar(&flagVer, "v", false, "Show version and exit")
	flag.BoolVar(&flagDebug, "d", false, "Set debug mode")
	flag.StringVar(&flagLogDir, "log_dir", "/tmp/log", "Log to dir")
	flag.BoolVar(&flagLogDB, "log_db", false, "Enable log db operation")
	flag.StringVar(&flagLogLevel, "l", "i", "Log level 'v/i/d/w/e/f'")
	flag.BoolVar(&flagLogConsole, "c", false, "Enable log to console")
	flag.BoolVar(&flagLogHttp, "log_http", false, "Enable log http request info")
	flag.Parse()
	if flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	if flagVer {
		str := "v" + _VERSION_
		fmt.Println(str)
		os.Exit(0)
	}
}

// 初始化yml，反序列进config中
func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}
	if err := viper.Unmarshal(&config.Cnf); err != nil {
		panic(fmt.Errorf("fatal Unmarshal config file: %v", err))
	}
}

// 初始化log
func initLog() {
	if flagDebug {
		flagLogLevel = "d"
	}
	fmt.Println(flagLogLevel)
	bin = filepath.Base(os.Args[0])
	logger.Init(flagLogConsole, "["+bin+"]", flagLogLevel,
		logger.LoggerConf{
			Filename: filepath.Join(flagLogDir, bin+".log"),
			MaxSize:  10 << 20, // 10MB
			MaxDays:  7,        // 7天
			Color:    false,
			Perm:     "0666",
		},
	)
}

// 初始化存储
func initDB() {
	defer logger.Debugln("initDB ok")
	db.InitMysql(true)
	db.InitRedis()
}
