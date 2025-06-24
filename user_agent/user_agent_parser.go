package useragent

import (
	"errors"
	"regexp"
	"strings"
)

// UserAgentInfo 存储解析后的用户代理信息
type UserAgentInfo struct {
	OS             string // 操作系统名称
	OSVersion      string // 操作系统版本
	Browser        string // 浏览器名称
	BrowserVersion string // 浏览器版本
	Engine         string // 渲染引擎名称
	EngineVersion  string // 渲染引擎版本
	DeviceType     string // 设备类型(desktop/mobile/tablet/other)
}

// 定义解析规则结构体
type parseRule struct {
	regexp          *regexp.Regexp
	osIndex         int
	osVerIndex      int
	browserIndex    int
	browserVerIndex int
	engineIndex     int
	engineVerIndex  int
	browserName     string // 显式指定浏览器名称
}

// 解析规则列表 - 按优先级排序
var parseRules = []parseRule{
	// Chrome/Chromium 规则
	{
		regexp:          regexp.MustCompile(`Chrome/([\d.]+)`), // 支持任意格式的版本号
		browserIndex:    -1,                                    // 不使用捕获组作为浏览器名称
		browserVerIndex: 1,
		osIndex:         -1,
		osVerIndex:      -1,
		engineIndex:     -1,
		engineVerIndex:  -1,
		browserName:     "Chrome", // 显式指定浏览器名称
	},
	// Safari 版本规则 (优先匹配Version字段)
	{
		regexp:          regexp.MustCompile(`Version/([\d.]+)`), // 支持多段版本号
		browserIndex:    -1,                                     // 不使用捕获组作为浏览器名称
		browserVerIndex: 1,
		osIndex:         -1,
		osVerIndex:      -1,
		engineIndex:     -1,
		engineVerIndex:  -1,
		browserName:     "Safari", // 显式指定为Safari浏览器
	},
	// Safari 后备规则
	{
		regexp:          regexp.MustCompile(`Safari/(\d+\.\d+)`),
		browserIndex:    0,
		browserVerIndex: 1,
		osIndex:         -1,
		osVerIndex:      -1,
		engineIndex:     -1,
		engineVerIndex:  -1,
	},
	// Firefox 规则
	{
		regexp:          regexp.MustCompile(`Firefox/([\d.]+)`),
		browserIndex:    -1,
		browserVerIndex: 1,
		osIndex:         -1,
		osVerIndex:      -1,
		engineIndex:     -1,
		engineVerIndex:  -1,
		browserName:     "Firefox",
	},
	// Edge 规则
	{
		regexp:          regexp.MustCompile(`Edge/(\d+\.\d+\.\d+)`),
		browserIndex:    0,
		browserVerIndex: 1,
		osIndex:         -1,
		osVerIndex:      -1,
		engineIndex:     -1,
		engineVerIndex:  -1,
	},
	// AppleWebKit 引擎规则
	{
		regexp:          regexp.MustCompile(`(AppleWebKit)/([\d.]+)`),
		engineIndex:     1,
		engineVerIndex:  2,
		osIndex:         -1,
		osVerIndex:      -1,
		browserIndex:    -1,
		browserVerIndex: -1,
	},
	// Gecko 引擎规则
	{
		regexp:          regexp.MustCompile(`(Gecko)/(\d+)`), // 支持纯数字版本号
		engineIndex:     1,
		engineVerIndex:  2,
		osIndex:         -1,
		osVerIndex:      -1,
		browserIndex:    -1,
		browserVerIndex: -1,
	},
}

// 操作系统匹配规则
var osRules = []struct {
	regexp       *regexp.Regexp
	osName       string
	versionIndex int
}{
	{regexp.MustCompile(`Windows NT (\d+\.\d+)`), "Windows", 1},
	{regexp.MustCompile(`Mac OS X (\d+_\d+_\d+)`), "macOS", 1},
	{regexp.MustCompile(`Android (\d+(\.\d+)*)`), "Android", 1},
	{regexp.MustCompile(`iPad; CPU OS (\d+_\d+)`), "iOS", 1},
	{regexp.MustCompile(`iOS (\d+\.\d+)`), "iOS", 1},
	{regexp.MustCompile(`Linux`), "Linux", -1},
}

