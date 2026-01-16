package ecode

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	e := New(100, "test error")
	assert.Equal(t, 100, e.Code)
	assert.Equal(t, "test error", e.Msg)
	assert.Equal(t, "code: 100, msg: test error", e.Error())
}

func TestErrorChaining(t *testing.T) {
	base := New(100, "base")

	// 测试不可变性 (WithDetails 不应修改 base)
	details := map[string]string{"field": "err"}
	e1 := base.WithDetails(details)

	assert.Nil(t, base.Details) // base 应该保持干净
	assert.Equal(t, details, e1.Details)

	// 测试 WithCause
	causeErr := errors.New("root cause")
	e2 := base.WithCause(causeErr)
	assert.Equal(t, causeErr, e2.Unwrap())
}

func TestGetStatus_Convention(t *testing.T) {
	// 1. 测试成功码
	assert.Equal(t, http.StatusOK, GetStatus(0))

	// 2. 测试 Level 2: 范围约定 (2xxxxxx -> 200)
	assert.Equal(t, http.StatusOK, GetStatus(2001001))
	assert.Equal(t, http.StatusOK, GetStatus(2999999))

	// 3. 测试 Level 2: 范围约定 (1xxxxxx -> 500)
	assert.Equal(t, http.StatusInternalServerError, GetStatus(1001001))

	// 4. 测试 Level 3: 兜底 (乱七八糟的码 -> 500)
	assert.Equal(t, http.StatusInternalServerError, GetStatus(999))
}

func TestGetStatus_Register(t *testing.T) {
	// 测试 Level 1: 手动注册优先级最高

	code := 2001004                               // 正常按约定应该是 200
	RegisterStatus(code, http.StatusUnauthorized) // 手动注册为 401

	assert.Equal(t, http.StatusUnauthorized, GetStatus(code))
}
