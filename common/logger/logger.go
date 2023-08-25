package logger

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
)

type LoggerConf struct {
	Filename string `json:"filename"`
	MaxSize  int64  `json:"maxsize"`
	MaxDays  int    `json:"maxdays"`
	Color    bool   `json:"color"`
	Perm     string `json:"perm"`
}

func (c LoggerConf) JSONString() string {
	if dst, err := json.Marshal(c); err == nil {
		return string(dst)
	}
	return "{}"
}

func Init(console bool, prefix, level string, conf LoggerConf) error {
	if console {
		logs.NewLogger()
		logs.SetLogger(logs.AdapterConsole)
	} else {
		logs.NewLogger(100)
		logs.SetLogger(logs.AdapterFile, conf.JSONString())
	}

	logs.SetPrefix(prefix)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(4)

	lv := logs.LevelInformational
	if len(level) > 0 {
		switch level[0] {
		case 'd':
			lv = logs.LevelDebug
		case 'w':
			lv = logs.LevelWarning
		case 'e':
			lv = logs.LevelError
		case 'f':
			lv = logs.LevelEmergency
		default:
			lv = logs.LevelInformational
		}
	}
	logs.SetLevel(lv)

	return nil
}

func GetLogger() *logs.BeeLogger {
	return logs.GetBeeLogger()
}

func Debugf(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
}

func Debugln(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
}

func Infof(f interface{}, v ...interface{}) {
	logs.Info(f, v...)
}

func Infoln(f interface{}, v ...interface{}) {
	logs.Info(f, v...)
}

func Warning(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)
}

func Warningf(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)
}

func Warningln(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)
}

func Errorln(f interface{}, v ...interface{}) {
	logs.Error(f, v...)
}

func Errorf(f interface{}, v ...interface{}) {
	logs.Error(f, v...)
}

func Panic(f interface{}, v ...interface{}) {
	Errorln(f, v...)
	panic("panic")
}

func Panicln(f interface{}, v ...interface{}) {
	Errorln(f, v...)
	panic("panic")
}

func Flush() {
	logs.GetBeeLogger().Flush()
}
