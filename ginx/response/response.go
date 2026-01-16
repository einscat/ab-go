package response

import (
	"net/http"

	"github.com/einscat/ab-go/ginx/ecode"
	"github.com/gin-gonic/gin"
)

// Result 标准响应结构体
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// respond 内部基础响应方法
func respond(c *gin.Context, httpStatus int, res Result) {
	c.JSON(httpStatus, res)
}

// Success 成功响应 (带数据)
func Success(c *gin.Context, data any) {
	respond(c, http.StatusOK, Result{
		Code: ecode.Success.Code,
		Msg:  ecode.Success.Msg,
		Data: data,
	})
}

// SuccessMsg 成功响应 (自定义消息，无数据)
func SuccessMsg(c *gin.Context, msg string) {
	respond(c, http.StatusOK, Result{
		Code: ecode.Success.Code,
		Msg:  msg,
		Data: nil,
	})
}

// Fail 失败响应 (自动识别 error 类型)
func Fail(c *gin.Context, err error) {
	// 1. 默认兜底：系统未知错误 (HTTP 500)
	resp := Result{
		Code: ecode.ServerError.Code,
		Msg:  ecode.ServerError.Msg,
		Data: nil,
	}
	httpStatus := http.StatusInternalServerError

	// 2. 尝试断言为自定义业务错误 (*ecode.Error)
	if e, ok := err.(*ecode.Error); ok {
		resp.Code = e.Code
		resp.Msg = e.Msg
		resp.Data = e.Details // 自动带上 Details (如表单校验错误详情)

		// 3. 【关键】动态获取 HTTP 状态码
		// 依次走：注册表 -> 范围判断 -> 默认500
		httpStatus = ecode.GetStatus(e.Code)
	}

	// 发送响应
	respond(c, httpStatus, resp)
}
