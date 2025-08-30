# GoDad 一期功能测试计划

## 测试概述

本文档针对GoDad育儿知识分享平台一期功能制定全面的测试计划，确保系统的功能性、安全性、性能和用户体验达到预期标准。

## 🎯 测试目标

1. **功能完整性**: 确保所有一期功能按需求正常工作
2. **安全性**: 验证用户认证、权限控制、数据安全
3. **性能**: 确保系统在正常负载下稳定运行
4. **用户体验**: 验证界面友好性和交互流畅性
5. **兼容性**: 确保跨浏览器、设备兼容性

## 🧪 测试分类与策略

### 1. 单元测试 (Unit Testing)

#### 后端单元测试 (Go)

**测试框架**: 使用Go内置的`testing`包

**需要测试的模块**:

```bash
# 安装测试依赖
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
```

**测试范围**:
- **Service层测试**
  - `UserService`: 用户CRUD操作、密码验证
  - `ArticleService`: 文章管理、统计数据
  - `AuthService`: 登录验证、JWT生成
  - `CategoryService`: 分类管理
  - `CommentService`: 评论系统
  - `AdminService`: 管理员统计数据

- **Controller层测试**
  - API接口参数验证
  - 响应数据格式
  - 错误处理
  - 权限控制

**测试文件示例**:
```
backend/
├── tests/
│   ├── services/
│   │   ├── user_service_test.go
│   │   ├── auth_service_test.go
│   │   ├── article_service_test.go
│   │   └── admin_service_test.go
│   ├── controllers/
│   │   ├── user_controller_test.go
│   │   ├── auth_controller_test.go
│   │   └── admin_controller_test.go
│   └── utils/
│       ├── jwt_test.go
│       └── password_test.go
```

#### 前端单元测试 (Vue 3)

**测试框架**: Vitest + Vue Test Utils

```bash
# 安装前端测试依赖
npm install -D vitest @vue/test-utils jsdom
npm install -D @vitest/ui @vitest/coverage-c8
```

**测试范围**:
- **组件测试**
  - `RegisterPage`: 密码强度验证、表单提交
  - `LoginPage`: 登录流程、错误处理
  - `AdminDashboard`: 数据显示、权限检查
  - `Navbar`: 菜单显示、权限控制

- **Store测试** (Pinia)
  - `authStore`: 登录状态、权限检查
  - API调用逻辑
  - 状态更新机制

### 2. 集成测试 (Integration Testing)

#### API集成测试

**测试工具**: Postman + Newman 或 Go HTTP测试

**测试场景**:
```bash
# 用户注册流程
POST /api/auth/register → 成功注册
POST /api/auth/login → 登录获取token
GET /api/user/profile → 获取用户信息

# 管理员功能测试
POST /api/auth/login (admin) → 管理员登录
GET /api/admin/stats → 获取统计数据
POST /api/admin/categories → 创建分类

# 权限测试
GET /api/admin/stats (普通用户) → 403 Forbidden
```

#### 数据库集成测试

**测试内容**:
- 数据库连接和迁移
- CRUD操作完整性
- 事务处理正确性
- 数据一致性验证

### 3. 端到端测试 (E2E Testing)

**测试框架**: Playwright 或 Cypress

```bash
# 安装E2E测试工具
npm install -D playwright @playwright/test
```

**测试场景**:

1. **用户注册流程**
   - 填写注册表单
   - 密码强度验证提示
   - 注册成功跳转登录页
   - 登录成功进入首页

2. **管理员功能流程**
   - 管理员登录
   - 验证"后台管理"菜单显示
   - 进入后台查看统计数据
   - 分类管理操作

3. **权限控制测试**
   - 普通用户无法看到管理员菜单
   - 直接访问管理员URL被拦截
   - 未登录用户访问限制

## 📋 具体测试用例

### 1. 密码强度验证测试

| 测试用例 | 输入 | 期望结果 |
|----------|------|----------|
| 密码过短 | `123` | ❌ "密码长度至少需要8位" |
| 无大写字母 | `password123` | ❌ "密码必须包含至少一个大写字母" |
| 无小写字母 | `PASSWORD123` | ❌ "密码必须包含至少一个小写字母" |
| 符合要求 | `Password123` | ✅ 注册成功 |

### 2. 管理员权限测试

