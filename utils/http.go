package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpPostRequestJson(httpUrl string, jsonData []byte) (string, error) {
	req, err := http.NewRequest("POST", httpUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置请求头，指定内容类型为application/json
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// PostRequest 发送HTTP POST请求
// url: 请求地址
// data: 请求参数(自动转换为JSON)
// headers: 自定义Header头(map类型)
// 返回: 响应内容, 状态码, 错误信息
func HttpPostRequest(url string, data interface{}, headers map[string]string) ([]byte, int, error) {
	// 1. 序列化请求参数
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 2. 创建请求对象
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("创建请求失败: %v", err)
	}

	// 3. 设置默认Header
	req.Header.Set("Content-Type", "application/json")

	// 4. 添加自定义Header
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 5. 配置HTTP客户端（设置超时时间）
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 6. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("读取响应失败: %v", err)
	}

	return body, resp.StatusCode, nil
}

/**
 * @Description: 通用POST请求方法
 * @param targetUrl 请求地址
 * @param headers 请求头(map类型)
 * @param data 请求数据
 * @param result 响应结果结构体指针
 * @param timeout 自定义超时时间（建议至少30秒）
 * @return error 错误信息
 */
func HttpPost(targetUrl string, headers map[string]string, data interface{}, result interface{}, timeout time.Duration) error {
	// 1. 序列化请求参数
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 2. 创建请求对象
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 3. 设置默认Header
	req.Header.Set("Content-Type", "application/json")

	// 4. 添加自定义Header
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 5. 配置HTTP客户端（设置超时时间）
	client := &http.Client{
		Timeout: timeout,
	}

	// 6. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP错误: 状态码 %d, 响应内容: %s",
			resp.StatusCode, string(body))
	}

	// 8. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	// 9. 自动处理JSON反序列化
	if result != nil {
		// 使用decoder来处理可能存在的特殊字符
		decoder := json.NewDecoder(bytes.NewReader(body))
		// decoder.DisallowUnknownFields() // 严格模式
		if err := decoder.Decode(result); err != nil {
			return fmt.Errorf("JSON解析失败: %v", err)
		}
	}

	return nil
}

/**
 * @Description: 通用GET请求方法
 * @param targetUrl 请求地址
 * @param headers 请求头(map类型)
 * @param queryParams URL参数(map类型)
 * @param result 响应结果结构体指针
 * @param timeout 自定义超时时间（建议至少30秒）
 * @return error 错误信息
 */
func HttpGet(targetUrl string, headers map[string]string, queryParams map[string]string, result interface{}, timeout time.Duration) error {
	// 创建带有超时的HTTP客户端
	client := &http.Client{
		Timeout: timeout,
	}

	// 处理URL参数
	params := url.Values{}
	for k, v := range queryParams {
		params.Add(k, v)
	}
	fullUrl := targetUrl
	if len(params) > 0 {
		fullUrl = fmt.Sprintf("%s?%s", targetUrl, params.Encode())
	}

	// 创建请求对象
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 默认设置User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GoHttp/1.1)")
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP错误: 状态码 %d, 响应内容: %s",
			resp.StatusCode, string(body))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	// 自动处理JSON反序列化
	if result != nil {
		// 使用decoder来处理可能存在的特殊字符
		decoder := json.NewDecoder(bytes.NewReader(body))
		// decoder.DisallowUnknownFields() // 严格模式
		if err := decoder.Decode(result); err != nil {
			return fmt.Errorf("JSON解析失败: %v", err)
		}
	}

	return nil
}

/**
 * @Description: 通用POST请求方法（带重试机制）
 * @param targetUrl 请求地址
 * @param headers 请求头(map类型)
 * @param data 请求数据
 * @param result 响应结果结构体指针
 * @param timeout 单次请求的超时时间（建议至少30秒）
 * @param maxRetries 最大重试次数
 * @param retryInterval 每次重试的间隔时间
 * @return error 错误信息
 */
func HttpPostWithRetry(targetUrl string, headers map[string]string, data interface{}, result interface{}, timeout time.Duration, maxRetries int, retryInterval time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 第一次立即请求，后续重试才等待
		if attempt > 0 {
			time.Sleep(retryInterval)
		}

		// 1. 序列化请求参数
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("JSON序列化失败: %v", err)
		}

		// 2. 创建请求对象
		req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("创建请求失败: %v", err)
		}

		// 3. 设置默认Header
		req.Header.Set("Content-Type", "application/json")

		// 4. 添加自定义Header
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// 5. 配置HTTP客户端
		client := &http.Client{
			Timeout: timeout,
		}

		// 6. 发送请求
		resp, err := client.Do(req)
		if err != nil {
			// 判断是否是临时性网络错误或超时
			if isRetryableError(err) {
				lastErr = fmt.Errorf("请求发送失败（第 %d 次重试）: %v", attempt+1, err)
				continue
			}
			return fmt.Errorf("请求发送失败: %v", err)
		}

		defer resp.Body.Close()

		// 7. 检查状态码
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			// 针对 408 状态码进行重试
			if resp.StatusCode == http.StatusRequestTimeout || strings.Contains(string(body), "timeout") {
				lastErr = fmt.Errorf("HTTP错误: 状态码 %d, 响应内容: %s", resp.StatusCode, string(body))
				continue
			}
			return fmt.Errorf("HTTP错误: 状态码 %d, 响应内容: %s", resp.StatusCode, string(body))
		}

		// 8. 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("读取响应体失败: %v", err)
		}

		// 9. 自动处理JSON反序列化
		if result != nil {
			decoder := json.NewDecoder(bytes.NewReader(body))
			if err := decoder.Decode(result); err != nil {
				return fmt.Errorf("JSON解析失败: %v", err)
			}
		}

		// 请求成功，返回 nil
		return nil
	}

	// 所有重试失败
	return fmt.Errorf("请求失败，重试 %d 次仍未成功: %v", maxRetries, lastErr)
}

// 判断错误是否可重试（网络超时、临时性错误等）
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}
	// 可根据需要扩展判断逻辑
	return true
}
