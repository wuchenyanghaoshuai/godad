# GoDad 🚀

<div align="center">
  <img src="assets/logo-design.html" alt="GoDad Logo" width="200"/>

  **专为男性育儿知识分享打造的现代化社区平台**

  [![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
  [![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
  [![TypeScript](https://img.shields.io/badge/TypeScript-5.3+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
</div>

## 📖 项目简介

GoDad（出发爸）是一个专注于男性育儿知识分享的现代化社区平台。我们相信每一位父亲都应该积极参与孩子的成长过程，通过知识分享、经验交流和互助支持，让每一个家庭都能更好地"一起出发，共同成长"。

### 🌟 核心特色

- **💬 实时交流** - 集成聊天系统，支持私信、表情、图片分享
- **📝 知识分享** - 支持Markdown的富文本编辑器，文章分类管理
- **🔍 智能搜索** - 全文搜索，关键词高亮，个性化推荐
- **👥 社交互动** - 点赞、评论、关注、收藏完整社交体系
- **🔔 消息通知** - 实时通知系统，不错过任何重要互动
- **📱 响应式设计** - 完美适配桌面端和移动端

## 🛠 技术栈

### 后端服务
- **语言**: Go 1.23+
- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0+ + Redis
- **ORM**: GORM v1.30+
- **认证**: JWT (golang-jwt/jwt/v5)
- **存储**: 阿里云 OSS
- **邮件**: GoMail v2

### 前端应用
- **框架**: Vue 3.4 + TypeScript 5.3+
- **构建工具**: Vite 5.0+
- **样式**: Tailwind CSS 3.4+
- **状态管理**: Pinia 3.0+
- **路由**: Vue Router 4.2+
- **图标**: Heroicons, Lucide


## 🚀 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 1. 克隆项目

```bash
git clone https://github.com/your-username/godad.git
cd godad
```

### 2. 后端设置

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库和其他服务

# 运行后端服务
go run main.go
```

### 3. 前端设置

```bash
# 进入前端目录
cd ui

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 4. 数据库初始化

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE godad CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 运行迁移（首次启动后端时自动执行）
```

## 📋 功能特性

### 👤 用户系统
- [x] 用户注册/登录 (JWT认证)
- [x] 个人资料管理
- [x] 头像上传
- [x] 密码重置
- [x] 角色权限管理

### 📝 内容管理
- [x] Markdown文章编写
- [x] 文章分类标签
- [x] 草稿自动保存
- [x] 文章搜索
- [x] 内容审核

### 💬 社交功能
- [x] 文章点赞评论
- [x] 用户关注系统
- [x] 私信聊天
- [x] 实时通知
- [x] 收藏书签

### 🔍 搜索发现
- [x] 全文搜索
- [x] 分类浏览
- [x] 热门推荐
- [x] 搜索历史

### 🛡 管理后台
- [x] 用户管理
- [x] 内容审核
- [x] 分类管理
- [x] 系统统计

### 🤖 智能功能
- [x] 推荐算法
- [x] 消息限流
- [x] 缓存优化

## 📁 项目结构

```
godad/
├── backend/              # Go后端服务
│   ├── config/          # 配置管理
│   ├── controllers/     # API控制器
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   ├── middleware/      # 中间件
│   ├── routes/          # 路由配置
│   └── utils/           # 工具函数
├── ui/                  # Vue前端应用
│   ├── src/
│   │   ├── api/         # API客户端
│   │   ├── components/  # Vue组件
│   │   ├── pages/       # 页面组件
│   │   ├── stores/      # Pinia状态
│   │   ├── router/      # 路由配置
│   │   └── composables/ # 组合式函数
├── supabase/           # 数据库迁移
└── assets/             # 静态资源
```

## 🔧 开发指南

### API 文档

后端服务提供 RESTful API，主要端点包括：

- `POST /api/auth/login` - 用户登录
- `GET /api/articles` - 获取文章列表
- `POST /api/articles` - 创建文章
- `GET /api/users/profile` - 用户资料
- `POST /api/chat/send` - 发送消息

### 环境配置

后端 `.env` 文件示例：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=godad

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT密钥
JWT_SECRET=your_jwt_secret

# 阿里云OSS
OSS_ENDPOINT=your_oss_endpoint
OSS_ACCESS_KEY=your_access_key
OSS_SECRET_KEY=your_secret_key
OSS_BUCKET=your_bucket
```

### 代码规范

- **Go**: 遵循 Go 官方代码规范
- **Vue**: 使用 Composition API，TypeScript 严格模式
- **CSS**: 使用 Tailwind CSS 类名
- **提交**: 遵循 Conventional Commits 规范

## 🤝 贡献指南

我们欢迎所有形式的贡献！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发流程

1. 确保所有测试通过
2. 更新相关文档
3. 遵循代码规范
4. 提供清晰的提交信息

## 📊 路线图

### v1.0 (当前开发中)
- [x] 基础用户系统
- [x] 文章发布系统
- [x] 聊天功能
- [x] 搜索功能
- [ ] 移动端优化
- [ ] 性能优化

### v1.1 (计划中)
- [ ] 直播功能
- [ ] 小组讨论
- [ ] 积分系统
- [ ] 推送通知
- [ ] API 开放平台

### v2.0 (未来计划)
- [ ] AI 智能推荐
- [ ] 视频内容支持
- [ ] 多语言支持
- [ ] 移动应用

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有为 GoDad 项目做出贡献的开发者和用户！

特别感谢：
- [Gin](https://gin-gonic.com/) - 优秀的 Go Web 框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Tailwind CSS](https://tailwindcss.com/) - 实用的 CSS 框架

## 📞 联系我们

- **项目主页**: [https://github.com/your-username/godad](https://github.com/your-username/godad)
- **问题反馈**: [Issues](https://github.com/your-username/godad/issues)
- **讨论交流**: [Discussions](https://github.com/your-username/godad/discussions)

---

<div align="center">
  <strong>让我们一起出发，共同成长！ 🚀</strong>
</div>