// ParseUserAgent 解析用户代理字符串并返回结构化信息
// uaStr: 用户代理字符串
// 返回解析后的信息和可能的错误
func ParseUserAgent(uaStr string) (*UserAgentInfo, error) {
	if uaStr == "" {
		return nil, errors.New("用户代理字符串不能为空")
	}

	info := &UserAgentInfo{}

	// 解析操作系统信息
	info.OS, info.OSVersion = parseOS(uaStr)

	// 解析渲染引擎
	info.Engine, info.EngineVersion = parseEngine(uaStr)

	// 解析浏览器
	info.Browser, info.BrowserVersion = parseBrowser(uaStr)

	// 确定设备类型
	info.DeviceType = determineDeviceType(uaStr, info.OS)

	return info, nil
}

// parseOS 解析操作系统信息
func parseOS(uaStr string) (osName, osVersion string) {
	for _, rule := range osRules {
		matches := rule.regexp.FindStringSubmatch(uaStr)
		if len(matches) > 0 {
			osName = rule.osName
			if rule.versionIndex > 0 && len(matches) > rule.versionIndex {
				osVersion = matches[rule.versionIndex]
				// 处理macOS版本格式(将下划线转为点)
				if osName == "macOS" || osName == "iOS" {
					osVersion = strings.ReplaceAll(osVersion, "_", ".")
				}
			}
			return
		}
	}
	return "Unknown", ""
}

// parseEngine 解析渲染引擎信息
func parseEngine(uaStr string) (engineName, engineVersion string) {
	for _, rule := range parseRules {
		if rule.engineIndex >= 0 {
			matches := rule.regexp.FindStringSubmatch(uaStr)
			if len(matches) > 0 {
				// 从捕获组提取引擎名称
				if rule.engineIndex < len(matches) {
					engineName = matches[rule.engineIndex]
				} else {
					engineName = "Unknown"
				}
				if rule.engineVerIndex > 0 && len(matches) > rule.engineVerIndex {
					engineVersion = matches[rule.engineVerIndex]
				}
				return
			}
		}
	}
	return "Unknown", ""
}

// parseBrowser 解析浏览器信息
func parseBrowser(uaStr string) (browserName, browserVersion string) {
	for _, rule := range parseRules {
		// 处理具有浏览器版本索引的规则（无论是否有显式浏览器名称）
		if rule.browserVerIndex >= 0 {
			matches := rule.regexp.FindStringSubmatch(uaStr)
			if len(matches) > 0 {
				// 优先使用规则中定义的显式浏览器名称
				if rule.browserName != "" {
					browserName = rule.browserName
				} else if rule.browserIndex >= 0 && rule.browserIndex < len(matches) {
					browserName = matches[rule.browserIndex]
				} else {
					// 根据正则表达式确定浏览器名称
					switch {
					case strings.Contains(rule.regexp.String(), "Chrome"):
						browserName = "Chrome"
					case strings.Contains(rule.regexp.String(), "Safari") && !strings.Contains(uaStr, "Chrome"):
						browserName = "Safari"
					case strings.Contains(rule.regexp.String(), "Firefox"):
						browserName = "Firefox"
					case strings.Contains(rule.regexp.String(), "Edge"):
						browserName = "Edge"
					default:
						browserName = "Unknown"
					}
				}
				if rule.browserVerIndex < len(matches) {
					browserVersion = matches[rule.browserVerIndex]
				}
				return
			}
		}
	}
	return "Unknown", ""
}

// determineDeviceType 确定设备类型
func determineDeviceType(uaStr, osName string) string {
	lowerUA := strings.ToLower(uaStr)
	// 优先检测平板设备
	if strings.Contains(lowerUA, "tablet") || (osName == "iOS" && strings.Contains(lowerUA, "ipad")) {
		return "tablet"
	} else if strings.Contains(lowerUA, "mobile") || (osName == "Android" && !strings.Contains(lowerUA, "tablet")) {
		// 检测移动设备
		return "mobile"
	} else if osName == "Windows" || osName == "macOS" || osName == "Linux" {
		// 桌面设备
		return "desktop"
	}
	return "other"
}
