package main

import (
	"github.com/Axope/JOJ/configs"

	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/middleware/rabbitmq"
	"github.com/Axope/JOJ/internal/middleware/redis"
	"github.com/Axope/JOJ/internal/router"
	"github.com/Axope/JOJ/internal/service"

	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/log"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// load configs
	configs.InitConfigs()

	// log
	log.InitLogger()
	defer log.Logger.Sync()
	log.Logger.Info("Golbal log init success")
	// mysql
	if err := dao.InitMysql(); err != nil {
		log.Logger.Error("mysql init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("mysql init success")
	}
	// mongoDB
	if err := dao.InitMongo(); err != nil {
		log.Logger.Error("mongoDB init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("mongoDB init success")
	}
	// rabbitMQ
	recvMsgs, err := rabbitmq.InitMQ()
	if err != nil {
		log.Logger.Error("RabbitMQ init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("RabbitMQ init success")
		go service.SubmitService.HandleSubmitResult(recvMsgs)
	}
	// jwt
	if err := jwt.InitJWT(); err != nil {
		log.Logger.Error("jwt init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("jwt init success")
	}
	// redis
	if err := redis.InitRedis(); err != nil {
		log.Logger.Error("redis init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("redis init success")
	}

	router := router.NewRouter()
	router.Run(":9876")
}
