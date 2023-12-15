# uniTranslate

[中文](./README.md) | [English](./README_EN.md)

# 项目简介 📒
该项目是一个支持多平台翻译和将翻译结果写入 Redis 缓存的工具。

## WEB管理
[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## 功能特点 ✨
- 支持百度、有道、谷歌和 Deepl 平台的翻译接入
- 支持设置翻译 API 的等级优先调用配置的低等级 API
- 同一个 API 提供商可配置不限次 可设置为不同等级
- 在配置多个 API 时如果调用当前 API 失败自动切换到下一个
- 可以将翻译过的内容写入 `Redis` 缓存重复翻译内容降低翻译 API 重复调用

## 未来支持 (优先级按照顺序,打勾为已实现) ✈️
- [x] 持久化已翻译到 `MySQL`
- [x] web 控制页面
- [x] ChatGPT AI翻译

## 基础类型 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT`



## 配置解析 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis , mem , off 模式 mem 会将翻译结果存储到程序内存中 模式 off 不写入任何缓存
  cachePlatform: false # 执行缓存key生成是否包含平台 (会影响项目启动时自动初始化存储的key)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # http api 对接时的密钥
```



## 翻译的内容不支持??? 🤔
本程序所有支持的语言根据 [translate.json](./translate.json) 文件进行国家语言**标识**统一使用 _有道_ 翻译 API 标识符作为基准

请根据 _有道_ 翻译 API 文档支持的标识作为基准修改 `translate.json` 文件

## API 文档 🌍
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)
