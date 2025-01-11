using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Text;

namespace UniTranslate
{
    public class UniTranslateSDK
    {
        /// <summary>
        /// Generate HMAC-SHA256 signature for the given parameters
        /// </summary>
        /// <param name="key">Secret key</param>
        /// <param name="params">Parameters to sign</param>
        /// <returns>Hex-encoded signature</returns>
        public static string AuthEncrypt(string key, Dictionary<string, object> @params)
        {
            string data = SortMapToStr(@params);
            using (var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(key)))
            {
                byte[] hash = hmac.ComputeHash(Encoding.UTF8.GetBytes(data));
                return BitConverter.ToString(hash).Replace("-", "").ToLower();
            }
        }

        /// <summary>
        /// Convert map to sorted string representation
        /// </summary>
        /// <param name="data">Input data</param>
        /// <returns>Sorted string representation</returns>
        private static string SortMapToStr(Dictionary<string, object> data)
        {
            if (data == null || !data.Any())
            {
                return string.Empty;
            }

            var sortedKeys = data.Keys.OrderBy(k => k).ToList();
            var parts = new List<string>();

            foreach (var k in sortedKeys)
            {
                var v = data[k];
                string value;

                if (v is Dictionary<string, object> dict)
                {
                    // Handle nested map
                    value = "{" + SortMapToStr(dict) + "}";
                }
                else if (v is IList<object> list)
                {
                    // Handle array
                    var arrayParts = new List<string>();
                    foreach (var item in list)
                    {
                        if (item is Dictionary<string, object> itemDict)
                        {
                            arrayParts.Add("{" + SortMapToStr(itemDict) + "}");
                        }
                        else
                        {
                            arrayParts.Add(item?.ToString() ?? "");
                        }
                    }
                    value = "[" + string.Join(",", arrayParts) + "]";
                }
                else
                {
                    value = v?.ToString() ?? "";
                }

                parts.Add($"{k}:{value}");
            }

            return string.Join("&", parts);
        }
    }
}
