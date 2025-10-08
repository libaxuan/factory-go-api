#!/usr/bin/env python3
"""
测试 Factory Proxy 的 OpenAI 兼容接口
使用标准 OpenAI Python SDK
"""

from openai import OpenAI
import sys

def test_basic_chat():
    """测试基础对话"""
    print("🧪 测试 1: 基础对话")
    print("-" * 50)
    
    client = OpenAI(
        api_key="your-factory-key",
        base_url="http://localhost:8003/v1"
    )
    
    try:
        response = client.chat.completions.create(
            model="claude-sonnet-4-5-20250929",
            messages=[
                {"role": "user", "content": "你好，请用中文简短回复"}
            ],
            max_tokens=50
        )
        
        print(f"✅ 模型: {response.model}")
        print(f"✅ ID: {response.id}")
        print(f"✅ 回复: {response.choices[0].message.content}")
        print(f"✅ Token使用: {response.usage.total_tokens} (输入: {response.usage.prompt_tokens}, 输出: {response.usage.completion_tokens})")
        print()
        return True
    except Exception as e:
        print(f"❌ 错误: {e}")
        return False

def test_with_system():
    """测试带 system 消息"""
    print("🧪 测试 2: System 消息")
    print("-" * 50)
    
    client = OpenAI(
        api_key="your-factory-key",
        base_url="http://localhost:8003/v1"
    )
    
    try:
        response = client.chat.completions.create(
            model="claude-sonnet-4-5-20250929",
            messages=[
                {"role": "system", "content": "You are a helpful math tutor. Answer briefly."},
                {"role": "user", "content": "What is 7 * 8?"}
            ],
            max_tokens=30,
            temperature=0.7
        )
        
        print(f"✅ 回复: {response.choices[0].message.content}")
        print(f"✅ Finish reason: {response.choices[0].finish_reason}")
        print()
        return True
    except Exception as e:
        print(f"❌ 错误: {e}")
        return False

def test_multi_turn():
    """测试多轮对话"""
    print("🧪 测试 3: 多轮对话")
    print("-" * 50)
    
    client = OpenAI(
        api_key="your-factory-key",
        base_url="http://localhost:8003/v1"
    )
    
    try:
        response = client.chat.completions.create(
            model="claude-sonnet-4-5-20250929",
            messages=[
                {"role": "user", "content": "My name is Alice"},
                {"role": "assistant", "content": "Nice to meet you, Alice!"},
                {"role": "user", "content": "What's my name?"}
            ],
            max_tokens=20
        )
        
        content = response.choices[0].message.content
        print(f"✅ 回复: {content}")
        
        # 检查是否记住了名字
        if "Alice" in content or "alice" in content.lower():
            print("✅ 成功记住上下文！")
        else:
            print("⚠️  可能未正确处理上下文")
        print()
        return True
    except Exception as e:
        print(f"❌ 错误: {e}")
        return False

if __name__ == "__main__":
    print("=" * 50)
    print("Factory Proxy - OpenAI SDK 兼容性测试")
    print("=" * 50)
    print()
    
    results = []
    results.append(test_basic_chat())
    results.append(test_with_system())
    results.append(test_multi_turn())
    
    print("=" * 50)
    print(f"📊 测试结果: {sum(results)}/{len(results)} 通过")
    print("=" * 50)
    
    sys.exit(0 if all(results) else 1)