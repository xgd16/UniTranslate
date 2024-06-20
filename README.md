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
- 支持百度、有道、谷歌和 Deepl 腾讯 ChatGPT 火山 讯飞 PaPaGo 平台的翻译接入
- 支持设置翻译 API 的等级优先调用配置的低等级 API
- 同一个 API 提供商可配置不限次 可设置为不同等级
- 在配置多个 API 时如果调用当前 API 失败自动切换到下一个
- 可以将翻译过的内容写入 `Redis` `Memory` 缓存重复翻译内容降低翻译 API 重复调用

## 批量翻译支持情况

|  平台   | 是否支持批量翻译 | 是否完美支持 | 准确的源语言 |                             备注                             |
| :-----: | :--------------: | :----------: | :----------: | :----------------------------------------------------------: |
|  百度   |        是        |      否      |      否      |            不支持精确返回具体每条结果的源语言类型            |
| Google  |        是        |      是      |      是      |                                                              |
|  有道   |        是        |      否      |      否      |                     源语言类型识别不准确                     |
|  火山   |        是        |      是      |      是      |                                                              |
|  Deepl  |        是        |      否      |      否      |                     源语言类型识别不准确                     |
|  讯飞   |        是        |      否      |      否      | 官方不支持批量翻译通过特殊字符 № 切割实现 且 可能出现结果非多条 |
| PaPaGo  |        是        |      否      |      否      |         基于 \n 切割实现 且不可识别不同的源语言类型          |
| ChatGPT |        是        |      是      |      是      |                                                              |

## 未来支持 (优先级按照顺序,打勾为已实现) ✈️
- [x] 持久化已翻译到 `MySQL`
- [x] web 控制页面
- [x] ChatGPT AI翻译
- [x] 讯飞翻译
- [x] 更合理安全的身份验证
- [x] 腾讯翻译
- [x] 火山翻译
- [x] PaPaGo
- [x] 支持更多国家语言
- [x] 支持模拟 `LibreTranslate` 翻译接口
- [x] 支持终端交互翻译
- [ ] 客户端更多翻译功能支持

## 基础类型 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu` `Tencent` `HuoShan` `PaPaGo`

## Docker 启动 🚀
```shell
# 项目目录下
docker build -t uni-translate:latest .
# 然后执行 (最好创建一个 network 将 mysql 和 redis 放在同一个下 然后配置里直接用容器名字访问应用即可)
docker run -d --name uniTranslate -v {本机目录}/config.yaml:/app/config.yaml -p 9431:{你在config.yaml中配置的port} --network baseRun uni-translate:latest
```

## 终端交互方式
在 `config.yaml` 配置完成后执行

```bash
./UniTranslate translate auto en
```

## 配置解析 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis , mem , off 模式 mem 会将翻译结果存储到程序内存中 模式 off 不写入任何缓存
  cachePlatform: false # 执行缓存key生成是否包含平台 (会影响项目启动时自动初始化存储的key)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # http api 对接时的密钥
  keyMode: 1 # 模式 1 直接传入 key 做验证 模式 2 使用 key 加密加签数据进行验证
```

## API 文档 🌍
[在线文档](https://apifox.com/apidoc/shared-335b66b6-90dd-42af-8a1b-f7d1a2c3f351)
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)

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

## 感谢 [Jetbrains](https://www.jetbrains.com/?from=UniTranslate) 提供免费的 IDE
