# uniTranslate

[中文](./README.md) | [English](./README_EN.md)

# 项目简介 📒
该项目是一个支持多平台翻译和将翻译结果写入 Redis 缓存的工具。

## 功能特点 ✨
- 支持百度、有道、谷歌和 Deepl 平台的翻译接入
- 支持设置翻译 API 的等级优先调用配置的低等级 API
- 同一个 API 提供商可配置不限次 可设置为不同等级
- 在配置多个 API 时如果调用当前 API 失败自动切换到下一个
- 可以将翻译过的内容写入 `Redis` 缓存重复翻译内容降低翻译 API 重复调用

## 未来支持 (优先级按照顺序) ✈️
- [ ] 持久化已翻译到 `MySQL`
- [ ] web 控制页面

## 基础类型 🪨
`YouDao` `Baidu` `Google` `Deepl`

## 翻译的内容不支持??? 🤔
本程序所有支持的语言根据 [translate.json](./translate.json) 文件进行国家语言**标识**统一使用 _有道_ 翻译 API 标识符作为基准

请根据 _有道_ 翻译 API 文档支持的标识作为基准修改 `translate.json` 文件

## API 文档 🌍
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)
