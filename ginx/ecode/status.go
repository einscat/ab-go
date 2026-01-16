package ecode

import (
	"net/http"
	"sync"
)

// statusMap 存储特殊指定的映射关系: 业务码 -> HTTP状态码
var statusMap sync.Map

// RegisterStatus 手动注册业务码对应的 HTTP 状态码
// 场景：用于注册 401, 403, 404 等特殊状态
func RegisterStatus(code int, httpStatus int) {
	statusMap.Store(code, httpStatus)
}

// GetStatus 核心方法：获取业务码对应的 HTTP 状态码
func GetStatus(code int) int {
	// 1. 成功码 (0) -> HTTP 200
	if code == 0 {
		return http.StatusOK
	}

	// 2. Level 1: 优先查询手动注册表
	if status, ok := statusMap.Load(code); ok {
		return status.(int)
	}

	// 3. Level 2: 根据错误码范围自动推断 (约定优于配置)
	return statusFromCodeRange(code)
}

// statusFromCodeRange 根据错误码数值范围推断 HTTP 状态
func statusFromCodeRange(code int) int {
	// === 业务错误 (Business Error) ===
	// 约定：2 开头的错误码 (如 2001001) 表示业务逻辑错误
	// 这类错误通常返回 HTTP 200，由前端解析 JSON 中的 code 进行弹窗提示
	if code >= 2000000 {
		return http.StatusOK
	}

	// === 系统错误 (System Error) ===
	// 约定：1 开头的错误码 (如 10000001) 表示服务端内部错误
	// 这类错误返回 HTTP 500
	if code >= 1000000 {
		return http.StatusInternalServerError
	}

	// === 兜底 ===
	// 其他未知情况，默认当作服务器挂了
	return http.StatusInternalServerError
}
