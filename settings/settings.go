package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var GlobalConfig = new(AppConfig)

type AppConfig struct {
	Name             string `mapstructure:"name"`
	Mode             string `mapstructure:"mode"`
	Version          string `mapstructure:"version"`
	Port             int    `mapstructure:"port"`
	*LogConfig       `mapstructure:"log"`
	*MysqlConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*SnowflakeConfig `mapstructure:"snowflake"`
}

type LogConfig struct {
	Level        string `mapstructure:"level"`
	Dir          string `mapstructure:"dir"`
	FileName     string `mapstructure:"filename"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxAge       int    `mapstructure:"max_age"`
	MaxBackups   int    `mapstructure:"max_backups"`
	LogInConsole bool   `mapstructure:"logInConsole"`
}

type MysqlConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	MaxOpenConnects int    `mapstructure:"max_open_connects"`
	MaxIdleConnects int    `mapstructure:"max_idle_connects"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Dbname   int    `mapstructure:"dbname"`
	PoolSize int    `mapstructure:"pool_size"`
}

type SnowflakeConfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_ID"`
}

// Viper
// 优先级: 命令行 > 环境变量 > 默认值

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err: %v\n", err)
		return
	}

	viper.WatchConfig()

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		fmt.Println("viper.Unmarshal failed,err:", err)
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件已修改...")

		if err := viper.Unmarshal(&GlobalConfig); err != nil {
			fmt.Println("viper.  failed,err:", err)
		}
	})

	return
}
