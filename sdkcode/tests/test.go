package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type TestCase struct {
	Name         string                 `json:"name"`
	Params       map[string]interface{} `json:"params"`
	ExpectedHash string                 `json:"expectedHash"`
}

type TestData struct {
	Key       string     `json:"key"`
	TestCases []TestCase `json:"testCases"`
}

// AuthEncrypt 实现
func AuthEncrypt(key string, params map[string]interface{}) string {
	data := SortMapToStr(params)
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SortMapToStr 实现
func SortMapToStr(data map[string]interface{}) string {
	if len(data) == 0 {
		return ""
	}
	
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.Grow(len(data) * 32)

	for i, k := range keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte(':')
		
		switch v := data[k].(type) {
		case map[string]interface{}:
			b.WriteByte('{')
			b.WriteString(SortMapToStr(v))
			b.WriteByte('}')
		case []interface{}:
			b.WriteByte('[')
			for j, item := range v {
				if j > 0 {
					b.WriteByte(',')
				}
				switch iv := item.(type) {
				case map[string]interface{}:
					b.WriteByte('{')
					b.WriteString(SortMapToStr(iv))
					b.WriteByte('}')
				default:
					b.WriteString(fmt.Sprint(iv))
				}
			}
			b.WriteByte(']')
		default:
			b.WriteString(fmt.Sprint(v))
		}
	}
	
	return b.String()
}

func main() {
	// 读取测试数据
	data, err := ioutil.ReadFile("test_data.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var testData TestData
	err = json.Unmarshal(data, &testData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// 执行测试并更新预期值
	results := make([]map[string]interface{}, 0)
	for _, tc := range testData.TestCases {
		hash := AuthEncrypt(testData.Key, tc.Params)
		results = append(results, map[string]interface{}{
			"name": tc.Name,
			"hash": hash,
			"params": tc.Params,
		})
	}

	// 更新测试数据文件
	testData.TestCases = make([]TestCase, len(results))
	for i, r := range results {
		testData.TestCases[i] = TestCase{
			Name:         r["name"].(string),
			Params:       r["params"].(map[string]interface{}),
			ExpectedHash: r["hash"].(string),
		}
	}

	// 写回文件
	updatedData, err := json.MarshalIndent(testData, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = ioutil.WriteFile("test_data.json", updatedData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Test data updated with correct hash values")
	fmt.Println("\nResults:")
	for _, r := range results {
		fmt.Printf("\nTest: %s\nHash: %s\n", r["name"], r["hash"])
	}
}
