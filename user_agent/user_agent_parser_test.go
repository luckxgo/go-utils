package useragent

import (
	"testing"
)

// TestParseUserAgent 测试用户代理字符串解析功能
func TestParseUserAgent(t *testing.T) {
	// 测试用例集合
	testCases := []struct {
		name      string
		uaStr     string
		expected  *UserAgentInfo
		expectErr bool
	}{
		{
			name:  "Chrome on Windows",
			uaStr: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1",
			expected: &UserAgentInfo{
				OS:             "Windows",
				OSVersion:      "6.1",
				Browser:        "Chrome",
				BrowserVersion: "14.0.835.163",
				Engine:         "AppleWebKit",
				EngineVersion:  "535.1",
				DeviceType:     "desktop",
			},
			expectErr: false,
		},
		{
			name:  "Safari on macOS",
			uaStr: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",
			expected: &UserAgentInfo{
				OS:             "macOS",
				OSVersion:      "10.15.7",
				Browser:        "Safari",
				BrowserVersion: "15.0",
				Engine:         "AppleWebKit",
				EngineVersion:  "605.1.15",
				DeviceType:     "desktop",
			},
			expectErr: false,
		},
		{
			name:  "Firefox on Linux",
			uaStr: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0",
			expected: &UserAgentInfo{
				OS:             "Linux",
				OSVersion:      "",
				Browser:        "Firefox",
				BrowserVersion: "94.0",
				Engine:         "Gecko",
				EngineVersion:  "20100101",
				DeviceType:     "desktop",
			},
			expectErr: false,
		},
		{
			name:  "Mobile Chrome on Android",
			uaStr: "Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Mobile Safari/537.36",
			expected: &UserAgentInfo{
				OS:             "Android",
				OSVersion:      "11",
				Browser:        "Chrome",
				BrowserVersion: "96.0.4664.45",
				Engine:         "AppleWebKit",
				EngineVersion:  "537.36",
				DeviceType:     "mobile",
			},
			expectErr: false,
		},
		{
			name:  "iPad Safari",
			uaStr: "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
			expected: &UserAgentInfo{
				OS:             "iOS",
				OSVersion:      "15.0",
				Browser:        "Safari",
				BrowserVersion: "15.0",
				Engine:         "AppleWebKit",
				EngineVersion:  "605.1.15",
				DeviceType:     "tablet",
			},
			expectErr: false,
		},
		{
			name:      "空UA字符串",
			uaStr:     "",
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseUserAgent(tc.uaStr)

			// 检查错误是否符合预期
			if (err != nil) != tc.expectErr {
				t.Errorf("错误预期: %v, 实际: %v", tc.expectErr, err)
				return
			}

			// 如果预期有错误，不需要继续检查结果
			if tc.expectErr {
				return
			}

			// 检查解析结果是否符合预期
			if result.OS != tc.expected.OS {
				t.Errorf("OS预期: %s, 实际: %s", tc.expected.OS, result.OS)
			}
			if result.OSVersion != tc.expected.OSVersion {
				t.Errorf("OSVersion预期: %s, 实际: %s", tc.expected.OSVersion, result.OSVersion)
			}
			if result.Browser != tc.expected.Browser {
				t.Errorf("Browser预期: %s, 实际: %s", tc.expected.Browser, result.Browser)
			}
			if result.BrowserVersion != tc.expected.BrowserVersion {
				t.Errorf("BrowserVersion预期: %s, 实际: %s", tc.expected.BrowserVersion, result.BrowserVersion)
			}
			if result.Engine != tc.expected.Engine {
				t.Errorf("Engine预期: %s, 实际: %s", tc.expected.Engine, result.Engine)
			}
			if result.EngineVersion != tc.expected.EngineVersion {
				t.Errorf("EngineVersion预期: %s, 实际: %s", tc.expected.EngineVersion, result.EngineVersion)
			}
			if result.DeviceType != tc.expected.DeviceType {
				t.Errorf("DeviceType预期: %s, 实际: %s", tc.expected.DeviceType, result.DeviceType)
			}
		})
	}
}

// BenchmarkParseUserAgent 基准测试解析性能
func BenchmarkParseUserAgent(b *testing.B) {
	uaStr := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseUserAgent(uaStr)
	}
}
