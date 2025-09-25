package routes

import (
    "godad-backend/controllers"
    "godad-backend/middleware"

    "github.com/gin-gonic/gin"
)

func SetupMentionRoutes(router *gin.Engine) {
    ctrl := controllers.NewMentionController()
    g := router.Group("/api/mentions")
    g.Use(middleware.AuthMiddleware())
    {
        g.GET("/suggestions", ctrl.Suggest)
    }
}

