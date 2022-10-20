package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// GlobalConfig 配置文件序列化结构体
var GlobalConfig = new(AppConfig)

//注意 当我们需要将viper读取的配置反序列到我们定义的结构体变量中时，一定要使用mapstructuretag！

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

// Init 优先级: 命令行 > 环境变量 > 默认值
func Init() (err error) {
	//指定配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	//读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err: %v\n", err)
		return
	}
	//监听配置文件变化
	viper.WatchConfig()

	//注意 当我们需要将viper读取的配置反序列到我们定义的结构体变量中时，一定要使用mapstructuretag！
	//读取配置文件反序列化
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		fmt.Println("viper.Unmarshal failed,err:", err)
	}
	//监听配置文件热加载
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件已修改...")
		if err := viper.Unmarshal(&GlobalConfig); err != nil {
			fmt.Println("viper.OnConfigChange.Unmarshal failed,err:", err)
		}
	})

	return
}
