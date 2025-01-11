import * as crypto from 'crypto';

export class UniTranslate {
    /**
     * Generate HMAC-SHA256 signature for the given parameters
     * @param key Secret key
     * @param params Parameters to sign
     * @returns Hex-encoded signature
     */
    static authEncrypt(key: string, params: Record<string, any>): string {
        const data = this.sortMapToStr(params);
        return crypto
            .createHmac('sha256', key)
            .update(data)
            .digest('hex');
    }

    /**
     * Convert map to sorted string representation
     * @param data Input data
     * @returns Sorted string representation
     */
    private static sortMapToStr(data: Record<string, any>): string {
        if (!data || Object.keys(data).length === 0) {
            return '';
        }

        const keys = Object.keys(data).sort();
        return keys.map((k, i) => {
            let value: string;
            if (Array.isArray(data[k])) {
                // Handle array
                value = '[' + data[k].map((item: any) => {
                    if (item && typeof item === 'object' && !Array.isArray(item)) {
                        return '{' + this.sortMapToStr(item) + '}';
                    }
                    return String(item);
                }).join(',') + ']';
            } else if (data[k] && typeof data[k] === 'object') {
                // Handle nested map
                value = '{' + this.sortMapToStr(data[k]) + '}';
            } else {
                value = String(data[k]);
            }
            return `${k}:${value}`;
        }).join('&');
    }
}
