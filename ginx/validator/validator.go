package validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 包级单例，保证全局唯一
var (
	trans ut.Translator
	once  sync.Once
)

// init 自动初始化：配置 JSON Tag 读取器和中文翻译器
// 引用此包时自动执行，无需用户手动初始化
func init() {
	once.Do(func() {
		if v, ok := binding.Validator.Engine().(*val.Validate); ok {
			// 1. 注册 Tag Name 函数：优先读取 json tag
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})

			// 2. 注册中文翻译
			zhT := zh.New()
			uni := ut.New(zhT, zhT)
			trans, _ = uni.GetTranslator("zh")
			_ = zh_translations.RegisterDefaultTranslations(v, trans)
		}
	})
}

// RegisterRule 提供给业务方注册自定义规则的入口
// tag: 验证标签，如 "mobile"
// msg: 错误提示模版，如 "手机号格式不正确"
// fn:  验证函数逻辑
func RegisterRule(tag string, fn val.Func, msg string) {
	if v, ok := binding.Validator.Engine().(*val.Validate); ok {
		// 注册验证逻辑
		_ = v.RegisterValidation(tag, fn)

		// 注册翻译逻辑
		_ = v.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
			// true 表示覆盖可能存在的默认翻译
			return ut.Add(tag, msg, true)
		}, func(ut ut.Translator, fe val.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		})
	}
}

// BindAndValid 通用绑定验证器
// 返回值 bool: true=通过, false=失败
// 返回值 map:  失败详情 (key=字段名, value=错误信息)
func BindAndValid(c *gin.Context, obj any) (bool, map[string]string) {
	err := c.ShouldBind(obj)
	if err == nil {
		return true, nil
	}

	errs := make(map[string]string)

	// 1. 处理校验规则错误 (ValidationErrors)
	var verrs val.ValidationErrors
	if errors.As(err, &verrs) {
		for key, value := range verrs.Translate(trans) {
			// key 格式通常是 "StructName.FieldName"
			// 我们需要去掉前缀，只保留字段名
			k := key
			if idx := strings.Index(k, "."); idx != -1 {
				k = k[idx+1:]
			}
			errs[k] = value
		}
		return false, errs
	}

	// 2. 处理 JSON 类型不匹配
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		errs[ute.Field] = "数据类型错误"
		return false, errs
	}

	// 3. 其他错误
	errs["request"] = "请求参数格式错误"
	return false, errs
}
