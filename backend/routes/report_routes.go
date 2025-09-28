package routes

import (
    "godad-backend/controllers"
    "godad-backend/middleware"
    "github.com/gin-gonic/gin"
)

func SetupReportRoutes(router *gin.Engine) {
    rc := controllers.NewReportController()

    v1 := router.Group("/api")
    // 用户侧（需登录）
    auth := v1.Group("")
    auth.Use(middleware.AuthMiddleware())
    {
        auth.POST("/reports", rc.CreateReport)
        auth.GET("/reports/my", rc.MyReports)
    }

    // 管理侧
    admin := v1.Group("/admin")
    admin.Use(middleware.AuthMiddleware())
    admin.Use(middleware.AdminMiddleware())
    {
        admin.GET("/reports", rc.AdminList)
        admin.PUT("/reports/:id/status", rc.AdminUpdateStatus)
    }
}

