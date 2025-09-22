package services

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"godad-backend/config"
	"godad-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db           *gorm.DB
	cacheService *CacheService
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		db:           config.GetDB(),
		cacheService: NewCacheService(),
	}
}

// Register 用户注册
func (s *UserService) Register(req *models.UserRegisterRequest) (*models.User, error) {
	// 验证输入
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	var existUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", req.Email).First(&existUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 加密密码
	hashedPassword := s.hashPassword(req.Password)

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   1,
		Role:     1,
	}

	// 保存到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *models.UserLoginRequest) (*models.User, error) {
	// 验证输入
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	query := s.db.Where("status = ?", 1)
	if s.isEmail(req.Username) {
		query = query.Where("email = ?", req.Username)
	} else {
		query = query.Where("username = ?", req.Username)
	}

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或已被禁用")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 验证密码
	if !s.checkPassword(req.Password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 更新用户信息
	user.UpdatedAt = time.Now()

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	// 尝试从缓存获取
	if cachedUser, err := s.cacheService.GetUser(id); err == nil {
		return cachedUser, nil
	}

	var user models.User
	if err := s.db.Where("id = ? AND status = ?", id, 1).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 缓存用户信息
	s.cacheService.SetUser(id, &user, 30*time.Minute)

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}

	var user models.User
	if err := s.db.Where("username = ? AND status = ?", username, 1).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 缓存用户信息
	s.cacheService.SetUser(user.ID, &user, 30*time.Minute)

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *models.UserUpdateRequest) (*models.User, error) {
	// 获取用户
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 验证输入
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}



	// 检查手机号是否被其他用户使用
	if req.Phone != "" && req.Phone != user.Phone {
		var existUser models.User
		if err := s.db.Where("phone = ? AND id != ?", req.Phone, userID).First(&existUser).Error; err == nil {
			return nil, errors.New("手机号已被其他用户使用")
		}
	}

	// 更新字段
	updateData := make(map[string]interface{})
	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}
	updateData["gender"] = req.Gender
	if req.Birthday != nil {
		updateData["birthday"] = req.Birthday
	}
	if req.Bio != "" {
		updateData["bio"] = req.Bio
	}
	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}

	// 执行更新
	if len(updateData) > 0 {
		if err := s.db.Model(user).Updates(updateData).Error; err != nil {
			return nil, fmt.Errorf("更新用户信息失败: %v", err)
		}

		// 清理用户缓存
		s.cacheService.Delete(fmt.Sprintf("user:%d", userID))
	}

	// 重新查询用户信息
	return s.GetUserByID(userID)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.checkPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 验证新密码
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword := s.hashPassword(newPassword)

	// 更新密码
	if err := s.db.Model(user).Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("修改密码失败: %v", err)
	}

	// 清理用户缓存
	s.cacheService.Delete(fmt.Sprintf("user:%d", userID))

	return nil
}

// CheckNicknameExists 检查昵称是否已存在
func (s *UserService) CheckNicknameExists(nickname string) (bool, error) {
	if nickname == "" {
		return false, errors.New("昵称不能为空")
	}

	var count int64
	if err := s.db.Model(&models.User{}).Where("nickname = ? AND status = ?", nickname, 1).Count(&count).Error; err != nil {
		return false, fmt.Errorf("查询昵称失败: %v", err)
	}

	return count > 0, nil
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(page, size int, keyword string) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := s.db.Model(&models.User{}).Where("status = ?", 1)

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %v", err)
	}

	return users, total, nil
}

// validateRegisterRequest 验证注册请求
func (s *UserService) validateRegisterRequest(req *models.UserRegisterRequest) error {
	if req.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return errors.New("用户名长度必须在3-20个字符之间")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.Username) {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	if req.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if !s.isEmail(req.Email) {
		return errors.New("邮箱格式不正确")
	}

	if err := s.validatePassword(req.Password); err != nil {
		return err
	}

	if req.Nickname == "" {
		req.Nickname = req.Username // 默认使用用户名作为昵称
	}

	return nil
}

