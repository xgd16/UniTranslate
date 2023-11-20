# uniTranslate

[中文](./README.md) | [English](./README_EN.md)

# Project Overview 📒
This project is a tool that supports translation across multiple platforms and writes the translation results into Redis cache.

## Key Features ✨
- Supports integration with translation platforms like Baidu, Youdao, Google, and Deepl
- Supports setting the priority level for calling translation APIs, allowing lower-level APIs to be configured for priority usage
- Configuration for unlimited calls for the same API provider, adjustable at different levels
- When configuring multiple APIs, automatically switches to the next one if the current API call fails
- Ability to write translated content into `Redis` cache to reduce repeated calls to translation APIs for duplicated content

## Future Support (priority in order) ✈️
- [ ] Persist translated content into `MySQL`
- [ ] Web control panel

## Base Types 🪨
`YouDao` `Baidu` `Google` `Deepl`

## Unsupported Content for Translation??? 🤔
All supported languages in this program are uniformly identified based on the _Youdao_ translation API identifier, according to the country language **identification** in the [translate.json](./translate.json) file.

Please modify the `translate.json` file based on the identification supported by the _Youdao_ translation API documentation.

## API Documentation 🌍
[Open Api File](./uniTranslate%20(Unified%20Translation).openapi.json)