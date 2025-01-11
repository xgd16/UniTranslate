import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.util.*;

public class UniTranslate {
    /**
     * Generate HMAC-SHA256 signature for the given parameters
     * @param key Secret key
     * @param params Parameters to sign
     * @return Hex-encoded signature
     * @throws Exception if encryption fails
     */
    public static String authEncrypt(String key, Map<String, Object> params) throws Exception {
        String data = sortMapToStr(params);
        Mac sha256Hmac = Mac.getInstance("HmacSHA256");
        SecretKeySpec secretKey = new SecretKeySpec(key.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
        sha256Hmac.init(secretKey);
        byte[] hash = sha256Hmac.doFinal(data.getBytes(StandardCharsets.UTF_8));
        return bytesToHex(hash);
    }

    /**
     * Convert map to sorted string representation
     * @param data Input data
     * @return Sorted string representation
     */
    private static String sortMapToStr(Map<String, Object> data) {
        if (data == null || data.isEmpty()) {
            return "";
        }

        List<String> keys = new ArrayList<>(data.keySet());
        Collections.sort(keys);
        StringBuilder result = new StringBuilder();

        for (int i = 0; i < keys.size(); i++) {
            String k = keys.get(i);
            if (i > 0) {
                result.append('&');
            }
            result.append(k).append(':');

            Object v = data.get(k);
            if (v instanceof Map) {
                @SuppressWarnings("unchecked")
                Map<String, Object> map = (Map<String, Object>) v;
                result.append('{').append(sortMapToStr(map)).append('}');
            } else if (v instanceof List) {
                result.append('[');
                List<?> list = (List<?>) v;
                for (int j = 0; j < list.size(); j++) {
                    if (j > 0) {
                        result.append(',');
                    }
                    Object item = list.get(j);
                    if (item instanceof Map) {
                        @SuppressWarnings("unchecked")
                        Map<String, Object> map = (Map<String, Object>) item;
                        result.append('{').append(sortMapToStr(map)).append('}');
                    } else {
                        result.append(String.valueOf(item));
                    }
                }
                result.append(']');
            } else {
                result.append(String.valueOf(v));
            }
        }

        return result.toString();
    }

    private static String bytesToHex(byte[] bytes) {
        StringBuilder result = new StringBuilder();
        for (byte b : bytes) {
            result.append(String.format("%02x", b));
        }
        return result.toString();
    }
}
