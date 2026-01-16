package validator_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/einscat/ab-go/ginx/validator"
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert" // 建议使用 testify 库进行断言
	// 引入我们要测试的包 (请替换为实际 module 路径)
)

// 模拟的请求结构体
type TestUser struct {
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"gte=18"`
	Email string `json:"email" binding:"email"`
	Phone string `json:"phone" binding:"mobile"` // 测试自定义规则
}

// SetupTestRouter 初始化一个用于测试的 Gin 路由
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 注册一个用于测试的路由
	r.POST("/test", func(c *gin.Context) {
		var req TestUser
		if valid, errs := validator.BindAndValid(c, &req); !valid {
			c.JSON(http.StatusBadRequest, errs)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

// 注册自定义规则 (模拟业务方注册)
func init() {
	validator.RegisterRule("mobile", func(fl val.FieldLevel) bool {
		ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, fl.Field().String())
		return ok
	}, "手机号格式错误")
}

func TestBindAndValid(t *testing.T) {
	router := SetupTestRouter()

	tests := []struct {
		name         string
		input        interface{}
		expectedCode int
		// 期望返回的错误包含的 Key (用于验证 map 是否正确)
		expectedErrKey string
	}{
		{
			name: "Valid Input",
			input: TestUser{
				Name:  "Tom",
				Age:   20,
				Email: "tom@example.com",
				Phone: "13800138000",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing Required Field",
			input: TestUser{
				// Name is missing
				Age: 18,
			},
			expectedCode:   http.StatusBadRequest,
			expectedErrKey: "name", // 期望 name 字段报错
		},
		{
			name: "Validation Failed (Age)",
			input: TestUser{
				Name: "Tom",
				Age:  10, // < 18
			},
			expectedCode:   http.StatusBadRequest,
			expectedErrKey: "age",
		},
		{
			name: "Custom Rule Failed (Mobile)",
			input: TestUser{
				Name:  "Tom",
				Age:   20,
				Phone: "110", // Invalid mobile
			},
			expectedCode:   http.StatusBadRequest,
			expectedErrKey: "phone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构造 JSON 请求
			jsonBytes, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 断言状态码
			assert.Equal(t, tt.expectedCode, w.Code)

			// 如果期望失败，检查返回的 Map 中是否包含对应的 Key
			if tt.expectedCode == http.StatusBadRequest {
				var respMap map[string]string
				_ = json.Unmarshal(w.Body.Bytes(), &respMap)
				_, exists := respMap[tt.expectedErrKey]
				assert.True(t, exists, "Expected error key '%s' not found in response", tt.expectedErrKey)
			}
		})
	}
}

// TestTypeError 测试类型错误 (例如 string 传 int)
func TestTypeError(t *testing.T) {
	router := SetupTestRouter()

	// 故意构造错误的 JSON: Age 应该是 int，这里传了 string
	jsonStr := `{"name": "Tom", "age": "abc"}`

	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var respMap map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &respMap)

	// 验证是否捕获到了 unmarshal error
	assert.Equal(t, "数据类型错误", respMap["age"])
}
