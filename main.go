package main

import (
	"github.com/Axope/JOJ/configs"

	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/router"

	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/log"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	configs.InitConfigs()

	log.InitLogger()
	defer log.Logger.Sync()
	log.Logger.Info("Golbal log init success")
	if err := dao.InitMysql(); err != nil {
		log.Logger.Error("mysql init failed")
		return
	} else {
		log.Logger.Info("mysql init success")
	}

	jwt.InitJWT()

	router := router.NewRouter()
	router.Run(":9876")
}
