package controllers

import (
    "strconv"
    "godad-backend/middleware"
    "godad-backend/services"
    "godad-backend/utils"
    "github.com/gin-gonic/gin"
)

type AppealController struct { svc *services.AppealService }

func NewAppealController() *AppealController { return &AppealController{ svc: services.NewAppealService() } }

// Create 用户创建申诉
func (c *AppealController) Create(ctx *gin.Context) {
    userID, ok := middleware.GetCurrentUserID(ctx)
    if !ok { utils.Error(ctx, utils.CodeUnauthorized, "请先登录"); return }
    var req services.CreateAppealRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.Error(ctx, utils.CodeBadRequest, "参数错误"); return }
    a, err := c.svc.CreateAppeal(userID, &req)
    if err != nil { utils.Error(ctx, utils.CodeBadRequest, err.Error()); return }
    utils.SuccessWithMessage(ctx, "申诉已提交", a)
}

// AdminList 管理员查看申诉列表
func (c *AppealController) AdminList(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
    status := ctx.Query("status")
    items, total, err := c.svc.ListAppeals(&services.AppealListParams{ Page: page, Size: size, Status: status })
    if err != nil { utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return }
    utils.SuccessWithMessage(ctx, "获取成功", utils.PagedResponse{ Items: items, Total: total, Page: page, Size: size, Pages: (total+int64(size)-1)/int64(size) })
}

// AdminUpdateStatus 管理员处理申诉
func (c *AppealController) AdminUpdateStatus(ctx *gin.Context) {
    id64, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil { utils.Error(ctx, utils.CodeBadRequest, "无效ID"); return }
    adminID, _ := middleware.GetCurrentUserID(ctx)
    var body struct { Status string `json:"status"`; HandledNote string `json:"handled_note"` }
    if err := ctx.ShouldBindJSON(&body); err != nil { utils.Error(ctx, utils.CodeBadRequest, "参数错误"); return }
    a, e := c.svc.UpdateStatus(uint(id64), adminID, body.Status, body.HandledNote)
    if e != nil { utils.Error(ctx, utils.CodeBadRequest, e.Error()); return }
    utils.Success(ctx, a)
}

