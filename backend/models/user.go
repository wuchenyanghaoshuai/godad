package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null;comment:密码"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(50);comment:昵称"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255);comment:头像URL"`
	Phone     string         `json:"phone" gorm:"type:varchar(20);comment:手机号"`
	Gender    int8           `json:"gender" gorm:"type:tinyint;default:0;comment:性别 0-未知 1-男 2-女"`
	Birthday  *time.Time     `json:"birthday" gorm:"comment:生日"`
	Bio       string         `json:"bio" gorm:"type:text;comment:个人简介"`
	Status    int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-正常"`
	Role      int8           `json:"role" gorm:"type:tinyint;default:1;comment:角色 1-普通用户 2-内容管理员 3-系统管理员"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Articles  []Article  `json:"articles,omitempty" gorm:"foreignKey:AuthorID"`
	Comments  []Comment  `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	Favorites []Favorite `json:"favorites,omitempty" gorm:"foreignKey:UserID"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"testuser"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=6,max=50" example:"123456"`
	Nickname string `json:"nickname" binding:"max=50" example:"测试用户"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Nickname string     `json:"nickname" binding:"max=50"`
	Avatar   string     `json:"avatar" binding:"max=255"`
	Phone    string     `json:"phone" binding:"max=20"`
	Gender   int8       `json:"gender" binding:"min=0,max=2"`
	Birthday *time.Time `json:"birthday"`
	Bio      string     `json:"bio" binding:"max=500"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Phone     string     `json:"phone"`
	Gender    int8       `json:"gender"`
	Birthday  *time.Time `json:"birthday"`
	Bio       string     `json:"bio"`
	Status    int8       `json:"status"`
	Role      int8       `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}