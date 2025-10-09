//go:build openai
// +build openai

package main

import (
	"factory-go-api/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	// 在测试开始前加载配置
	_, err := config.LoadConfig("config.json")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "环境变量存在",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "环境变量不存在",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}
			if got := getEnv(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	// 重置 startTime 为当前时间以便测试
	startTime = time.Now()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}
}

func TestIsModelSupported(t *testing.T) {
	tests := []struct {
		name    string
		modelID string
		want    bool
	}{
		{
			name:    "支持的模型 - claude-opus-4-1",
			modelID: "claude-opus-4-1-20250805",
			want:    true,
		},
		{
			name:    "支持的模型 - claude-sonnet-4",
			modelID: "claude-sonnet-4-20250514",
			want:    true,
		},
		{
			name:    "支持的模型 - claude-sonnet-4-5",
			modelID: "claude-sonnet-4-5-20250929",
			want:    true,
		},
		{
			name:    "不支持的模型",
			modelID: "gpt-4",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.IsModelSupported(tt.modelID); got != tt.want {
				t.Errorf("config.IsModelSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "a < b",
			a:    1,
			b:    2,
			want: 1,
		},
		{
			name: "a > b",
			a:    5,
			b:    3,
			want: 3,
		},
		{
			name: "a == b",
			a:    4,
			b:    4,
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.a, tt.b); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseRecorder(t *testing.T) {
	rec := httptest.NewRecorder()
	recorder := &responseRecorder{
		ResponseWriter: rec,
		statusCode:     http.StatusOK,
	}

	// 测试 WriteHeader
	recorder.WriteHeader(http.StatusNotFound)
	if recorder.statusCode != http.StatusNotFound {
		t.Errorf("WriteHeader() statusCode = %v, want %v", recorder.statusCode, http.StatusNotFound)
	}
}