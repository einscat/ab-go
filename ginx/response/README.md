# ginx/response

`response` 是一个基于 Gin 的 HTTP 响应封装库。它依赖 `ginx/ecode` 来实现自动化的错误处理和统一的 JSON 格式输出。

## ✨ 特性

- **统一格式**：所有响应严格遵循 `{code, msg, data}` 结构。
- **智能错误处理**：`Fail` 方法能自动识别 `ecode.Error` 还是普通 `error`。
- **详情透传**：自动将 `ecode` 中的 `Details` 字段填充到 JSON 的 `data` 字段中（用于表单错误提示）。

## 📖 使用指南

### 1. 成功响应

```go
// 无数据
response.SuccessMsg(c, "操作成功")

// 有数据
response.Success(c, userProfile)
```

### 2. 失败响应

#### 场景 A：业务错误 (自动识别)
```go
// 假设 ecode.UserDuplicate 定义为 2001001
// 响应 HTTP 200: { "code": 2001001, "msg": "用户已存在", "data": null }
response.Fail(c, ecode.UserDuplicate)
```

#### 场景 B：带详情的校验错误
```go
// 响应 HTTP 400: { "code": 10000001, "msg": "参数错误", "data": {"age": "太小了"} }
err := ecode.InvalidParams.WithDetails(map[string]string{"age": "太小了"})
response.Fail(c, err)
```

#### 场景 C：系统未知错误

```go
// 响应 HTTP 500: { "code": 10000000, "msg": "服务内部错误", "data": null }
response.Fail(c, errors.New("db connection lost"))
```