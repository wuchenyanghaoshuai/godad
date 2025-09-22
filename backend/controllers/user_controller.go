package controllers

import (
    "net/http"
    "strconv"

    "godad-backend/middleware"
    "godad-backend/models"
    "godad-backend/services"
    "godad-backend/utils"

    "github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService    *services.UserService
	articleService *services.ArticleService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService:    services.NewUserService(),
		articleService: services.NewArticleService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body models.UserRegisterRequest true "注册信息"
// @Success 200 {object} utils.Response{data=models.UserResponse} "注册成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Router /api/user/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req models.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层注册用户
	user, err := c.userService.Register(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	// 返回用户信息（不包含密码）
	utils.SuccessWithMessage(ctx, "注册成功", user.ToResponse())
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=map[string]interface{}} "登录成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "认证失败"
// @Router /api/user/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req models.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层登录
	user, err := c.userService.Login(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeUnauthorized, err.Error())
		return
	}

    // 生成JWT令牌（access + refresh）
    token, err := middleware.GenerateToken(user)
    if err != nil {
        utils.Error(ctx, utils.CodeInternalError, "生成令牌失败")
        return
    }
    refresh, err := middleware.GenerateRefreshToken(user)
    if err != nil {
        utils.Error(ctx, utils.CodeInternalError, "生成刷新令牌失败")
        return
    }
    // 设置 httpOnly Cookie（开发环境 secure=false）
    // 并返回用户信息（为兼容前端旧逻辑，保留 token 字段）
    // 写入 Cookie
    // 通过中间件方法设置
    // 需要导出一个函数或在此处复制逻辑，这里直接调用未导出的函数不行，改为在此处设置 cookie
    // 由于 setAuthCookies 在 middleware 包内未导出，这里重复设置
    ctx.SetSameSite(http.SameSiteLaxMode)
    ctx.SetCookie("access_token", token, 0, "/", "", false, true)
    ctx.SetCookie("refresh_token", refresh, 0, "/", "", false, true)

	// 返回用户信息和令牌
    utils.SuccessWithMessage(ctx, "登录成功", gin.H{
        "user":  user.ToResponse(),
        "token": token,
    })
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=models.UserResponse} "获取成功"
// @Failure 401 {object} utils.Response "未授权"
// @Router /api/user/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "未授权")
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, err.Error())
		return
	}

	utils.Success(ctx, user.ToResponse())
}

// UpdateProfile 更新当前用户信息
// @Summary 更新当前用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.UserUpdateRequest true "更新信息"
// @Success 200 {object} utils.Response{data=models.UserResponse} "更新成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Router /api/user/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "未授权")
		return
	}

	var req models.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 更新用户信息
	user, err := c.userService.UpdateUser(userID, &req)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "更新成功", user.ToResponse())
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body map[string]string true "密码信息"
// @Success 200 {object} utils.Response "修改成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Router /api/user/change-password [post]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	// 从中间件获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "未授权")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 修改密码
	if err := c.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "密码修改成功", nil)
}

// CheckNickname 检查昵称是否可用
// @Summary 检查昵称是否可用
// @Description 检查指定昵称是否已被使用
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param nickname query string true "要检查的昵称"
// @Success 200 {object} utils.Response{data=map[string]interface{}} "检查成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Router /api/user/check-nickname [get]
func (c *UserController) CheckNickname(ctx *gin.Context) {
	nickname := ctx.Query("nickname")
	if nickname == "" {
		utils.Error(ctx, utils.CodeBadRequest, "昵称参数不能为空")
		return
	}

	// 检查昵称长度
	if len(nickname) < 2 || len(nickname) > 12 {
		utils.Error(ctx, utils.CodeBadRequest, "昵称长度必须在2-12个字符之间")
		return
	}

	// 检查昵称是否已存在
	exists, err := c.userService.CheckNicknameExists(nickname)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"nickname":  nickname,
		"exists":    exists,
		"available": !exists,
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 刷新用户的JWT令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=map[string]string} "刷新成功"
// @Failure 401 {object} utils.Response "未授权"
// @Router /api/user/refresh-token [post]
func (c *UserController) RefreshToken(ctx *gin.Context) {
	// 从中间件获取当前用户信息
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil {
		utils.Error(ctx, utils.CodeUnauthorized, "未授权")
		return
	}

	// 生成新的JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, "生成令牌失败")
		return
	}

	utils.SuccessWithMessage(ctx, "令牌刷新成功", gin.H{"token": token})
}

