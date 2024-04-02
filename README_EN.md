# uniTranslate

<img src="https://github.com/xgd16/UniTranslate/assets/42709773/3d879e22-fe2c-4238-aabb-39ab478fbd20" alt="UniTranslate" width="300" height="300">

[中文](./README.md) | [English](./README_EN.md)

# Project Introduction 📒
This project is a tool that supports translation across multiple platforms and writes translation results to Redis cache.

## Dependencies
`MySQL: 8.*` `redis`

Optional

`graylog`

## WEB Management
[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## Features ✨
- Supports integration of translation services from Baidu, Youdao, Google, and Deepl platforms.
- Supports setting the priority level of translation APIs for calling lower-level APIs.
- Multiple calls are allowed for the same API provider; different levels can be set.
- Automatically switches to the next API if the current API call fails when configuring multiple APIs.
- Translated content can be written to `Redis` or `Memory` cache to reduce redundant API calls for repeated translation content.

## Future Support (priority in order, checked if implemented) ✈️
- [x] Persistence of translated content to `MySQL`
- [x] Web control page
- [x] ChatGPT AI translation
- [x] XunFei translation
- [x] More reasonable and secure authentication
- [x] Tencent translation
- [x] HuoShan translation
- [x] Support for more languages ​​in different countries
- [ ] More translation feature support for clients

## Basic Types 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu` `Tencent` `HuoShan`

## Docker Startup 🚀
```shell
# In the project directory
docker build -t uni-translate:latest .
# Then execute (it is best to create a network to put mysql and redis under the same network, then directly use the container name to access the application in the configuration)
docker run -d --name uniTranslate -v {local directory}/config.yaml:/app/config.yaml -p 9431:{port configured in your config.yaml} --network baseRun uni-translate:latest
```

## Configuration Parsing 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis, mem, off modes mem will store translation results in program memory, mode off does not write any cache
  cachePlatform: false # Does the execution cache key generation contain the platform (will affect the automatic initialization of stored keys when the project starts)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # Key for http API docking authentication
  keyMode: 1 # Mode 1 directly passes in the key for verification, mode 2 uses key encryption and signing data for verification
```

## Interface Authentication ts Example
```typescript
import { MD5 } from "crypto-js";

/**
 * 
 * @param key Key set by the platform
 * @param params Request parameters
 * @return Generated authentication code
 */
function AuthEncrypt(key: string, params: { [key: string]: any }): string {
  return MD5(key + sortMapToStr(params)).toString();
}


const sortMapToStr = (map: { [key: string]: any }): string => {
  let mapArr = new Array();
  for (const key in map) {
    const item = map[key];
    if (Array.isArray(item)) {
      mapArr.push(`${key}:${item.join(",")}`);
      continue;
    }
    if (typeof item === "object") {
      mapArr.push(`${key}:|${sortMapToStr(item)}|`);
      continue;
    }
    mapArr.push(`${key}:${item}`);
  }

  return mapArr.sort().join("&");
};

const params: { [key: string]: any } = {
    c: {
        cc: 1,
        cb: 2,
        ca: 3,
        cd: 4,
    },
    a: 1,
    b: [4, 1, 2],
};

console.log(AuthEncrypt("123456", params));
```

Request Example

```shell
curl --location --request POST 'http://127.0.0.1:9431/api/translate' \
--header 'auth_key: xxxxxxxxx{place AuthEncrypt function result here}' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--data '{
    "from": "auto",
    "to": "en",
    "text": "Test",
    "platform": "YouDao"
}'
```


## Unsupported Translated Content ??? 🤔
All supported languages ​​in this program are uniformly marked based on the identifier of the _Youdao_ translation API in the [translate.json](./translate.json) file.

Please modify the `translate.json` file based on the identifiers supported by the _Youdao_ translation API documentation.

## Basic Language Identifiers

| English Name              | Chinese Name     | Code      |
|-------------------------|-------------|---------|
| Arabic                  | 阿拉伯语        | ar      |
| German                  | 德语          | de      |
| English                 | 英语          | en      |
| Spanish                 | 西班牙语        | es      |
| French                  | 法语          | fr      |
| Hindi                   | 印地语         | hi      |
| Indonesian              | 印度尼西亚语      | id      |
| Italian                 | 意大利语        | it      |
| Japanese                | 日语          | ja      |
| Korean                  | 韩语          | ko      |
| Dutch                   | 荷兰语         | nl      |
| Portuguese              | 葡萄牙语        | pt      |
| Russian                 | 俄语          | ru      |
| Thai                    | 泰语          | th      |
| Vietnamese              | 越南语         | vi      |
| Chinese                 | 简体中文        | zh-CHS  |
| Chinese                 | 繁体中文        | zh-CHT  |
| Afrikaans               | 南非荷兰语       | af      |
| Amharic                 | 阿姆哈拉语       | am      |
| Azerbaijani             | 阿塞拜疆语       | az      |
| Belarusian              | 白俄罗斯语       | be      |
| Bulgarian               | 保加利亚语       | bg      |
| Bengali                 | 孟加拉语        | bn      |
| Bosnian (Latin)         | 波斯尼亚语       | bs      |
| Catalan                 | 加泰隆语        | ca      |
| Cebuano                 | 宿务语         | ceb     |
| Corsican                | 科西嘉语        | co      |
| Czech                   | 捷克语         | cs      |
| Welsh                   | 威尔士语        | cy      |
| Danish                  | 丹麦语         | da      |
| Greek                   | 希腊语         | el      |
| Esperanto               | 世界语         | eo      |
| Estonian                | 爱沙尼亚语       | et      |
| Basque                  | 巴斯克语        | eu      |
| Persian                 | 波斯语         | fa      |
| Finnish                 | 芬兰语         | fi      |
| Fijian                  | 斐济语         | fj      |
| Frisian                 | 弗里西语        | fy      |
| Irish                   | 爱尔兰语        | ga      |
| Scots                   | 苏格兰盖尔语      | gd      |
| Galician                | 加利西亚语       | gl      |
| Gujarati                | 古吉拉特语       | gu      |
| Hausa                   | 豪萨语         | ha      |
| Hawaiian                | 夏威夷语        | haw     |
| Hebrew                  | 希伯来语        | he      |
| Hindi                   | 印地语         | hi      |
| Croatian                | 克罗地亚语       | hr      |
| Haitian                 | 海地克里奥尔语     | ht      |
| Hungarian               | 匈牙利语        | hu      |
| Armenian                | 亚美尼亚语       | hy      |
| Igbo                    | 伊博语         | ig      |
| Icelandic               | 冰岛语         | is      |
| Javanese                | 爪哇语         | jw      |
| Georgian                | 格鲁吉亚语       | ka      |
| Kazakh                  | 哈萨克语        | kk      |
| Khmer                   | 高棉语         | km      |
| Kannada                 | 卡纳达语        | kn      |
| Kurdish                 | 库尔德语        | ku      |
| Kyrgyz                  | 柯尔克孜语       | ky      |
| Latin                   | 拉丁语         | la      |
| Luxembourgish           | 卢森堡语        | lb      |
| Lao                     | 老挝语         | lo      |
| Lithuanian              | 立陶宛语        | lt      |
| Latvian                 | 拉脱维亚语       | lv      |
| Malagasy                | 马尔加什语       | mg      |
| Maori                   | 毛利语         | mi      |
| Macedonian              | 马其顿语        | mk      |
| Malayalam               | 马拉雅拉姆语      | ml      |
| Mongolian               | 蒙古语         | mn      |
| Marathi                 | 马拉地语        | mr      |
| Malay                   | 马来语         | ms      |
| Maltese                 | 马耳他语        | mt      |
| Hmong                   | 白苗语         | mww     |
| Myanmar (Burmese)       | 缅甸语         | my      |
| Nepali                  | 尼泊尔语        | ne      |
| Norwegian               | 挪威语         | no      |
| Nyanja (Chichewa)       | 齐切瓦语        | ny      |
| Querétaro Otomi         | 克雷塔罗奥托米语    | otq     |
| Punjabi                 | 旁遮普语        | pa      |
| Polish                  | 波兰语         | pl      |
| Pashto                  | 普什图语        | ps      |
| Romanian                | 罗马尼亚语       | ro      |
| Sindhi                  | 信德语         | sd      |
| Sinhala (Sinhalese)     | 僧伽罗语        | si      |
| Slovak                  | 斯洛伐克语       | sk      |
| Slovenian               | 斯洛文尼亚语      | sl      |
| Samoan                  | 萨摩亚语        | sm      |
| Shona                   | 修纳语         | sn      |
| Somali                  | 索马里语        | so      |
| Albanian                | 阿尔巴尼亚语      | sq      |
| Serbian (Cyrillic)      | 塞尔维亚语(西里尔文) | sr-Cyrl |
| Serbian (Latin)         | 塞尔维亚语(拉丁文)  | sr-Latn |
| Sesotho                 | 塞索托语        | st      |
| Sundanese               | 巽他语         | su      |
| Swedish                 | 瑞典语         | sv      |
| Kiswahili               | 斯瓦希里语       | sw      |
| Tamil                   | 泰米尔语        | ta      |
| Telugu                  | 泰卢固语        | te      |
| Tajik                   | 塔吉克语        | tg      |
| Filipino                | 菲律宾语        | tl      |
| Klingon                 | 克林贡语        | tlh     |
| Tongan                  | 汤加语         | to      |
| Turkish                 | 土耳其语        | tr      |
| Tahitian                | 塔希提语        | ty      |
| Ukrainian               | 乌克兰语        | uk      |
| Urdu                    | 乌尔都语        | ur      |
| Uzbek                   | 乌兹别克语       | uz      |
| Xhosa                   | 南非科萨语       | xh      |
| Yiddish                 | 意第绪语        | yi      |
| Yoruba                  | 约鲁巴语        | yo      |
| Yucatec                 | 尤卡坦玛雅语      | yua     |
| Cantonese (Traditional) | 粤语          | yue     |
| Zulu                    | 南

非祖鲁语       | zu      |
| 自动识别                    | auto        |         |

## API Documentation 🌍
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)