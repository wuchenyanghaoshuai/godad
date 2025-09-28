package routes

import (
    "godad-backend/controllers"
    "godad-backend/middleware"
    "github.com/gin-gonic/gin"
)

func SetupAppealRoutes(router *gin.Engine) {
    c := controllers.NewAppealController()
    v1 := router.Group("/api")
    {
        auth := v1.Group("")
        auth.Use(middleware.AuthMiddleware())
        auth.POST("/appeals", c.Create)

        admin := v1.Group("/admin")
        admin.Use(middleware.AuthMiddleware())
        admin.Use(middleware.AdminMiddleware())
        admin.GET("/appeals", c.AdminList)
        admin.PUT("/appeals/:id/status", c.AdminUpdateStatus)
    }
}

