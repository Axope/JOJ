package router

import (
	v1 "github.com/Axope/JOJ/api/v1"
	_ "github.com/Axope/JOJ/docs"
	"github.com/Axope/JOJ/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// GetIndex
//
//	@Router	/index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome!",
	})
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/index", GetIndex)

	publicUserGroup := router.Group("/user")
	{
		publicUserGroup.POST("/register", v1.UserAPI.Register)
		publicUserGroup.POST("/login", v1.UserAPI.Login)
	}

	privateUserGroup := router.Group("/user")
	privateUserGroup.Use(middleware.JWTAuth())
	{
		privateUserGroup.POST("/changePassword", v1.UserAPI.ChangePassword)
	}

	return router
}
