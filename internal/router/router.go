package router

import (
	"net/http"

	v1 "github.com/Axope/JOJ/api/v1"
	_ "github.com/Axope/JOJ/docs"
	mjwt "github.com/Axope/JOJ/internal/middleware/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	publicProblemGroup := router.Group("/problem")
	{
		publicProblemGroup.GET("getProblemList", v1.ProblemAPI.GetProblemList)
		publicProblemGroup.GET("getProblem", v1.ProblemAPI.GetProblem)
	}
	publicSubmissionGroup := router.Group("/submission")
	{
		publicSubmissionGroup.GET("getSubmissionList", v1.SubmissionAPI.GetSubmissionList)
	}

	privateUserGroup := router.Group("/user")
	privateUserGroup.Use(mjwt.JWTAuth(false))
	{
		privateUserGroup.POST("/changePassword", v1.UserAPI.ChangePassword)
	}
	// privateProblemGroup := router.Group("/problem")
	// privateProblemGroup.Use(middleware.JWTAuth(true))
	// {
	// 	privateProblemGroup.POST("/createProblem", v1.ProblemAPI.CreateProblem)
	// 	privateProblemGroup.PUT("/updateProblem", v1.ProblemAPI.UpdateProblem)
	// 	privateProblemGroup.DELETE("/deleteProblem", v1.ProblemAPI.DeleteProblem)
	// }

	router.POST("/submit", mjwt.JWTAuth(false), v1.SubmitAPI.Submit)

	return router
}
