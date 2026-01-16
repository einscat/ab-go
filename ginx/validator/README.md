# ginx/validator

`ginx/validator` 是一个基于 [Gin](https://github.com/gin-gonic/gin) 和 [go-playground/validator](https://github.com/go-playground/validator) 的轻量级参数校验封装库。

它的核心目标是：**让参数校验变得简单、直观且“开箱即用”。**

## ✨ 特性

- **⚡️ 开箱即用**：自动初始化中文翻译器，自动配置 JSON Tag 优先。
- **🌏 自动翻译**：内置中文错误提示，无需手动编写翻译映射逻辑。
- **🛡️ 格式统一**：错误信息以 `map[string]string` 格式返回，Key 为 JSON 字段名，Value 为错误详情，前端友好。
- **🔌 插件化扩展**：提供 `RegisterRule` 接口，允许业务方在各模块中轻松注册自定义正则或逻辑。
- **🚀 零侵入**：完全兼容 Gin 原生 `binding` 标签。

## 📦 安装

```bash
go get -u [github.com/einscat/ab-go](https://github.com/einscat/ab-go)
```