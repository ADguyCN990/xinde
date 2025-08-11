#!/bin/bash

# --- 脚本说明 ---
# 功能: 生成 Swagger 文档并启动 Go 应用服务
# 环境: 需要在 Git Bash 或 WSL 中运行

# 设置脚本在遇到错误时立即退出
set -e

echo "==> 1. Generating Swagger documentation..."

# 运行 swag init 命令
# -g: 指定主程序文件
# --output: 指定文档输出目录
# --parseDependency: 解析依赖项
# --parseInternal: 解析 internal 包
# 我们从脚本所在目录返回到项目根目录 (cd ..) 来执行
(cd .. && swag init -g cmd/app/main.go --output docs --parseDependency --parseInternal)

echo "==> Swagger documentation generated successfully in 'docs' folder."
echo "--------------------------------------------------"
echo "==> 2. Starting Go application..."

# 运行 Go 应用
# 同样，返回到项目根目录来执行
(cd .. && go run cmd/app/main.go)