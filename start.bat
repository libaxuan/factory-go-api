@echo off
REM 设置 UTF-8 编码以支持中文显示
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM Factory Go API - 多模型支持启动脚本 (Windows)

echo [Factory Go API] - 多模型支持
echo ==================================

REM 加载 .env 文件（如果存在）
if exist .env (
    echo [INFO] 加载 .env 配置文件...
    for /f "usebackq tokens=1,* delims==" %%a in (.env) do (
        set "line=%%a"
        REM 跳过注释和空行
        if not "!line:~0,1!"=="#" if not "!line!"=="" (
            set "%%a=%%b"
        )
    )
    echo [OK] 环境变量已加载
) else (
    echo [WARN] 未找到 .env 文件，将使用环境变量或默认值
    echo    提示: 复制 .env.example 为 .env 并配置 FACTORY_API_KEY
)

REM 检查必需的环境变量
if "%FACTORY_API_KEY%"=="" (
    echo [ERROR] 错误: 未设置 FACTORY_API_KEY 环境变量
    echo    请在 .env 文件中设置或通过环境变量设置
    pause
    exit /b 1
)
echo [OK] FACTORY_API_KEY 已配置（源头 Key）

if not "%PROXY_API_KEY%"=="" (
    echo [OK] PROXY_API_KEY 已配置（对外代理 Key）
) else (
    echo [WARN] 未设置 PROXY_API_KEY，将使用直连模式
)

REM 检查 Go 是否安装
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Go 未安装，请先安装 Go: https://golang.org/dl/
    pause
    exit /b 1
)

for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo [OK] Go 版本: %GO_VERSION%

REM 检查并关闭占用 8003 端口的进程
echo [INFO] 检查端口 8003...
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8003 ^| findstr LISTENING') do (
    set PORT_PID=%%a
)

if defined PORT_PID (
    echo [WARN] 端口 8003 已被进程 !PORT_PID! 占用
    echo [INFO] 自动终止旧进程...
    taskkill /F /PID !PORT_PID! >nul 2>&1
    timeout /t 1 /nobreak >nul
    echo [OK] 旧进程已终止
) else (
    echo [OK] 端口 8003 可用
)

REM 安装依赖
echo [INFO] 安装依赖...
go mod tidy

REM 构建多模型版本
echo [INFO] 构建多模型支持版本...
go build -o factory-api.exe main_multimodel.go

if %errorlevel% equ 0 (
    echo [OK] 构建成功！
    echo.
    echo [服务信息]
    echo    地址: http://localhost:8003
    echo    文档: http://localhost:8003/docs
    echo    配置: config.json
    echo.
    echo [快速测试]
    echo    curl http://localhost:8003/health
    if not "%PROXY_API_KEY%"=="" (
        echo    curl http://localhost:8003/v1/models -H "Authorization: Bearer %PROXY_API_KEY%"
    ) else (
        echo    curl http://localhost:8003/v1/models -H "Authorization: Bearer %FACTORY_API_KEY%"
    )
    echo.
    echo [完整文档] type README.md
    echo ==================================
    echo.
    
    REM 启动服务
    factory-api.exe
) else (
    echo [ERROR] 构建失败
    pause
    exit /b 1
)