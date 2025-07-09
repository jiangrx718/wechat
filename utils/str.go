package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 定义全局字符集（大小写字母 + 数字）
const charset = "abcdefghijklmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func MD5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateSecureRandomString 生成安全的随机字符串
func GenerateSecureRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// 将随机字节转换为Base64字符串（调整字符集）
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateTimestampJoinSalt(length int) string {
	timestampStr := strconv.FormatInt(GetCurrentUnixTimestamp(), 10)
	saltStr := GenerateRandomString(length)
	timestampStr = saltStr + timestampStr
	return timestampStr
}

// RemoveEmptyStrings 移除字符串切片中的空字符串
func RemoveEmptyStrings(input []string) []string {
	result := make([]string, 0, len(input))
	for _, str := range input {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

// 字符串切片去重去空，并行处理版本 (适用于超大切片)
func DeduplicateParallel(input []string) []string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		seen := make(map[string]struct{})
		for _, s := range input {
			trimmed := strings.TrimSpace(s)
			if trimmed == "" {
				continue
			}
			if _, exists := seen[trimmed]; !exists {
				seen[trimmed] = struct{}{}
				ch <- trimmed
			}
		}
	}()

	result := make([]string, 0)
	for s := range ch {
		result = append(result, s)
	}
	return result
}

// DifferenceSlice 计算两个字符串切片的差集
func DifferenceSlice(slice1, slice2 []string) []string {
	set := make(map[string]struct{})
	for _, s := range slice2 {
		set[s] = struct{}{}
	}

	var diff []string
	for _, s := range slice1 {
		if _, found := set[s]; !found {
			diff = append(diff, s)
		}
	}

	return diff
}
