package redis

import (
	"bluebell_app/settings"
	"fmt"
	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		//Addr: fmt.Sprintf("%s:%d", viper.GetString(""+
		//	"redis.host"), viper.GetInt("redis.port")),
		//Password: viper.GetString("redis.password"), // no password set
		//DB:       viper.GetInt("redis.dbname"),      // use default DB
		//PoolSize: viper.GetInt("redis.pool_size"),

		// 反序列化
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.Dbname,   // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}