| 测试用例 | 用户类型 | 操作 | 期望结果 |
|----------|----------|------|----------|
| 普通用户 | role=0 | 访问 `/management-dashboard` | ❌ 404或重定向 |
| 普通用户 | role=0 | 调用 `/api/admin/stats` | ❌ 403 Forbidden |
| 管理员 | role=2 | 查看用户菜单 | ✅ 显示"后台管理" |
| 管理员 | role=2 | 访问后台页面 | ✅ 显示统计数据 |

### 3. 注册登录流程测试

| 步骤 | 操作 | 期望结果 |
|------|------|----------|
| 1 | 用户填写注册信息 | 表单验证通过 |
| 2 | 提交注册 | 跳转到登录页面 |
| 3 | 显示成功提示 | "注册成功！请使用您的账号密码登录" |
| 4 | 使用新账号登录 | 登录成功，进入首页 |
| 5 | 检查权限 | 普通用户看不到管理员菜单 |

## 🛠️ 测试环境配置

### 1. 测试数据库设置

```bash
# 创建测试数据库
CREATE DATABASE godad_test;

# 运行迁移
DB_NAME=godad_test go run scripts/migrate.go --action=up
```

### 2. 测试配置文件

创建 `.env.test`:
```env
DB_NAME=godad_test
DB_HOST=127.0.0.1
DB_PORT=3307
DB_USER=root
DB_PASSWORD=123456
JWT_SECRET=test-jwt-secret-key
SERVER_PORT=8889
```

### 3. 前端测试配置

`vitest.config.ts`:
```typescript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      reporter: ['text', 'json', 'html']
    }
  }
})
```

## 🚀 测试执行计划

### Phase 1: 基础功能测试 (Week 1)
- [ ] 后端单元测试 (Service层)
- [ ] 基础API接口测试
- [ ] 密码验证逻辑测试
- [ ] 用户认证流程测试

### Phase 2: 权限和安全测试 (Week 2)
- [ ] 管理员权限测试
- [ ] 路由守卫测试
- [ ] API权限控制测试
- [ ] 安全漏洞扫描

### Phase 3: 集成和E2E测试 (Week 3)
- [ ] 前后端集成测试
- [ ] 完整用户流程测试
- [ ] 浏览器兼容性测试
- [ ] 移动端适配测试

### Phase 4: 性能和压力测试 (Week 4)
- [ ] API性能测试
- [ ] 数据库查询性能
- [ ] 前端页面加载测试
- [ ] 并发用户测试

## 📊 测试报告和指标

### 测试覆盖率目标
- **后端代码覆盖率**: ≥ 80%
- **前端组件覆盖率**: ≥ 70%
- **API接口覆盖率**: 100%
- **核心业务流程**: 100%

### 性能指标
- **API响应时间**: < 200ms (95th percentile)
- **页面加载时间**: < 3秒
- **数据库查询**: < 100ms (平均)
- **并发用户**: 支持100+同时在线

## 🔧 测试工具和脚本

### 自动化测试脚本

```bash
#!/bin/bash
# test.sh - 一键执行所有测试

echo "🧪 开始执行GoDad测试套件..."

# 后端测试
echo "📡 执行后端单元测试..."
cd backend && go test ./... -v -cover

# 前端测试
echo "🖥️ 执行前端单元测试..."
cd ui && npm run test

# E2E测试
echo "🔄 执行端到端测试..."
cd ui && npm run test:e2e

# API测试
echo "🌐 执行API集成测试..."
newman run tests/postman/GoDad-API-Tests.json

echo "✅ 所有测试完成！"
```

### 持续集成配置

`.github/workflows/test.yml`:
```yaml
name: Test Suite
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: 123456
          MYSQL_DATABASE: godad_test
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.23
      - name: Run Backend Tests
        run: go test ./...
      - name: Setup Node
        uses: actions/setup-node@v3
        with:
          node-version: 18
      - name: Run Frontend Tests
        run: |
          cd ui
          npm install
          npm run test
```

## ✅ 验收标准

### 功能验收
- [ ] 所有核心功能正常运行
- [ ] 权限控制准确有效
- [ ] 错误处理友好合理
- [ ] 数据验证严格完整

### 质量验收
- [ ] 代码覆盖率达标
- [ ] 性能指标达标
- [ ] 安全测试通过
- [ ] 兼容性测试通过

### 用户体验验收
- [ ] 注册流程顺畅
- [ ] 登录体验良好
- [ ] 管理功能易用
- [ ] 错误提示清晰

---

**文档创建时间**: 2025年8月30日  
**负责团队**: GoDad开发团队  
**预估测试周期**: 4周  
**测试环境**: 开发环境 + 测试环境