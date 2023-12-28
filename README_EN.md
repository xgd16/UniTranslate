# uniTranslate

<img src="https://github.com/xgd16/UniTranslate/assets/42709773/3d879e22-fe2c-4238-aabb-39ab478fbd20" alt="UniTranslate" width="300" height="300">

[中文](./README.md) | [English](./README_EN.md)

# Project Introduction 📒
This project is a tool that supports multi-platform translation and writes translation results into Redis cache.

## WEB Management
[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## Features ✨
- Support translation integration with Baidu, Youdao, Google, and Deepl platforms
- Support setting the priority of translation APIs, favoring lower-level APIs in the configuration
- Configuration of unlimited calls for the same API provider; can be set to different levels
- When configuring multiple APIs, automatically switch to the next one if the current API call fails
- Translated content can be written into `Redis` cache to reduce repetitive calls to the translation API

## Future Support (Priority in order, ✔️ indicates implemented) ✈️
- [x] Persist translated content to `MySQL`
- [x] Web control page
- [x] ChatGPT AI translate

## Basic Types 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu`

## Configuration Parsing 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis, mem, off modes where mem stores translation results in program memory and off doesn't write any cache
  cachePlatform: false # Whether to include the platform in cache key generation (affects automatic initialization of stored keys when the project starts)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # Key for HTTP API integration
```

## Unsupported Content for Translation??? 🤔
All supported languages in this program are unified using the _Youdao_ translation API identifier as a reference based on the `translate.json` file.

Please modify the `translate.json` file based on the identifiers supported in the _Youdao_ translation API documentation.

## API Documentation 🌍
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)