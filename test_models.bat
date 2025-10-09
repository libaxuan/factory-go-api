@echo off
REM 设置 UTF-8 编码以支持中文显示
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM Factory Go API - 模型测试脚本 (Windows)

echo.
echo ========================================
echo   Factory Go API - 模型测试
echo ========================================
echo.

REM 加载 .env 文件
if exist .env (
    for /f "usebackq tokens=1,* delims==" %%a in (.env) do (
        set "line=%%a"
        if not "!line:~0,1!"=="#" if not "!line!"=="" (
            set "%%a=%%b"
        )
    )
)

REM 获取 API Key
if not "%PROXY_API_KEY%"=="" (
    set "API_KEY=%PROXY_API_KEY%"
    echo [INFO] 使用 PROXY_API_KEY
) else (
    set "API_KEY=%FACTORY_API_KEY%"
    echo [INFO] 使用 FACTORY_API_KEY
)

if "%API_KEY%"=="" (
    echo [ERROR] 错误: 未设置 API Key
    pause
    exit /b 1
)

REM 服务地址
set "BASE_URL=http://localhost:8003"

REM 测试提示词（使用简单的数学问题验证完整输出）
set "TEST_PROMPT=What is 123 + 456? Just give me the number."

REM Extended Thinking 模型的 max_tokens
set MAX_TOKENS_HIGH=30000
set MAX_TOKENS_LOW=100

REM 测试计数
set /a TOTAL=0
set /a SUCCESS=0
set /a FAILED=0

echo.
echo [INFO] 开始测试所有模型...
echo ========================================

REM 创建临时文件
set "TEMP_FILE=%TEMP%\factory_test_%RANDOM%.json"
set "RESPONSE_FILE=%TEMP%\factory_response_%RANDOM%.json"

echo.
echo [Claude 系列 - Anthropic]
echo ----------------------------------------

REM 测试 Claude Opus 4.1 (流式和非流式)
call :test_model "claude-opus-4-1-20250805" "Claude Opus 4.1 (非流式)" "false" %MAX_TOKENS_LOW%
call :test_model "claude-opus-4-1-20250805" "Claude Opus 4.1 (流式)" "true" %MAX_TOKENS_LOW%

REM 测试 Claude Sonnet 4 (流式和非流式)
call :test_model "claude-sonnet-4-20250514" "Claude Sonnet 4 (非流式)" "false" %MAX_TOKENS_HIGH%
call :test_model "claude-sonnet-4-20250514" "Claude Sonnet 4 (流式)" "true" %MAX_TOKENS_HIGH%

REM 测试 Claude Sonnet 4.5 (流式和非流式)
call :test_model "claude-sonnet-4-5-20250929" "Claude Sonnet 4.5 (非流式)" "false" %MAX_TOKENS_HIGH%
call :test_model "claude-sonnet-4-5-20250929" "Claude Sonnet 4.5 (流式)" "true" %MAX_TOKENS_HIGH%

echo.
echo [GPT 系列 - OpenAI]
echo ----------------------------------------

REM 测试 GPT-5 (流式和非流式)
call :test_model "gpt-5-2025-08-07" "GPT-5 (非流式)" "false" %MAX_TOKENS_HIGH%
call :test_model "gpt-5-2025-08-07" "GPT-5 (流式)" "true" %MAX_TOKENS_HIGH%

REM 测试 GPT-5 Codex (流式和非流式)
call :test_model "gpt-5-codex" "GPT-5 Codex (非流式)" "false" %MAX_TOKENS_LOW%
call :test_model "gpt-5-codex" "GPT-5 Codex (流式)" "true" %MAX_TOKENS_LOW%

REM 清理临时文件
if exist "%TEMP_FILE%" del "%TEMP_FILE%"
if exist "%RESPONSE_FILE%" del "%RESPONSE_FILE%"

REM 汇总结果
echo.
echo ========================================
echo   测试结果汇总
echo ========================================
echo 总测试数: %TOTAL%
echo 成功: %SUCCESS%
echo 失败: %FAILED%
set /a SUCCESS_RATE=%SUCCESS%*100/%TOTAL%
echo 成功率: %SUCCESS_RATE%%%
echo.

if %FAILED% equ 0 (
    echo [OK] 所有模型测试通过！
    pause
    exit /b 0
) else (
    echo [WARN] 部分模型测试失败，请检查日志
    pause
    exit /b 1
)

:test_model
set /a TOTAL+=1
set "MODEL_ID=%~1"
set "MODEL_NAME=%~2"
set "USE_STREAM=%~3"
set "MAX_TOKENS=%~4"

echo.
echo [%TOTAL%] 测试: %MODEL_NAME%
echo     模型: %MODEL_ID%
echo     流式: %USE_STREAM% ^| max_tokens: %MAX_TOKENS%

REM 构建请求 JSON
echo {"model": "%MODEL_ID%", "messages": [{"role": "user", "content": "%TEST_PROMPT%"}], "stream": %USE_STREAM%, "max_tokens": %MAX_TOKENS%} > "%TEMP_FILE%"

REM 发送请求
curl -s -X POST "%BASE_URL%/v1/chat/completions" ^
    -H "Content-Type: application/json" ^
    -H "Authorization: Bearer %API_KEY%" ^
    --max-time 30 ^
    -d @"%TEMP_FILE%" ^
    -o "%RESPONSE_FILE%" 2>nul

if %errorlevel% neq 0 (
    echo     状态: [FAIL] 请求失败
    set /a FAILED+=1
    goto :eof
)

REM 检查响应
if "%USE_STREAM%"=="true" (
    REM 流式响应检查
    findstr /C:"data:" "%RESPONSE_FILE%" >nul 2>&1
    if %errorlevel% equ 0 (
        echo     状态: [OK] 成功（流式）
        REM 显示响应的前100个字符
        for /f "usebackq delims=" %%i in (`type "%RESPONSE_FILE%"`) do (
            set "response=%%i"
            echo     响应片段: !response:~0,80!
            goto :stream_done
        )
        :stream_done
        set /a SUCCESS+=1
    ) else (
        findstr /C:"\"error\"" "%RESPONSE_FILE%" >nul 2>&1
        if %errorlevel% equ 0 (
            echo     状态: [FAIL] 失败
            type "%RESPONSE_FILE%"
            set /a FAILED+=1
        ) else (
            echo     状态: [WARN] 无数据
            set /a FAILED+=1
        )
    )
) else (
    REM 非流式响应检查
    findstr /C:"\"choices\"" "%RESPONSE_FILE%" >nul 2>&1
    if %errorlevel% equ 0 (
        REM 检查是否有content
        findstr /C:"\"content\"" "%RESPONSE_FILE%" >nul 2>&1
        if %errorlevel% equ 0 (
            echo     状态: [OK] 成功
            REM 显示响应的前100个字符
            for /f "usebackq delims=" %%i in (`type "%RESPONSE_FILE%"`) do (
                set "response=%%i"
                echo     响应: !response:~0,100!
                goto :done
            )
            :done
            set /a SUCCESS+=1
        ) else (
            echo     状态: [WARN] 响应为空
            set /a FAILED+=1
        )
    ) else (
        findstr /C:"\"error\"" "%RESPONSE_FILE%" >nul 2>&1
        if %errorlevel% equ 0 (
            echo     状态: [FAIL] 失败
            type "%RESPONSE_FILE%"
            set /a FAILED+=1
        ) else (
            echo     状态: [FAIL] 超时或无响应
            set /a FAILED+=1
        )
    )
)

goto :eof