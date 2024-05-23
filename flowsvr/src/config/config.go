package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/niuniumart/gosdk/martlog"
)

var Conf *TomlConfig

var (
	TestFilePath string
)

// TomlConfig
type TomlConfig struct {
	Common commonConfig
	MySQL  mysqlConfig
	Redis  redisConfig
	Task   TaskConfig
}

type commonConfig struct {
	Port    int  `toml:"port"`
	OpenTLS bool `toml:"open_tls"`
}

type mysqlConfig struct {
	Url    string `toml:"url"`
	User   string `toml:"user"`
	Pwd    string `toml:"pwd"`
	Dbname string `toml:"db_name"`
}

type redisConfig struct {
	Url                    string `toml:"url"`
	Auth                   string `toml:"auth"`
	MaxIdle                int    `toml:"max_idle"`
	MaxActive              int    `toml:"max_active"`
	IdleTimeout            int    `toml:"idle_timeout"`
	CacheTimeout           int    `toml:"cache_timeout"`
	CacheTimeoutVerifyCode int    `toml:"cache_timeout_verify_code"`
	CacheTimeoutDay        int    `toml:"cache_timeout_day"`
}

type TaskConfig struct {
	TableMaxRows        int   `toml:"table_max_rows"`
	AliveThreshold      int   `toml:"alive_threshold"`
	SplitInterval       int   `toml:"split_interval"`
	LongProcessInterval int   `toml:"long_process_interval"`
	MoveInterval        int   `toml:"move_interval"`
	MaxProcessTime      int64 `toml:"max_process_time"`
}

// LoadConfig
func (c *TomlConfig) LoadConfig(env string) {
	if env == "" {
		env = "test"
	}

	filePath := "../config/config-" + env + ".toml"
	if TestFilePath != "" {
		filePath = TestFilePath
	}

	if _, err := os.Stat(filePath); err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		panic(err)
	}
}

const (
	USAGE = "Usage: asyncflow [-e <test|prod>]"
)

// GetConfEnv
func GetConfEnv() string {
	usage := "./main {$env} "

	env := os.Getenv("ENV")
	if env == "" {
		if len(os.Args) < 2 {
			fmt.Println("not enough params, usage:  ", usage)
			os.Exit(1)
		}
		if len(os.Args) >= 4 {
			env = "test"
		} else {
			env = os.Args[1]
		}
	}

	return env
}

func Init() {
	//
	env := GetConfEnv()
	InitConf(env)
}

// InitConf
func InitConf(env string) {
	Conf = new(TomlConfig)
	Conf.LoadConfig(env)
	printLog()
}

func printLog() {
	martlog.Infof("======== [Common] ========")
	martlog.Infof("%+v", Conf.Common)
	martlog.Infof("======== [MySQL] ========")
	martlog.Infof("%+v", Conf.MySQL)
	martlog.Infof("======== [Redis] ========")
	martlog.Infof("%+v", Conf.Redis)
}
