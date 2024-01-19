// @ts-ignore
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