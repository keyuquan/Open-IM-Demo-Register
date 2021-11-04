package main

import (
	"Company-Chat-Register/cmd/register"
	"Company-Chat-Register/common/config"
	"Company-Chat-Register/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(utils.CorsHandler())

	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/code", register.SendVerificationCode)
		authRouterGroup.POST("/verify", register.Verify)
		authRouterGroup.POST("/password", register.SetPassword)
		authRouterGroup.POST("/login", register.Login)
	}
	r.Run(":" + config.Config.Api.Port)
}
