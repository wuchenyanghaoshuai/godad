# GoDad 后端项目 Makefile

.PHONY: help run build test clean migrate migrate-status dev prod

# 默认目标
help:
	@echo "GoDad 后端项目管理工具"
	@echo ""
	@echo "可用命令:"
	@echo "  help          显示帮助信息"
	@echo "  dev           开发模式启动（自动迁移）"
	@echo "  prod          生产模式启动（不自动迁移）"
	@echo "  run           直接启动服务"
	@echo "  build         编译项目"
	@echo "  test          运行测试"
	@echo "  clean         清理编译文件"
	@echo "  migrate       运行数据库迁移"
	@echo "  migrate-status 查看数据库迁移状态"
	@echo ""
	@echo "环境变量:"
	@echo "  SERVER_ENV=development|production|testing"
	@echo "  DB_AUTO_MIGRATE=true|false"

# 开发模式启动
dev:
	@echo "🚀 启动开发模式..."
	@export SERVER_ENV=development && \
	export DB_AUTO_MIGRATE=true && \
	go run main.go

# 生产模式启动
prod:
	@echo "🚀 启动生产模式..."
	@export SERVER_ENV=production && \
	export DB_AUTO_MIGRATE=false && \
	go run main.go

# 直接启动服务
run:
	@echo "🚀 启动服务..."
	go run main.go

# 编译项目
build:
	@echo "🔨 编译项目..."
	go build -o bin/godad-backend main.go
	@echo "✅ 编译完成，可执行文件: bin/godad-backend"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test ./...

# 清理编译文件
clean:
	@echo "🧹 清理编译文件..."
	rm -rf bin/
	go clean

# 运行数据库迁移
migrate:
	@echo "🗄️  执行数据库迁移..."
	go run scripts/migrate.go --action=up

# 强制运行数据库迁移（生产环境）
migrate-force:
	@echo "🗄️  强制执行数据库迁移..."
	go run scripts/migrate.go --action=up --force

# 查看数据库迁移状态
migrate-status:
	@echo "🔍 查看数据库状态..."
	go run scripts/migrate.go --action=status

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	go mod tidy
	go mod download

# 代码格式化
fmt:
	@echo "✨ 格式化代码..."
	go fmt ./...

# 代码检查
vet:
	@echo "🔍 代码检查..."
	go vet ./...

# 创建 .env 文件
init-env:
	@if [ ! -f .env ]; then \
		echo "📄 创建环境配置文件..."; \
		cp .env.example .env; \
		echo "✅ 已创建 .env 文件，请编辑配置信息"; \
	else \
		echo "⚠️  .env 文件已存在"; \
	fi

# 完整的开发环境设置
setup: init-env deps migrate-status
	@echo "🎉 开发环境设置完成!"
	@echo "使用 'make dev' 启动开发服务器"