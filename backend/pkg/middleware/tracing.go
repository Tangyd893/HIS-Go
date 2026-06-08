// Package middleware 分布式链路追踪中间件
// 提供 TraceID 传播、Span 关联、日志注入能力
// 生产环境可引入 OpenTelemetry SDK 替换增强
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ctxKeyTraceID = "traceID"
	ctxKeySpanID  = "spanID"
	fmtMethodPath = "%s %s"
)

// TraceConfig 追踪配置
type TraceConfig struct {
	ServiceName string // 服务名称
	Enabled     bool   // 是否启用追踪
}

// DefaultTraceConfig 默认配置
var DefaultTraceConfig = TraceConfig{
	ServiceName: "his-go",
	Enabled:     true,
}

// 当前各服务的 Span 上下文（简化版，生产用 OpenTelemetry 替换）
var (
	traceMu     sync.RWMutex
	activeSpans = make(map[string]*Span)
)

// Span 表示一个追踪 Span（简化版）
type Span struct {
	TraceID    string
	SpanID     string
	ParentID   string
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Attributes map[string]string
	StatusCode int
	mu         sync.Mutex
}

// GenerateTraceID 生成 TraceID（32位十六进制）
func GenerateTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateSpanID 生成 SpanID（16位十六进制）
func GenerateSpanID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// StartSpan 创建新 Span
func StartSpan(name string, parentID string) *Span {
	span := &Span{
		TraceID:    GenerateTraceID(),
		SpanID:     GenerateSpanID(),
		ParentID:   parentID,
		Name:       name,
		StartTime:  time.Now(),
		Attributes: make(map[string]string),
	}

	if parentID != "" {
		traceMu.RLock()
		if parent, ok := activeSpans[parentID]; ok {
			span.TraceID = parent.TraceID
		}
		traceMu.RUnlock()
	}

	traceMu.Lock()
	activeSpans[span.SpanID] = span
	traceMu.Unlock()

	return span
}

// EndSpan 结束 Span 并记录状态码
func (s *Span) EndSpan(statusCode int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EndTime = time.Now()
	s.StatusCode = statusCode
}

// SetAttribute 设置属性
func (s *Span) SetAttribute(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Attributes[key] = value
}

// Duration Span 耗时
func (s *Span) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

const (
	HeaderTraceID      = "X-Trace-ID"
	HeaderSpanID       = "X-Span-ID"
	HeaderParentSpanID = "X-Parent-Span-ID"
	HeaderTraceParent  = "traceparent"
)

// Tracing 分布式链路追踪中间件
// 优先级: 请求头传入 TraceID > 自动生成
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头提取或生成 TraceID
		traceID := c.GetHeader(HeaderTraceID)
		if traceID == "" {
			traceID = GenerateTraceID()
		}

		parentSpanID := c.GetHeader(HeaderSpanID)
		spanID := GenerateSpanID()

		// 设置请求上下文
		c.Set(ctxKeyTraceID, traceID)
		c.Set(ctxKeySpanID, spanID)

		// 创建服务端 Span
		spanName := fmt.Sprintf(fmtMethodPath, c.Request.Method, c.FullPath())
		if c.FullPath() == "" {
			spanName = fmt.Sprintf(fmtMethodPath, c.Request.Method, c.Request.URL.Path)
		}

		span := StartSpan(spanName, parentSpanID)
		span.TraceID = traceID
		span.SetAttribute("http.method", c.Request.Method)
		span.SetAttribute("http.url", c.Request.URL.String())
		span.SetAttribute("client.ip", c.ClientIP())

		// 响应头中返回 TraceID，方便客户端关联
		c.Header(HeaderTraceID, traceID)
		c.Header(HeaderSpanID, spanID)

		c.Next()

		// 记录 HTTP 状态码
		span.EndSpan(c.Writer.Status())
	}
}

// TracingSpan 增强版追踪中间件（待接入 OpenTelemetry SDK 后启用）
// 当前提供 TraceID 传播 + 请求属性记录
func TracingSpan(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(HeaderTraceID)
		if traceID == "" {
			traceID = GenerateTraceID()
		}

		spanID := GenerateSpanID()
		parentSpanID := c.GetHeader(HeaderSpanID)

		c.Set(ctxKeyTraceID, traceID)
		c.Set(ctxKeySpanID, spanID)
		c.Set("serviceName", serviceName)

		spanName := fmt.Sprintf(fmtMethodPath, c.Request.Method, c.FullPath())
		if c.FullPath() == "" {
			spanName = fmt.Sprintf(fmtMethodPath, c.Request.Method, c.Request.URL.Path)
		}

		span := StartSpan(spanName, parentSpanID)
		span.TraceID = traceID
		span.SetAttribute("http.method", c.Request.Method)
		span.SetAttribute("http.url", c.Request.URL.String())
		span.SetAttribute("http.user_agent", c.Request.UserAgent())
		span.SetAttribute("client.ip", c.ClientIP())
		span.SetAttribute("service.name", serviceName)

		c.Header(HeaderTraceID, traceID)
		c.Header(HeaderSpanID, spanID)
		c.Header(HeaderParentSpanID, parentSpanID)

		c.Next()

		span.EndSpan(c.Writer.Status())
	}
}

// InjectTraceHeaders 将 Trace 信息注入到 HTTP 请求头（用于跨服务调用）
func InjectTraceHeaders(c *gin.Context, header http.Header) {
	if traceID, ok := c.Get(ctxKeyTraceID); ok {
		header.Set(HeaderTraceID, traceID.(string))
	}
	if spanID, ok := c.Get(ctxKeySpanID); ok {
		header.Set(HeaderSpanID, spanID.(string))
	}
}

// GetTraceInfo 从 Gin 上下文获取 Trace 信息（用于日志）
func GetTraceInfo(c *gin.Context) (traceID, spanID string) {
	if v, ok := c.Get(ctxKeyTraceID); ok {
		traceID = v.(string)
	}
	if v, ok := c.Get(ctxKeySpanID); ok {
		spanID = v.(string)
	}
	return
}
