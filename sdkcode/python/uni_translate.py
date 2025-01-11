import hmac
import hashlib
from typing import Dict, Any, Union

class UniTranslate:
    @staticmethod
    def auth_encrypt(key: str, params: Dict[str, Any]) -> str:
        """
        Generate HMAC-SHA256 signature for the given parameters
        Args:
            key: Secret key
            params: Parameters to sign
        Returns:
            Hex-encoded signature
        """
        data = UniTranslate.sort_map_to_str(params)
        hmac_obj = hmac.new(
            key.encode('utf-8'),
            data.encode('utf-8'),
            hashlib.sha256
        )
        return hmac_obj.hexdigest()

    @staticmethod
    def sort_map_to_str(data: Dict[str, Any]) -> str:
        """
        Convert map to sorted string representation
        Args:
            data: Input data
        Returns:
            Sorted string representation
        """
        if not data:
            return ''

        parts = []
        for k in sorted(data.keys()):
            v = data[k]
            if isinstance(v, dict):
                # Handle nested map
                value = '{' + UniTranslate.sort_map_to_str(v) + '}'
            elif isinstance(v, (list, tuple)):
                # Handle array
                array_parts = []
                for item in v:
                    if isinstance(item, dict):
                        array_parts.append('{' + UniTranslate.sort_map_to_str(item) + '}')
                    else:
                        array_parts.append(str(item if item is not None else ''))
                value = '[' + ','.join(array_parts) + ']'
            else:
                value = str(v if v is not None else '')
            parts.append(f"{k}:{value}")

        return '&'.join(parts)
