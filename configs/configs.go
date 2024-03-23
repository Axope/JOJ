package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configs struct {
	Log   LogConfig
	Mysql MysqlConfig
	JWT   JWTConfig
	Mongo MongoConfig
}

var cfg Configs

func InitConfigs() {
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic("unmarshal error")
	}
}

func GetLogConfig() LogConfig {
	return cfg.Log
}
func GetMysqlConfig() MysqlConfig {
	return cfg.Mysql
}
func GetJWTConfig() JWTConfig {
	return cfg.JWT
}
func GetMongoConfig() MongoConfig {
	return cfg.Mongo
}
