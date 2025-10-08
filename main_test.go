//go:build !openai
// +build !openai

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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