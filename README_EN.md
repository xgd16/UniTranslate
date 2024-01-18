# uniTranslate

<img src="https://github.com/xgd16/UniTranslate/assets/42709773/3d879e22-fe2c-4238-aabb-39ab478fbd20" alt="UniTranslate" width="300" height="300">

[中文](./README.md) | [English](./README_EN.md)

# Project Introduction 📒
This project is a tool that supports translation on multiple platforms and writes the translation results to the Redis cache.

## Dependencies
`MySQL: 8.*` `redis`

Optional

`graylog`

## WEB Management
[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## Features ✨
- Supports translation integration with platforms such as Baidu, Youdao, Google, and Deepl.
- Allows setting the translation API's priority level for invoking lower-level APIs configured.
- Configurable unlimited calls for the same API provider; can be set to different levels.
- Automatically switches to the next API if the current API call fails when configuring multiple APIs.
- Translated content can be written to `Redis` or `Memory` cache to reduce repetitive translation API calls.

## Future Support (Priority in Order, ✔️ indicates implemented)
- [x] Persist translated content to `MySQL`.
- [x] Web control page.
- [x] ChatGPT AI translation.
- [x] XunFei translation.
- [ ] More secure identity authentication.
- [ ] Tencent translation.
- [ ] Support for more languages from different countries.
- [ ] More translation functionality support for the client.

## Base Types 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu`

## Docker Launch 🚀
```shell
# In the project directory
docker build -t uni-translate:latest .
# Then execute (it's better to create a network and put mysql and redis in the same one, then directly access the app using container names in the configuration)
docker run -d --name uniTranslate -v {local directory}/config.yaml:/app/config.yaml -p 9431:{port configured in config.yaml} --network baseRun uni-translate:latest
```

## Configuration Parsing 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis, mem, off modes; mem stores translation results in program memory, off does not write to any cache
  cachePlatform: false # Does cache key generation include platform (affects automatic initialization of stored keys when the project starts)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # Key for HTTP API integration
```


## Unsupported Translations 🤔
All supported languages in this program are based on the **identifiers** of languages in the [translate.json](./translate.json) file, using the identifier of the Youdao translation API as the reference.

Modify the `translate.json` file based on the identifiers supported by the Youdao translation API documentation.

## Base Language Identifiers

Here are the results of converting the data you provided into a Markdown table:

| English Name         | Chinese Name  | Code    |
|---------------------|----------------|---------|
| Arabic              | 阿拉伯语         | ar      |
| German              | 德语            | de      |
| English             | 英语            | en      |
| Spanish             | 西班牙语         | es      |
| French              | 法语            | fr      |
| Hindi               | 印地语           | hi      |
| Indonesian          | 印度尼西亚语     | id      |
| Italian             | 意大利语         | it      |
| Japanese            | 日语            | ja      |
| Korean              | 韩语            | ko      |
| Dutch               | 荷兰语           | nl      |
| Portuguese          | 葡萄牙语         | pt      |
| Russian             | 俄语            | ru      |
| Thai                | 泰语            | th      |
| Vietnamese          | 越南语           | vi      |
| Chinese (Simplified)| 简体中文         | zh-CHS  |
| Chinese (Traditional)| 繁体中文       | zh-CHT  |
| Afrikaans           | 南非荷兰语       | af      |
| Amharic             | 阿姆哈拉语       | am      |
| Azerbaijani         | 阿塞拜疆语       | az      |
| Belarusian          | 白俄罗斯语       | be      |
| Bulgarian           | 保加利亚语       | bg      |
| Bengali             | 孟加拉语         | bn      |
| Bosnian (Latin)     | 波斯尼亚语       | bs      |
| Catalan             | 加泰隆语         | ca      |
| Cebuano             | 宿务语           | ceb     |
| Corsican            | 科西嘉语         | co      |
| Czech               | 捷克语           | cs      |
| Welsh               | 威尔士语         | cy      |
| Danish              | 丹麦语           | da      |
| Greek               | 希腊语           | el      |
| Esperanto           | 世界语           | eo      |
| Estonian            | 爱沙尼亚语       | et      |
| Basque              | 巴斯克语         | eu      |
| Persian             | 波斯语           | fa      |
| Finnish             | 芬兰语           | fi      |
| Fijian              | 斐济语           | fj      |
| Frisian             | 弗里西语         | fy      |
| Irish               | 爱尔兰语         | ga      |
| Scots               | 苏格兰盖尔语     | gd      |
| Galician            | 加利西亚语       | gl      |
| Gujarati            | 古吉拉特语       | gu      |
| Hausa               | 豪萨语           | ha      |
| Hawaiian            | 夏威夷语         | haw     |
| Hebrew              | 希伯来语         | he      |
| Hindi               | 印地语           | hi      |
| Croatian            | 克罗地亚语       | hr      |
| Haitian             | 海地克里奥尔语   | ht      |
| Hungarian           | 匈牙利语         | hu      |
| Armenian            | 亚美尼亚语       | hy      |
| Igbo                | 伊博语           | ig      |
| Icelandic           | 冰岛语           | is      |
| Javanese            | 爪哇语           | jw      |
| Georgian            | 格鲁吉亚语       | ka      |
| Kazakh              | 哈萨克语         | kk      |
| Khmer               | 高棉语           | km      |
| Kannada             | 卡纳达语         | kn      |
| Kurdish             | 库尔德语         | ku      |
| Kyrgyz              | 柯尔克孜语       | ky      |
| Latin               | 拉丁语           | la      |
| Luxembourgish       | 卢森堡语         | lb      |
| Lao                 | 老挝语           | lo      |
| Lithuanian          | 立陶宛语         | lt      |
| Latvian             | 拉脱维亚语       | lv      |
| Malagasy            | 马尔加什语       | mg      |
| Maori               | 毛利语           | mi      |
| Macedonian          | 马其顿语         | mk      |
| Malayalam           | 马拉雅拉姆语     | ml      |
| Mongolian           | 蒙古语           | mn      |
| Marathi             | 马拉地语         | mr      |
| Malay               | 马来语           | ms      |
| Maltese             | 马耳他语         | mt      |
| Hmong               | 白苗语           | mww     |
| Myanmar (Burmese)   | 缅甸语           | my      |
| Nepali              | 尼泊尔语         | ne      |
| Norwegian           | 挪威语           | no      |
| Nyanja (Chichewa)   | 齐切瓦语         | ny      |
| Querétaro Otomi     | 克雷塔罗奥托米语 | otq    |
| Punjabi             | 旁遮普语         | pa      |
| Polish              | 波兰语           | pl      |
| Pashto              | 普什图语         | ps      |
| Romanian            | 罗马尼亚语       | ro      |
| Sindhi              | 信德语           | sd      |
| Sinhala (Sinhalese) | 僧伽罗语         | si      |
| Slovak              | 斯洛伐克语       | sk      |
| Slovenian           | 斯洛文尼亚语     | sl      |
| Samoan              | 萨摩亚语         | sm      |
| Shona               | 修纳语           | sn      |
| Somali              | 索马里语         | so      |
| Albanian            | 阿尔巴尼亚语     | sq      |
| Serbian (Cyrillic)  | 塞尔维亚语(西里尔文)| sr-Cyrl |
| Serbian (Latin)     | 塞尔维亚语(拉丁文)| sr-Latn |
| Sesotho             | 塞索托语         | st      |
| Sundanese           | 巽他语           | su      |
| Swedish             | 瑞典语           | sv      |
| Kiswahili           | 斯瓦希里语       | sw      |
| Tamil               | 泰米尔语         | ta      |
| Telugu              | 泰卢固语         | te      |
| Tajik               | 塔吉克语         | tg      |
| Filipino            | 菲律宾语         | tl      |
| Klingon             | 克林贡语         | tlh     |
| Tongan              | 汤加语           | to      |
| Turkish             | 土耳其语         | tr      |
| Tahitian            | 塔希提语         | ty      |
| Ukrainian           | 乌克兰语         | uk      |
| Urdu                | 乌尔都语         | ur      |
| Uzbek               | 乌兹别克语       | uz      |
| Xhosa               | 南非科萨语       | xh      |
| Yiddish             | 意第绪语         | yi      |
| Yoruba              | 约鲁巴语         | yo      |
| Yucatec             | 尤卡坦玛雅语     | yua     |
| Cantonese (Traditional)| 粤语          | yue     |
| Zulu                | 南非祖鲁语       | zu      |
| Auto-detect         | 自动识别         | auto    |

## API Documentation 🌍
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)