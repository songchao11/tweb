package global

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/labstack/gommon/log"
	"os"
)

type Config struct {
	DebugWeb  bool    //web调试模式
	DebugDB   bool    //db调试模式
	LogStdout bool    //日志是否输出到标准输出
	LogLevel  log.Lvl //日志级别
	LogPath   string  //日志路径
	DBType    string  //数据库类型
	DBAddr    string  //数据库地址
	Endpoint  string  //web服务监听地址
}

func LoadConfig() (conf *Config, err error) {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalf("load file config.ini failed: %s", err)
	}

	conf = &Config{}
	conf.LogStdout, err = cfg.Section("log").Key("stdout").Bool()
	if err != nil {
		return nil, fmt.Errorf("invalid log/stdout config")
	}
	switch cfg.Section("log").Key("level").String() {
	case "debug":
		conf.LogLevel = log.DEBUG
	case "info":
		conf.LogLevel = log.INFO
	case "warn":
		conf.LogLevel = log.WARN
	case "error":
		conf.LogLevel = log.ERROR
	default:
		log.Errorf("invalid log level!")
		return nil, fmt.Errorf("invalid log level!")
	}
	conf.LogPath = cfg.Section("log").Key("path").String()
	conf.DebugWeb, err = cfg.Section("debug").Key("web").Bool()
	if err != nil {
		return nil, err
	}
	conf.DebugDB, err = cfg.Section("debug").Key("db").Bool()
	if err != nil {
		return nil, err
	}
	conf.Endpoint = cfg.Section("http").Key("endpoint").String()
	conf.DBType = cfg.Section("db").Key("type").String()
	conf.DBAddr = cfg.Section("db").Key("addr").String()
	return
}

func (c *Config) Apply() bool {
	var err error
	var logOutput *os.File

	// 调整日志级别和输出
	log.SetLevel(c.LogLevel)
	if c.LogStdout {
		logOutput = os.Stdout
	} else {
		logOutput, err = os.OpenFile(c.LogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Errorf("initialize log output file %s failed: %s", c.LogPath, err)
			return false
		}
	}
	fp, ok := log.Output().(*os.File)
	if ok && fp != os.Stdout {
		fp.Close()
	}
	log.SetOutput(logOutput)
	return true
}
