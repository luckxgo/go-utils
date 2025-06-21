package strutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty 判断字符串是否为空（长度为0）
// 参数:
//
//	s - 待检查的字符串
//
// 返回值:
//
//	如果字符串长度为0则返回true，否则返回false
//
// 示例:
//
//	IsEmpty("") → true
//	IsEmpty("hello") → false
func IsEmpty(s string) bool {
	return len(s) == 0
}

// UrlEncode 对字符串进行URL编码(Percent-Encoding)
// 参数:
//
//	s - 待编码的字符串
//
// 返回值:
//
//	URL编码后的字符串
//
// 示例:
//
//	UrlEncode("hello world!") → "hello+world%21"
func UrlEncode(s string) string {
	return url.QueryEscape(s)
}

// UrlDecode 对URL编码的字符串进行解码
// 参数:
//
//	s - URL编码的字符串
//
// 返回值:
//
//	解码后的原始字符串和可能的错误
//
// 示例:
//
//	UrlDecode("hello+world%21") → "hello world!", nil
func UrlDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

// WordCount 统计单词数量，支持自定义分隔符
// 参数:
//
//	s - 待统计的字符串
//	separators - 可变参数，自定义分隔符集合
//
// 返回值:
//
//	单词数量
//
// 示例:
//
//	WordCount("hello,world,go", ',') → 3
//	WordCount("hello world\tgo", ' ', '\t') → 3
func WordCount(s string, separators ...rune) int {
	if IsEmpty(s) {
		return 0
	}

	sepSet := make(map[rune]bool)
	for _, sep := range separators {
		sepSet[sep] = true
	}

	// 如果没有提供分隔符，默认使用空格
	if len(sepSet) == 0 {
		sepSet[' '] = true
	}

	count := 0
	inWord := false

	for _, c := range s {
		if sepSet[c] {
			if inWord {
				count++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	// 处理最后一个单词
	if inWord {
		count++
	}

	return count
}

// CharCount 统计字符出现频率
// 参数:
//
//	s - 待统计的字符串
//	caseSensitive - 是否区分大小写，true区分，false不区分
//
// 返回值:
//
//	字符到出现次数的映射
//
// 示例:
//
//	CharCount("Hello", true) → map[rune]int{'H':1, 'e':1, 'l':2, 'o':1}
//	CharCount("Hello", false) → map[rune]int{'h':1, 'e':1, 'l':2, 'o':1}
func CharCount(s string, caseSensitive bool) map[rune]int {
	counts := make(map[rune]int)
	for _, c := range s {
		if !caseSensitive {
			c = unicode.ToLower(c)
		}
		counts[c]++
	}
	return counts
}

// IsPalindrome 判断字符串是否为回文（正读反读一致）
// 参数:
//
//	s - 待检查的字符串
//
// 返回值:
//
//	如果是回文则返回true，否则返回false
//
// 示例:
//
//	IsPalindrome("abba") → true
//	IsPalindrome("abc") → false
func IsPalindrome(s string) bool {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// RegexMatch 检查字符串是否匹配正则表达式
// 参数:
//
//	pattern - 正则表达式模式
//	s - 待检查的字符串
//
// 返回值:
//
//	是否匹配和可能的错误
//
// 示例:
//
//	RegexMatch(`^[a-zA-Z0-9]+$`, "hello123") → true, nil
func RegexMatch(pattern, s string) (bool, error) {
	return regexp.MatchString(pattern, s)
}

// RegexExtract 从字符串中提取第一个匹配的子串
// 参数:
//
//	pattern - 正则表达式模式
//	s - 待提取的字符串
//
// 返回值:
//
//	第一个匹配的子串和可能的错误
//
// 示例:
//
//	RegexExtract(`\d+`, "age: 25, score: 90") → "25", nil
func RegexExtract(pattern, s string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	match := re.FindString(s)
	return match, nil
}

// RegexReplaceAll 使用正则表达式替换所有匹配项
// 参数:
//
//	pattern - 正则表达式模式
//	s - 待替换的字符串
//	replacement - 替换字符串
//
// 返回值:
//
//	替换后的字符串和可能的错误
//
// 示例:
//
//	RegexReplaceAll(`\d+`, "a1b2c3", "x") → "axbxcx", nil
func RegexReplaceAll(pattern, s, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	return re.ReplaceAllString(s, replacement), nil
}

// IsBlank 判断字符串是否为空白（全是空白字符或空字符串）
// 空白字符包括空格、制表符、换行符等unicode.IsSpace认定的字符
// 参数:
//
//	s - 待检查的字符串
//
// 返回值:
//
//	如果字符串为空或仅包含空白字符则返回true，否则返回false
//
// 示例:
//
//	IsBlank("") → true
//	IsBlank("  	") → true
//	IsBlank(" hello ") → false
func IsBlank(s string) bool {
	if IsEmpty(s) {
		return true
	}
	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}

// Trim 去除字符串两端的空白字符
// 使用unicode.IsSpace定义的空白字符，包括空格、制表符、换行符等
// 参数:
//
//	s - 待处理的字符串
//
// 返回值:
//
//	去除两端空白后的新字符串
//
// 示例:
//
//	Trim("  hello world  ") → "hello world"
//	Trim("\thello\n") → "hello"
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimAll 去除字符串中所有空白字符
func TrimAll(s string) string {
	var builder strings.Builder
	builder.Grow(len(s)) // 预分配空间，避免多次扩容
	for _, c := range s {
		if !unicode.IsSpace(c) {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

// Substring 安全地截取子字符串，支持负索引（从末尾开始计数）
// start 起始索引（包含），end 结束索引（不包含），如果为负数则从末尾开始计算
func Substring(s string, start, end int) (string, error) {
	runes := []rune(s)
	length := len(runes)

	// 处理负索引
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	// 边界检查
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		return "", fmt.Errorf("start index %d greater than end index %d", start, end)
	}

	return string(runes[start:end]), nil
}

// Split 分割字符串，支持多个分隔符，忽略空字符串结果
// 参数:
//
//	s - 待分割的字符串
//	separators - 可变参数，分隔符集合
//
// 返回值:
//
//	分割后的字符串切片，不包含空字符串
//
// 示例:
//
//	Split("a,b;c", ',', ';') → ["a", "b", "c"]
//	Split("a,,b", ',') → ["a", "b"]
func Split(s string, separators ...rune) []string {
	if IsEmpty(s) {
		return []string{}
	}

	// 创建分隔符集合
	sepSet := make(map[rune]bool)
	for _, sep := range separators {
		sepSet[sep] = true
	}

	var result []string
	var builder strings.Builder

	for _, c := range s {
		if sepSet[c] {
			if builder.Len() > 0 {
				result = append(result, builder.String())
				builder.Reset()
			}
		} else {
			builder.WriteRune(c)
		}
	}

	// 添加最后一个元素
	if builder.Len() > 0 {
		result = append(result, builder.String())
	}

	return result
}

// Join 连接字符串数组，使用指定的分隔符
// 参数:
//
//	strs - 字符串切片
//	sep - 分隔符
//
// 返回值:
//
//	用分隔符连接后的字符串
//
// 示例:
//
//	Join([]string{"a", "b", "c"}, ",") → "a,b,c"
//	Join([]string{}, ",") → ""
func Join(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// Equals 安全比较两个字符串
func Equals(a, b string) bool {
	return a == b
}

// DefaultIfEmpty 如果字符串为空则返回默认值
func DefaultIfEmpty(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsNotBlank 判断字符串是否非空白（包含非空白字符）
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// IsNotEmpty 判断字符串是否非空（长度大于0）
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Capitalize 首字母大写，其余字母小写
func Capitalize(s string) string {
	if IsEmpty(s) {
		return s
	}
	runes := []rune(s)
	// 首字母大写
	runes[0] = unicode.ToUpper(runes[0])
	// 其余字母小写
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// Uncapitalize 首字母小写，其余字母不变
func Uncapitalize(s string) string {
	if IsEmpty(s) {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// PadLeft 使用指定字符在左侧填充字符串至指定长度
func PadLeft(s string, length int, padChar rune) string {
	return pad(s, length, padChar, true)
}

// PadRight 使用指定字符在右侧填充字符串至指定长度
func PadRight(s string, length int, padChar rune) string {
	return pad(s, length, padChar, false)
}

// pad 内部填充实现
func pad(s string, length int, padChar rune, isLeft bool) string {
	runes := []rune(s)
	currentLen := len(runes)
	if currentLen >= length {
		return s
	}

	padLen := length - currentLen
	padRunes := make([]rune, padLen)
	for i := 0; i < padLen; i++ {
		padRunes[i] = padChar
	}

	if isLeft {
		return string(append(padRunes, runes...))
	} else {
		return string(append(runes, padRunes...))
	}
}

// RemovePrefix 移除字符串前缀
func RemovePrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// RemoveSuffix 移除字符串后缀
func RemoveSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// Replace 替换字符串中的所有匹配项
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Count 计算子字符串出现的次数
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// IsAllBlank 检查字符串数组中的所有元素是否都是空白字符串
// 如果数组为空或所有元素都是空白字符串(空字符串、null或仅包含空白字符)，返回true
// 示例:
//
//	IsAllBlank()                // true
//	IsAllBlank("", " ", "\t")   // true
//	IsAllBlank("123", " ")      // false
//	IsAllBlank("123", "abc")    // false
func IsAllBlank(strs ...string) bool {
	if len(strs) == 0 {
		return true
	}

	for _, s := range strs {
		if !IsBlank(s) {
			return false
		}
	}
	return true
}

// IsAllEmpty 检查字符串数组中的所有元素是否都是空字符串
// 如果数组为空或所有元素都是空字符串(空字符串或null)，返回true
// 示例:
//
//	IsAllEmpty()                // true
//	IsAllEmpty("", "")          // true
//	IsAllEmpty("123", "")       // false
//	IsAllEmpty("123", "abc")    // false
//	IsAllEmpty(" ", "\t")       // false
func IsAllEmpty(strs ...string) bool {
	if len(strs) == 0 {
		return true
	}

	for _, s := range strs {
		if !IsEmpty(s) {
			return false
		}
	}
	return true
}

// ToUpper 将字符串转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower 将字符串转换为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToCamelCase 将字符串转换为驼峰命名法
func ToCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// ToSnakeCase 将字符串转换为蛇形命名法
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// IsNumeric 检查字符串是否只包含ASCII数字字符(0-9)
// 注意: 此函数仅支持ASCII数字，不支持 Unicode 数字字符（如 ½、③等）
// 参数:
//
//	s - 待检查的字符串
//
// 返回值:
//
//	如果字符串非空且所有字符都是0-9的数字则返回true，否则返回false
//
// 示例:
//
//	IsNumeric("12345") → true
//	IsNumeric("12.34") → false (包含小数点)
//	IsNumeric("123a5") → false (包含字母)
//	IsNumeric("") → false (空字符串)
func IsNumeric(s string) bool {
	if IsEmpty(s) {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// IsAlpha 检查字符串是否只包含ASCII字母字符(a-z,A-Z)
// 注意: 此函数仅支持ASCII字母，不支持 Unicode 字母字符（如 é、ñ、ü等）
// 参数:
//
//	s - 待检查的字符串
//
// 返回值:
//
//	如果字符串非空且所有字符都是字母则返回true，否则返回false
//
// 示例:
//
//	IsAlpha("abcXYZ") → true
//	IsAlpha("abc123") → false (包含数字)
//	IsAlpha("") → false (空字符串)
func IsAlpha(s string) bool {
	if IsEmpty(s) {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}
	return true
}

// IsAlphanumeric 检查字符串是否只包含ASCII字母和数字字符
// 注意: 此函数仅支持ASCII字符，不支持 Unicode 字符
// 参数:
//
//	s - 待检查的字符串
//
// 返回值:
//
//	如果字符串非空且所有字符都是字母或数字则返回true，否则返回false
//
// 示例:
//
//	IsAlphanumeric("abc123") → true
//	IsAlphanumeric("abc_123") → false (包含下划线)
//	IsAlphanumeric("") → false (空字符串)
func IsAlphanumeric(s string) bool {
	if IsEmpty(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// Base64Encode 将字符串编码为base64
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64Decode 解码base64字符串
func Base64Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Mask 对字符串中的敏感信息进行掩码处理
// 示例: Mask("13812345678", 3, 4, '*') → "138****5678"
func Mask(s string, leftUnmaskLen, rightUnmaskLen int, maskChar rune) string {
	runes := []rune(s)
	if len(runes) <= leftUnmaskLen+rightUnmaskLen {
		return s
	}
	// Calculate total mask length
	maskLen := len(runes) - leftUnmaskLen - rightUnmaskLen
	if maskLen <= 0 {
		return s
	}
	// Replace with exact number of mask characters
	masked := string(runes[:leftUnmaskLen]) + strings.Repeat(string(maskChar), maskLen) + string(runes[len(runes)-rightUnmaskLen:])
	return masked
}

// RandomUUID 生成随机UUID (Version 4) 字符串
// 采用RFC 4122标准，格式为8-4-4-4-12的十六进制字符
// 返回值:
//
//	生成的UUID字符串
//
// 注意:
//
//	如果随机数生成失败会触发panic
//
// 示例:
//
//	RandomUUID() → "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
func RandomUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// Format 格式化字符串模板，用参数替换模板中的占位符{}
// 参数说明:
//
//	template - 包含占位符{}的模板字符串
//	params   - 可变参数，用于替换模板中的占位符
//
// 返回值:
//
//	格式化后的字符串
//
// 示例:
//
//	Format("Hello, {}!", "World") => "Hello, World!"
//	Format("Name: {}, Age: {}", "Alice") => "Name: Alice, Age: {}"
func Format(template string, params ...string) string {
	var result strings.Builder
	paramIndex := 0
	placeholderStart := -1

	for i, c := range template {
		if c == '{' && placeholderStart == -1 {
			placeholderStart = i
		} else if c == '}' && placeholderStart != -1 {
			if paramIndex < len(params) {
				result.WriteString(params[paramIndex])
				paramIndex++
			} else {
				result.WriteString(template[placeholderStart : i+1])
			}
			placeholderStart = -1
		} else if placeholderStart == -1 {
			result.WriteRune(c)
		}
	}

	// Handle unclosed placeholder at end of string
	if placeholderStart != -1 {
		result.WriteString(template[placeholderStart:])
	}

	return result.String()
}