// validateLoginRequest 验证登录请求
func (s *UserService) validateLoginRequest(req *models.UserLoginRequest) error {
	if req.Username == "" {
		return errors.New("用户名/邮箱不能为空")
	}
	if req.Password == "" {
		return errors.New("密码不能为空")
	}
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *UserService) validateUpdateRequest(req *models.UserUpdateRequest) error {
	if req.Phone != "" && !s.isPhone(req.Phone) {
		return errors.New("手机号格式不正确")
	}
	return nil
}

// validatePassword 验证密码
func (s *UserService) validatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	if len(password) < 6 || len(password) > 20 {
		return errors.New("密码长度必须在6-20个字符之间")
	}
	return nil
}

// hashPassword 加密密码
// HashPassword 加密密码（公开方法）- 使用bcrypt
func (s *UserService) HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		// 如果加密失败，记录错误但返回一个无效的哈希
		fmt.Printf("密码加密失败: %v\n", err)
		return ""
	}
	return string(hashedPassword)
}

func (s *UserService) hashPassword(password string) string {
	return s.HashPassword(password)
}

// checkPassword 验证密码 - 使用bcrypt
func (s *UserService) checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// IsEmail 验证邮箱格式（公开方法）
func (s *UserService) IsEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isEmail 验证邮箱格式（兼容性保留）
func (s *UserService) isEmail(email string) bool {
	return s.IsEmail(email)
}

// isPhone 验证手机号格式
func (s *UserService) isPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// 随机昵称生成相关的数据
var (
	adjectives = []string{
		"温柔的", "可爱的", "聪明的", "快乐的", "阳光的", "善良的", "优雅的", "活泼的",
		"勇敢的", "耐心的", "细心的", "贴心的", "甜美的", "温暖的", "开朗的", "乐观的",
		"智慧的", "美丽的", "幸福的", "幸运的", "勤劳的", "健康的", "有爱的", "温馨的",
	}
	
	nouns = []string{
		"小熊", "小兔", "小猫", "小狗", "小鸟", "星星", "月亮", "太阳", "花朵", "小树",
		"宝贝", "天使", "公主", "王子", "小鱼", "蝴蝶", "彩虹", "云朵", "雪花", "露珠",
		"妈妈", "爸爸", "宝宝", "小精灵", "小仙女", "小王子", "小公主", "小天使", "小宝贝", "小可爱",
	}
)

// GenerateRandomNickname 生成随机昵称
func (s *UserService) GenerateRandomNickname() (string, error) {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())
	
	maxAttempts := 50 // 最大尝试次数，防止无限循环
	
	for i := 0; i < maxAttempts; i++ {
		// 随机选择形容词和名词
		adjective := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		
		// 生成昵称
		nickname := adjective + noun
		
		// 如果昵称太长，添加随机数字
		if len(nickname) > 10 {
			nickname = adjective[:3] + noun + fmt.Sprintf("%02d", rand.Intn(100))
		}
		
		// 检查昵称是否已存在
		exists, err := s.CheckNicknameExists(nickname)
		if err != nil {
			return "", err
		}
		
		if !exists {
			return nickname, nil
		}
		
		// 如果昵称已存在，尝试添加数字后缀
		for j := 1; j <= 99; j++ {
			nicknameWithNumber := fmt.Sprintf("%s%02d", nickname, j)
			exists, err := s.CheckNicknameExists(nicknameWithNumber)
			if err != nil {
				return "", err
			}
			if !exists {
				return nicknameWithNumber, nil
			}
		}
	}
	
	// 如果所有尝试都失败，生成一个基于时间戳的昵称
	timestamp := time.Now().Unix() % 10000
	fallbackNickname := fmt.Sprintf("用户%04d", timestamp)
	
	return fallbackNickname, nil
}

// GetUserCount 获取用户总数
func (s *UserService) GetUserCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.User{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取用户总数失败: %v", err)
	}
	return count, nil
}