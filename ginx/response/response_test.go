package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/einscat/ab-go/ginx/ecode"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter 辅助函数
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSuccess(t *testing.T) {
	r := setupRouter()
	r.GET("/ok", func(c *gin.Context) {
		Success(c, map[string]string{"foo": "bar"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Result
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "成功", resp.Msg) // 假设 ecode.Success.Msg 是 "成功"

	data := resp.Data.(map[string]interface{})
	assert.Equal(t, "bar", data["foo"])
}

func TestFail_BusinessError(t *testing.T) {
	// 测试：2开头业务错误 -> HTTP 200
	r := setupRouter()
	bizErr := ecode.New(2001001, "业务错误")

	r.GET("/biz", func(c *gin.Context) {
		Fail(c, bizErr)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/biz", nil)
	r.ServeHTTP(w, req)

	// 验证 HTTP 状态码是否为 200 (约定优于配置)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Result
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 2001001, resp.Code)
	assert.Equal(t, "业务错误", resp.Msg)
}

func TestFail_SystemError(t *testing.T) {
	// 测试：普通 error -> HTTP 500 + ServerError Code
	r := setupRouter()

	r.GET("/sys", func(c *gin.Context) {
		Fail(c, errors.New("db panic"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sys", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp Result
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, ecode.ServerError.Code, resp.Code)
	assert.Nil(t, resp.Data)
}

func TestFail_WithDetails(t *testing.T) {
	// 测试：带 Details 的错误 -> data 字段有值
	r := setupRouter()

	r.GET("/details", func(c *gin.Context) {
		err := ecode.InvalidParams.WithDetails(map[string]string{"k": "v"})
		Fail(c, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/details", nil)
	r.ServeHTTP(w, req)

	// InvalidParams 默认注册了 400 (在 ecode/common.go 中)
	// 如果你之前代码注册了 RegisterStatus(InvalidParams.Code, http.StatusBadRequest)
	// 这里应该是 400，否则可能是 500。假设已注册：
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Result
	json.Unmarshal(w.Body.Bytes(), &resp)

	details := resp.Data.(map[string]interface{})
	assert.Equal(t, "v", details["k"])
}
