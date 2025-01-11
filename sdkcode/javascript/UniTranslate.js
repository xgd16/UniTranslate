const crypto = require('crypto');

class UniTranslate {
    /**
     * Generate HMAC-SHA256 signature for the given parameters
     * @param {string} key Secret key
     * @param {Object} params Parameters to sign
     * @returns {string} Hex-encoded signature
     */
    static authEncrypt(key, params) {
        const data = this.sortMapToStr(params);
        return crypto
            .createHmac('sha256', key)
            .update(data)
            .digest('hex');
    }

    /**
     * Convert map to sorted string representation
     * @param {Object} data Input data
     * @returns {string} Sorted string representation
     */
    static sortMapToStr(data) {
        if (!data || Object.keys(data).length === 0) {
            return '';
        }

        const keys = Object.keys(data).sort();
        return keys.map((k, i) => {
            let value;
            if (Array.isArray(data[k])) {
                // Handle array
                value = '[' + data[k].map(item => {
                    if (item && typeof item === 'object' && !Array.isArray(item)) {
                        return '{' + this.sortMapToStr(item) + '}';
                    }
                    return String(item);
                }).join(',') + ']';
            } else if (data[k] && typeof data[k] === 'object' && !Array.isArray(data[k])) {
                // Handle nested map
                value = '{' + this.sortMapToStr(data[k]) + '}';
            } else {
                value = String(data[k] === null ? '' : data[k]);
            }
            return `${k}:${value}`;
        }).join('&');
    }
}

module.exports = UniTranslate;
