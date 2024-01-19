# uniTranslate

<img src="https://github.com/xgd16/UniTranslate/assets/42709773/3d879e22-fe2c-4238-aabb-39ab478fbd20" alt="UniTranslate" width="300" height="300">

[中文](./README.md) | [English](./README_EN.md)

# 项目简介 📒
该项目是一个支持多平台翻译和将翻译结果写入 Redis 缓存的工具。

## 依赖
`MySQL: 8.*` `redis`

可选

`graylog`

## WEB管理
[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## 功能特点 ✨
- 支持百度、有道、谷歌和 Deepl 平台的翻译接入
- 支持设置翻译 API 的等级优先调用配置的低等级 API
- 同一个 API 提供商可配置不限次 可设置为不同等级
- 在配置多个 API 时如果调用当前 API 失败自动切换到下一个
- 可以将翻译过的内容写入 `Redis` `Memory` 缓存重复翻译内容降低翻译 API 重复调用

## 未来支持 (优先级按照顺序,打勾为已实现) ✈️
- [x] 持久化已翻译到 `MySQL`
- [x] web 控制页面
- [x] ChatGPT AI翻译
- [x] 讯飞翻译
- [x] 更合理安全的身份验证
- [ ] 腾讯翻译
- [ ] 支持更多国家语言
- [ ] 客户端更多翻译功能支持

## 基础类型 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu`

## Docker 启动 🚀
```shell
# 项目目录下
docker build -t uni-translate:latest .
# 然后执行 (最好创建一个 network 将 mysql 和 redis 放在同一个下 然后配置里直接用容器名字访问应用即可)
docker run -d --name uniTranslate -v {本机目录}/config.yaml:/app/config.yaml -p 9431:{你在config.yaml中配置的port} --network baseRun uni-translate:latest
```


## 配置解析 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis , mem , off 模式 mem 会将翻译结果存储到程序内存中 模式 off 不写入任何缓存
  cachePlatform: false # 执行缓存key生成是否包含平台 (会影响项目启动时自动初始化存储的key)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # http api 对接时的密钥
```

## 接口身份验证 ts 示例
```typescript
import { MD5 } from "crypto-js";

/**
 * 
 * @param key 平台设置的key
 * @param params 请求参数
 * @return 生成的身份验证码
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

请求示例

```shell
curl --location --request POST 'http://127.0.0.1:9431/api/translate' \
--header 'auth_key: xxxxxxxxx{AuthEncrypt函数结果放在此处}' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--data '{
    "from": "auto",
    "to": "en",
    "text": "测试一下",
    "platform": "YouDao"
}'
```


## 翻译的内容不支持??? 🤔
本程序所有支持的语言根据 [translate.json](./translate.json) 文件进行国家语言**标识**统一使用 _有道_ 翻译 API 标识符作为基准

请根据 _有道_ 翻译 API 文档支持的标识作为基准修改 `translate.json` 文件

## 基础语言标识

| 英文名             | 中文名         | 代码     |
|------------------|--------------|---------|
| Arabic           | 阿拉伯语       | ar      |
| German           | 德语           | de      |
| English          | 英语           | en      |
| Spanish          | 西班牙语       | es      |
| French           | 法语           | fr      |
| Hindi            | 印地语         | hi      |
| Indonesian       | 印度尼西亚语   | id      |
| Italian          | 意大利语       | it      |
| Japanese         | 日语           | ja      |
| Korean           | 韩语           | ko      |
| Dutch            | 荷兰语         | nl      |
| Portuguese       | 葡萄牙语       | pt      |
| Russian          | 俄语           | ru      |
| Thai             | 泰语           | th      |
| Vietnamese       | 越南语         | vi      |
| Chinese          | 简体中文       | zh-CHS  |
| Chinese          | 繁体中文       | zh-CHT  |
| Afrikaans        | 南非荷兰语     | af      |
| Amharic          | 阿姆哈拉语     | am      |
| Azerbaijani      | 阿塞拜疆语     | az      |
| Belarusian       | 白俄罗斯语     | be      |
| Bulgarian        | 保加利亚语     | bg      |
| Bengali          | 孟加拉语       | bn      |
| Bosnian (Latin)  | 波斯尼亚语     | bs      |
| Catalan          | 加泰隆语       | ca      |
| Cebuano          | 宿务语         | ceb     |
| Corsican         | 科西嘉语       | co      |
| Czech            | 捷克语         | cs      |
| Welsh            | 威尔士语       | cy      |
| Danish           | 丹麦语         | da      |
| Greek            | 希腊语         | el      |
| Esperanto        | 世界语         | eo      |
| Estonian         | 爱沙尼亚语     | et      |
| Basque           | 巴斯克语       | eu      |
| Persian          | 波斯语         | fa      |
| Finnish          | 芬兰语         | fi      |
| Fijian           | 斐济语         | fj      |
| Frisian          | 弗里西语       | fy      |
| Irish            | 爱尔兰语       | ga      |
| Scots            | 苏格兰盖尔语   | gd      |
| Galician         | 加利西亚语     | gl      |
| Gujarati         | 古吉拉特语     | gu      |
| Hausa            | 豪萨语         | ha      |
| Hawaiian         | 夏威夷语       | haw     |
| Hebrew           | 希伯来语       | he      |
| Hindi            | 印地语         | hi      |
| Croatian         | 克罗地亚语     | hr      |
| Haitian          | 海地克里奥尔语 | ht      |
| Hungarian        | 匈牙利语       | hu      |
| Armenian         | 亚美尼亚语     | hy      |
| Igbo             | 伊博语         | ig      |
| Icelandic        | 冰岛语         | is      |
| Javanese         | 爪哇语         | jw      |
| Georgian         | 格鲁吉亚语     | ka      |
| Kazakh           | 哈萨克语       | kk      |
| Khmer            | 高棉语         | km      |
| Kannada          | 卡纳达语       | kn      |
| Kurdish          | 库尔德语       | ku      |
| Kyrgyz           | 柯尔克孜语     | ky      |
| Latin            | 拉丁语         | la      |
| Luxembourgish    | 卢森堡语       | lb      |
| Lao              | 老挝语         | lo      |
| Lithuanian       | 立陶宛语       | lt      |
| Latvian          | 拉脱维亚语     | lv      |
| Malagasy         | 马尔加什语     | mg      |
| Maori            | 毛利语         | mi      |
| Macedonian       | 马其顿语       | mk      |
| Malayalam        | 马拉雅拉姆语   | ml      |
| Mongolian        | 蒙古语         | mn      |
| Marathi          | 马拉地语       | mr      |
| Malay            | 马来语         | ms      |
| Maltese          | 马耳他语       | mt      |
| Hmong            | 白苗语         | mww     |
| Myanmar (Burmese)| 缅甸语         | my      |
| Nepali           | 尼泊尔语       | ne      |
| Norwegian        | 挪威语         | no      |
| Nyanja (Chichewa)| 齐切瓦语       | ny      |
| Querétaro Otomi  | 克雷塔罗奥托米语 | otq    |
| Punjabi          | 旁遮普语       | pa      |
| Polish           | 波兰语         | pl      |
| Pashto           | 普什图语       | ps      |
| Romanian         | 罗马尼亚语     | ro      |
| Sindhi           | 信德语         | sd      |
| Sinhala (Sinhalese)| 僧伽罗语     | si      |
| Slovak           | 斯洛伐克语     | sk      |
| Slovenian        | 斯洛文尼亚语   | sl      |
| Samoan           | 萨摩亚语       | sm      |
| Shona            | 修纳语         | sn      |
| Somali           | 索马里语       | so      |
| Albanian         | 阿尔巴尼亚语   | sq      |
| Serbian (Cyrillic)| 塞尔维亚语(西里尔文)| sr-Cyrl |
| Serbian (Latin)  | 塞尔维亚语(拉丁文)| sr-Latn |
| Sesotho          | 塞索托语       | st      |
| Sundanese        | 巽他语         | su      |
| Swedish          | 瑞典语         | sv      |
| Kiswahili        | 斯瓦希里语     | sw      |
| Tamil            | 泰米尔语       | ta      |
| Telugu           | 泰卢固语       | te      |
| Tajik            | 塔吉克语       | tg      |
| Filipino         | 菲律宾语       | tl      |
| Klingon          | 克林贡语       | tlh     |
| Tongan           | 汤加语         | to      |
| Turkish          | 土耳其语       | tr      |
| Tahitian         | 塔希提语       | ty      |
| Ukrainian        | 乌克兰语       | uk      |
| Urdu             | 乌尔都语       | ur      |
| Uzbek            | 乌兹别克语     | uz      |
| Xhosa            | 南非科萨语     | xh      |
| Yiddish          | 意第绪语       | yi      |
| Yoruba           | 约鲁巴语       | yo      |
| Yucatec          | 尤卡坦玛雅语   | yua     |
| Cantonese (Traditional)| 粤语   | yue     |
| Zulu             | 南非祖鲁语     | zu      |
| 自动识别         | auto          |         |

## API 文档 🌍
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)
