@echo off

REM --- 脚本说明 ---
REM 功能: 生成 Swagger 文档并启动 Go 应用服务
REM 环境: 可以在 Windows CMD 或 PowerShell 中直接运行

echo ==^> 1. Generating Swagger documentation...

REM 切换到项目根目录
cd ..

REM 运行 swag init 命令
swag init -g cmd/app/main.go --output docs --parseDependency --parseInternal

REM 检查上一条命令是否成功
IF %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Swagger documentation generation failed. Aborting.
    goto :eof
)

echo.
echo ==^> Swagger documentation generated successfully in 'docs' folder.
echo --------------------------------------------------
echo ==^> 2. Starting Go application...
echo.

REM 运行 Go 应用
go run cmd/app/main.go

:eof