#!/usr/bin/env python3
"""
测试流式和非流式响应功能

依赖安装:
    pip install openai
"""

import os
import sys

try:
    from openai import OpenAI
except ImportError:
    print("❌ 错误: 未安装 openai 库")
    print("\n请运行以下命令安装依赖:")
    print("  pip install openai")
    print("\n或者:")
    print("  pip3 install openai")
    sys.exit(1)

# 配置
API_KEY = os.getenv("PROXY_API_KEY", "your_proxy_api_key")
BASE_URL = "http://localhost:8003/v1"
MODEL = "claude-sonnet-4-5-20250929"

client = OpenAI(
    api_key=API_KEY,
    base_url=BASE_URL
)

def test_non_stream():
    """测试非流式响应"""
    print("=" * 60)
    print("测试 1: 非流式响应")
    print("=" * 60)
    
    try:
        response = client.chat.completions.create(
            model=MODEL,
            messages=[
                {"role": "user", "content": "用一句话介绍自己"}
            ],
            max_tokens=100,
            stream=False
        )
        
        print(f"✅ 非流式响应成功")
        print(f"模型: {response.model}")
        print(f"内容: {response.choices[0].message.content}")
        print(f"Token 使用: {response.usage}")
        print()
        
    except Exception as e:
        print(f"❌ 非流式响应失败: {e}")
        print()

def test_stream():
    """测试流式响应"""
    print("=" * 60)
    print("测试 2: 流式响应")
    print("=" * 60)
    
    try:
        stream = client.chat.completions.create(
            model=MODEL,
            messages=[
                {"role": "user", "content": "用一句话介绍自己"}
            ],
            max_tokens=100,
            stream=True
        )
        
        print("✅ 流式响应成功，接收内容:")
        print("-" * 60)
        
        content = ""
        for chunk in stream:
            if chunk.choices[0].delta.content:
                text = chunk.choices[0].delta.content
                content += text
                print(text, end="", flush=True)
        
        print()
        print("-" * 60)
        print(f"完整内容: {content}")
        print()
        
    except Exception as e:
        print(f"❌ 流式响应失败: {e}")
        print()

def test_stream_with_system():
    """测试带系统提示的流式响应"""
    print("=" * 60)
    print("测试 3: 流式响应 + 系统提示")
    print("=" * 60)
    
    try:
        stream = client.chat.completions.create(
            model=MODEL,
            messages=[
                {"role": "system", "content": "你是一个专业的 Python 编程助手"},
                {"role": "user", "content": "写一个 Hello World 程序"}
            ],
            max_tokens=200,
            temperature=0.7,
            stream=True
        )
        
        print("✅ 流式响应成功，接收内容:")
        print("-" * 60)
        
        content = ""
        for chunk in stream:
            if chunk.choices[0].delta.content:
                text = chunk.choices[0].delta.content
                content += text
                print(text, end="", flush=True)
        
        print()
        print("-" * 60)
        print(f"完整内容长度: {len(content)} 字符")
        print()
        
    except Exception as e:
        print(f"❌ 流式响应失败: {e}")
        print()

if __name__ == "__main__":
    print("\n🚀 Factory Proxy API - 流式/非流式测试\n")
    
    # 测试非流式
    test_non_stream()
    
    # 测试流式
    test_stream()
    
    # 测试流式 + 系统提示
    test_stream_with_system()
    
    print("=" * 60)
    print("✅ 所有测试完成")
    print("=" * 60)