package ecode

import "net/http"

var (
	Success       = New(0, "成功")
	ServerError   = New(10000000, "服务内部错误") // 1开头 -> 自动 500
	InvalidParams = New(10000001, "参数错误")   // 1开头 -> 自动 500
	// 下面这些需要特殊对待，所以需要手动注册
	NotFound     = New(10000002, "资源不存在")
	Unauthorized = New(10000003, "未授权")
	Forbidden    = New(10000004, "禁止访问")
	TooManyReq   = New(10000007, "请求过多")
)

func init() {
	// 注册通用错误的 HTTP 状态码
	RegisterStatus(Success.Code, http.StatusOK)
	RegisterStatus(NotFound.Code, http.StatusNotFound)          // 404
	RegisterStatus(Unauthorized.Code, http.StatusUnauthorized)  // 401
	RegisterStatus(Forbidden.Code, http.StatusForbidden)        // 403
	RegisterStatus(TooManyReq.Code, http.StatusTooManyRequests) // 429

	// InvalidParams 如果你希望它返回 400 Bad Request，也可以注册：
	RegisterStatus(InvalidParams.Code, http.StatusBadRequest)
}
