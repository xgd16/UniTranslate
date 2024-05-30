# uniTranslate

<img src="https://github.com/xgd16/UniTranslate/assets/42709773/3d879e22-fe2c-4238-aabb-39ab478fbd20" alt="UniTranslate" width="300" height="300">

[中文](./README.md) | [English](./README_EN.md)

# Project Introduction 📒
This project is a tool that supports multi-platform translation and writes translation results into Redis cache.

## Dependencies
`MySQL: 8.*` `redis`

Optional

`graylog`

## Web Management
[UniTranslate-web-console](https://github.com/xgd16/UniTranslate-web-console)

## Features ✨
- Supports translation integration for platforms including Baidu, Youdao, Google, Deepl, Tencent, ChatGPT, Volcano, iFLYTEK, and PaPaGo
- Allows setting translation API priority levels, configuring lower-level APIs to be called first
- Unlimited configurations for the same API provider, can be set at different levels
- Automatically switches to the next API if the current one fails when multiple APIs are configured
- Can write translated content into `Redis` or `Memory` cache to reduce repeated calls to the translation API

## Batch Translation Support

|  Platform  | Supports Batch Translation | Perfect Support | Accurate Source Language |                             Notes                             |
| :--------: | :------------------------: | :-------------: | :----------------------: | :----------------------------------------------------------: |
|    Baidu   |            Yes             |        No       |           No             | Does not support accurately returning the source language type for each result |
|   Google   |            Yes             |       Yes       |          Yes             |                                                              |
|   Youdao   |            Yes             |        No       |           No             |        Source language type recognition is inaccurate         |
|  Volcano   |            Yes             |       Yes       |          Yes             |                                                              |
|   Deepl    |            Yes             |        No       |           No             |        Source language type recognition is inaccurate         |
|   iFLYTEK  |            Yes             |        No       |           No             | Officially does not support batch translation, achieved through special character № slicing, may result in non-multiple results |
|  PaPaGo    |            Yes             |        No       |           No             | Implemented through \n slicing, cannot recognize different source language types |
|  ChatGPT   |            Yes             |       Yes       |          Yes             |                                                              |

## Future Support (Priority Order, Checked Items are Implemented) ✈️
- [x] Persistent storage of translations to `MySQL`
- [x] Web control page
- [x] ChatGPT AI translation
- [x] iFLYTEK translation
- [x] More reasonable and secure authentication
- [x] Tencent translation
- [x] Volcano translation
- [x] PaPaGo
- [x] Support for more national languages
- [x] Support for simulating `LibreTranslate` translation interface
- [x] Support for terminal interactive translation
- [ ] More translation features support for the client

## Basic Types 🪨
`YouDao` `Baidu` `Google` `Deepl` `ChatGPT` `XunFei` `XunFeiNiu` `Tencent` `HuoShan` `PaPaGo`

## Docker Startup 🚀
```shell
# In the project directory
docker build -t uni-translate:latest .
# Then execute (preferably create a network to place mysql and redis in the same network, then directly access the application with the container name in the configuration)
docker run -d --name uniTranslate -v {local_directory}/config.yaml:/app/config.yaml -p 9431:{port_configured_in_config.yaml} --network baseRun uni-translate:latest
```

## Terminal Interaction
After configuring `config.yaml`, execute

```bash
./UniTranslate translate auto en
```

## Configuration Parsing 🗄️

```yaml
server:
  name: uniTranslate
  address: "0.0.0.0:9431"
  cacheMode: redis # redis , mem , off modes. mem stores translation results in program memory, off does not write to any cache
  cachePlatform: false # Whether to include the platform in the cache key generation (affects the key stored during project startup initialization)
  key: "hdasdhasdhsahdkasjfsoufoqjoje" # Secret key for http api integration
  keyMode: 1 # Mode 1 uses key for direct verification, Mode 2 uses key to encrypt and sign data for verification
```

## API Documentation 🌍
[Online Documentation](https://apifox.com/apidoc/shared-335b66b6-90dd-42af-8a1b-f7d1a2c3f351)
[Open Api File](./uniTranslate%20(统一翻译).openapi.json)

## Interface Authentication TS Example
```typescript
import { MD5 } from "crypto-js";

/**
 * 
 * @param key The key set by the platform
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


## Unsupported Translations??? 🤔
All supported languages in this program use the national language **identifiers** based on the _Youdao_ translation API identifiers as specified in the [translate.json](./translate.json) file.

Please refer to the identifiers supported by the _Youdao_ translation API documentation as the basis for modifying the `translate.json` file.