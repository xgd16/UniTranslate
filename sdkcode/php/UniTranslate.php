<?php

class UniTranslate {
    /**
     * Generate HMAC-SHA256 signature for the given parameters
     * @param string $key Secret key
     * @param array $params Parameters to sign
     * @return string Hex-encoded signature
     */
    public static function authEncrypt(string $key, array $params): string {
        $data = self::sortMapToStr($params);
        return hash_hmac('sha256', $data, $key);
    }

    /**
     * Convert map to sorted string representation
     * @param array $data Input data
     * @return string Sorted string representation
     */
    private static function sortMapToStr(array $data): string {
        if (empty($data)) {
            return '';
        }

        ksort($data);
        $parts = [];

        foreach ($data as $k => $v) {
            if (is_array($v)) {
                if (self::isAssoc($v)) {
                    // Handle nested map
                    $value = '{' . self::sortMapToStr($v) . '}';
                } else {
                    // Handle array
                    $arrayParts = [];
                    foreach ($v as $item) {
                        if (is_array($item) && self::isAssoc($item)) {
                            $arrayParts[] = '{' . self::sortMapToStr($item) . '}';
                        } else {
                            $arrayParts[] = is_null($item) ? '' : (string)$item;
                        }
                    }
                    $value = '[' . implode(',', $arrayParts) . ']';
                }
            } else {
                $value = is_null($v) ? '' : (string)$v;
            }
            $parts[] = $k . ':' . $value;
        }

        return implode('&', $parts);
    }

    /**
     * Check if array is associative
     * @param array $arr Array to check
     * @return bool True if associative
     */
    private static function isAssoc(array $arr): bool {
        if (empty($arr)) return false;
        return array_keys($arr) !== range(0, count($arr) - 1);
    }
}
