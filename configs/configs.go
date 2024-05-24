package configs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configs struct {
	Log      LogConfig
	Mysql    MysqlConfig
	JWT      JWTConfig
	Mongo    MongoConfig
	Rabbitmq RabbitmqConfig
	Datas    DatasConfig
	Redis    RedisConfig
}

var cfg Configs

// var tags map[string]string
var tags map[string]struct{}
var tagsList []string

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

	// tags config
	data, err := os.ReadFile("./configs/tags.json")
	if err != nil {
		panic(fmt.Sprintf("tags.json read error, err = %v", err))
	}
	err = json.Unmarshal(data, &tagsList)
	if err != nil {
		panic(fmt.Sprintf("decode tags json error, err = %v", err))
	}
	tags = make(map[string]struct{})
	for _, tag := range tagsList {
		tags[tag] = struct{}{}
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
func GetRBTConfig() RabbitmqConfig {
	return cfg.Rabbitmq
}
func GetDatasConfig() DatasConfig {
	return cfg.Datas
}
func GetRedisConfig() RedisConfig {
	return cfg.Redis
}

func GetTagColor(tag string) bool {
	_, ok := tags[tag]
	return ok
}
func GetTagsList() []string {
	return tagsList
}
