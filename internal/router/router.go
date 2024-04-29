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
	// router.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	corsMiddleware := cors.New(config)
	router.Use(corsMiddleware)

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
		// publicProblemGroup.GET("getProblem", v1.ProblemAPI.GetProblem)
		publicProblemGroup.GET("/:pid", v1.ProblemAPI.GetProblem)
	}
	publicSubmissionGroup := router.Group("/submission")
	{
		publicSubmissionGroup.GET("getSubmissionList", v1.SubmissionAPI.GetSubmissionList)
	}
	publicContestGroup := router.Group("/contest")
	{
		publicContestGroup.GET("getContestList", v1.ContestAPI.GetContestList)
		publicContestGroup.GET("/:cid", v1.ContestAPI.GetContest)
	}

	privateUserGroup := router.Group("/user")
	privateUserGroup.Use(mjwt.JWTAuth(false))
	{
		privateUserGroup.POST("/changePassword", v1.UserAPI.ChangePassword)
	}
	privateProblemGroup := router.Group("/problem")
	privateProblemGroup.Use(mjwt.JWTAuth(true))
	{
		privateProblemGroup.POST("/createProblem", v1.ProblemAPI.CreateProblem)
		// privateProblemGroup.PUT("/updateProblem", v1.ProblemAPI.UpdateProblem)
		// privateProblemGroup.DELETE("/deleteProblem", v1.ProblemAPI.DeleteProblem)
	}
	privateContestGroup := router.Group("/contest")
	privateContestGroup.Use(mjwt.JWTAuth(true))
	{
		privateContestGroup.POST("/createContest", v1.ContestAPI.CreateContest)
	}

	router.POST("/submit", mjwt.JWTAuth(false), v1.SubmitAPI.Submit)

	return router
}