// GetUserByID 根据ID获取用户信息（公开信息）
// @Summary 根据ID获取用户信息
// @Description 获取指定用户的公开信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response{data=models.UserResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /api/user/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "用户ID格式错误")
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, err.Error())
		return
	}

	// 返回公开信息（隐藏敏感信息）
	response := user.ToResponse()
	response.Email = "" // 隐藏邮箱
	response.Phone = "" // 隐藏手机号

	utils.Success(ctx, response)
}

// GetUserByUsername 根据用户名获取用户信息（公开信息）
// @Summary 根据用户名获取用户信息
// @Description 获取指定用户的公开信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} utils.Response{data=models.UserResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /api/user/profile/{username} [get]
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		utils.Error(ctx, utils.CodeBadRequest, "用户名不能为空")
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByUsername(username)
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, err.Error())
		return
	}

	// 返回公开信息（隐藏敏感信息）
	response := user.ToResponse()
	response.Email = "" // 隐藏邮箱
	response.Phone = "" // 隐藏手机号

	utils.Success(ctx, response)
}

// GetUserList 获取用户列表（管理员功能）
// @Summary 获取用户列表
// @Description 获取用户列表，支持分页和搜索（管理员功能）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} utils.PageResponse{data=[]models.UserResponse} "获取成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 403 {object} utils.Response "权限不足"
// @Router /api/user/list [get]
func (c *UserController) GetUserList(ctx *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	keyword := ctx.Query("keyword")

	// 参数验证
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取用户列表
	users, total, err := c.userService.GetUserList(page, size, keyword)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var userResponses []*models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	utils.SuccessPage(ctx, userResponses, total, page, size)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口（客户端删除token即可）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response "登出成功"
// @Router /api/user/logout [post]
func (c *UserController) Logout(ctx *gin.Context) {
    // 清除 Cookie（access_token, refresh_token）
    ctx.SetSameSite(http.SameSiteLaxMode)
    ctx.SetCookie("access_token", "", -1, "/", "", false, true)
    ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)
    utils.SuccessWithMessage(ctx, "登出成功", nil)
}

// GenerateRandomNickname 生成随机昵称
// @Summary 生成随机昵称
// @Description 为用户生成一个随机的、可爱的昵称
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=map[string]string} "生成成功"
// @Failure 500 {object} utils.Response "生成失败"
// @Router /api/user/generate-nickname [post]
func (c *UserController) GenerateRandomNickname(ctx *gin.Context) {
	// 生成随机昵称
	nickname, err := c.userService.GenerateRandomNickname()
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, "生成昵称失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"nickname": nickname,
	})
}

// GetUserArticles 获取用户文章列表
// @Summary 获取用户文章列表
// @Description 获取指定用户的文章列表，支持分页
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} utils.PageResponse{data=[]models.ArticleResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /api/user/{id}/articles [get]
func (c *UserController) GetUserArticles(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "用户ID格式错误")
		return
	}

	// 验证用户是否存在
	_, err = c.userService.GetUserByID(uint(id))
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, err.Error())
		return
	}

	// 绑定查询参数
	var req models.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// 设置查询条件：只显示已发布的文章和该用户的文章
	req.AuthorID = uint(id)
	req.Status = 1 // 只显示已发布的文章

	// 获取文章列表
	articles, total, err := c.articleService.GetArticleList(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse(false))
	}

	utils.SuccessPage(ctx, responses, total, req.Page, req.Size)
}

// GetUserArticlesByUsername 根据用户名获取用户文章列表
// @Summary 根据用户名获取用户文章列表
// @Description 获取指定用户的文章列表，支持分页
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} utils.PageResponse{data=[]models.ArticleResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /api/user/profile/{username}/articles [get]
func (c *UserController) GetUserArticlesByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		utils.Error(ctx, utils.CodeBadRequest, "用户名不能为空")
		return
	}

	// 验证用户是否存在
	user, err := c.userService.GetUserByUsername(username)
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, err.Error())
		return
	}

	// 绑定查询参数
	var req models.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// 设置查询条件：只显示已发布的文章和该用户的文章
	req.AuthorID = user.ID
	req.Status = 1 // 只显示已发布的文章

	// 获取文章列表
	articles, total, err := c.articleService.GetArticleList(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse(false))
	}

	utils.SuccessPage(ctx, responses, total, req.Page, req.Size)
}
