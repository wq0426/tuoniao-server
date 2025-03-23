package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var ConfigInstance *viper.Viper
var Rdb *redis.Client

type CustomConfig struct {
	*viper.Viper
	Id int
}

func NewConfig() *viper.Viper {
	if ConfigInstance == nil {
		envConf := os.Getenv("APP_CONF")
		if len(envConf) == 0 {
			flag.StringVar(&envConf, "conf", "config/local.yml", "config path, eg: -conf config/local.yml")
			flag.Parse()
		}
		fmt.Println("load conf file:", envConf)
		ConfigInstance = getConfig(envConf)
	}
	if Rdb == nil {
		Rdb = getRdb()
	}
	return ConfigInstance
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

func getRdb() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     ConfigInstance.GetString("data.db.redis.addr"),
		Password: ConfigInstance.GetString("data.db.redis.password"), // no password set
		DB:       ConfigInstance.GetInt("data.db.redis.db"),          // use default DB
	})
}
