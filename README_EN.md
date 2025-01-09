# uniTranslate

<img src="./logo.svg" alt="UniTranslate" width="300" height="300">

[中文](./README.md) | [English](./README_EN.md)

# Project Introduction 📒

This project is a tool that supports multi-platform translation and writing translation results into Redis cache.

## Dependencies

`MySQL: 8.*` `redis`

Optional

`graylog`

## WEB Management

[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## Features ✨

- Supports translation integration with Baidu, Youdao, Google, Deepl, Tencent, ChatGPT, Volcano, iFlytek, PaPaGo, and free Google platforms
- Allows setting translation API priorities to call lower priority APIs first
- Configurable unlimited instances for the same API provider at different priority levels
- Automatically switches to the next API if the current one fails when multiple APIs are configured
- Can write translated content into `Redis` and `Memory` cache to reduce repeated API calls for the same content

## Batch Translation Support

|    Platform    | Batch Translation Supported | Perfectly Supported | Accurate Source Language |                      Remarks                       |
| :------------: | :-------------------------: | :-----------------: | :----------------------: | :------------------------------------------------: |
|      Baidu     |             Yes             |         No          |           No             | Does not support accurately returning the source language for each result |
|     Google     |             Yes             |         Yes         |           Yes            |                                                    |
|     Youdao     |             Yes             |         No          |           No             |             Inaccurate source language recognition              |
|     Volcano    |             Yes             |         Yes         |           Yes            |                                                    |
|     Deepl      |             Yes             |         No          |           No             |             Inaccurate source language recognition              |
|     iFlytek    |             Yes             |         Yes         |           Yes            |                    Loop implementation                     |
|     PaPaGo     |             Yes             |         No          |           No             | Based on \n split implementation and cannot recognize different source languages |
|    ChatGPT     |             Yes             |         Yes         |           Yes            |                                                    |
|   FreeGoogle   |             Yes             |         Yes         |           Yes            |                    Loop implementation                     |

## Future Support (Priority in order, checked items are implemented) ✈️

- [x] Persist translated content to `MySQL`
- [x] Web control panel
- [x] ChatGPT AI translation
- [x] iFlytek translation
- [x] More reasonable and secure authentication
- [x] Tencent translation
- [x] Volcano translation
- [x] PaPaGo
- [x] Support for more languages
- [x] Support for simulating `LibreTranslate` translation interface
- [x] Support for terminal interactive translation
- [x] Free Google translation
- [x] SQL Lite support
- [ ] More client translation features support

## Basic Types 🪨

`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu` `Tencent` `HuoShan` `PaPaGo` `FreeGoogle`

## Docker Startup 🚀

```shell
# In the project directory
docker build -t uni-translate:latest .
# Then execute (it's better to create a network to put MySQL and Redis under the same network, and then access the application directly by the container name in the configuration)
docker run -d --name uniTranslate -v {local_directory}/config.yaml:/app/config.yaml -p 9431:{port_in_config.yaml} --network baseRun uni-translate:latest
```

## Terminal Interactive Mode

After completing the `config.yaml` configuration, execute

```bash
./UniTranslate translate auto en
```

## Configuration Parsing 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis, mem, off mode; mem stores translation results in program memory; off does not write to any cache
  cachePlatform: false # Whether to include the platform in the cache key generation (affects the key stored during project initialization)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # Key for HTTP API integration
  keyMode: 1 # Mode 1: direct key verification; Mode 2: use key to encrypt and sign data for verification
```

## API Documentation 🌍

[Online Documentation](https://apifox.com/apidoc/shared-335b66b6-90dd-42af-8a1b-f7d1a2c3f351)  
[Open API File](<./uniTranslate%20(统一翻译).openapi.json>)

## Interface Authentication TS Example

```typescript
import { MD5 } from "crypto-js";

/**
 *
 * @param key The key set on the platform
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
--header 'auth_key: xxxxxxxxx{AuthEncrypt function result here}' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--data '{
    "from": "auto",
    "to": "en",
    "text": "测试一下",
    "platform": "YouDao"
}'
```

## Unsupported Translation Content??? 🤔

All supported languages in this program are based on the country language identifiers in the [translate.json](./translate.json) file, using _Youdao_ translation API identifiers as the baseline.

Please refer to the identifiers supported by the _Youdao_ translation API documentation to modify the `translate.json` file.

## Thanks to [Jetbrains](https://www.jetbrains.com/?from=UniTranslate) for providing a free IDE