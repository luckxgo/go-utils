package strutil

import (
	"testing"
)

func TestUrlEncode(t *testing.T) {
	cases := []struct {
		name        string
		s           string
		wantEncoded string
	}{{
		name:        "normal string",
		s:           "hello world",
		wantEncoded: "hello+world",
	}, {
		name:        "special characters",
		s:           "&=?#% ",
		wantEncoded: "%26%3D%3F%23%25+",
	}, {
		name:        "empty string",
		s:           "",
		wantEncoded: "",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := UrlEncode(tc.s)
			if got != tc.wantEncoded {
				t.Errorf("UrlEncode(%q) = %q, want %q", tc.s, got, tc.wantEncoded)
			}
		})
	}
}

func TestUrlDecode(t *testing.T) {
	cases := []struct {
		name        string
		s           string
		wantDecoded string
	}{{
		name:        "normal string",
		s:           "hello+world",
		wantDecoded: "hello world",
	}, {
		name:        "special characters",
		s:           "%26%3D%3F%23%25+",
		wantDecoded: "&=?#% ",
	}, {
		name:        "empty string",
		s:           "",
		wantDecoded: "",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := UrlDecode(tc.s)
			if err != nil {
				t.Fatalf("UrlDecode(%q) returned error: %v", tc.s, err)
			}
			if got != tc.wantDecoded {
				t.Errorf("UrlDecode(%q) = %q, want %q", tc.s, got, tc.wantDecoded)
			}
		})
	}
}

func TestRegexReplaceAll(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		repl    string
		s       string
		want    string
	}{{
		name:    "replace numbers",
		pattern: `\d+`,
		repl:    "X",
		s:       "abc123def456",
		want:    "abcXdefX",
	}, {
		name:    "no match",
		pattern: `\d+`,
		repl:    "X",
		s:       "abcdef",
		want:    "abcdef",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RegexReplaceAll(tc.pattern, tc.s, tc.repl)
			if err != nil {
				t.Fatalf("RegexReplaceAll(%q, %q, %q) returned error: %v", tc.pattern, tc.repl, tc.s, err)
			}
			if got != tc.want {
				t.Errorf("RegexReplaceAll(%q, %q, %q) = %q, want %q", tc.pattern, tc.repl, tc.s, got, tc.want)
			}
		})
	}
}

func TestUrlEncodeDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantEncoded string
	}{{
		name:        "normal",
		input:       "hello world!",
		wantEncoded: "hello+world%21",
	}, {
		name:        "special_chars",
		input:       "&=?#% ",
		wantEncoded: "%26%3D%3F%23%25+",
	}, {
		name:        "chinese",
		input:       "你好世界",
		wantEncoded: "%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := UrlEncode(tt.input)
			if encoded != tt.wantEncoded {
				t.Errorf("UrlEncode() = %v, want %v", encoded, tt.wantEncoded)
			}

			decoded, err := UrlDecode(encoded)
			if err != nil {
				t.Fatalf("UrlDecode() error = %v", err)
			}
			if decoded != tt.input {
				t.Errorf("UrlDecode() = %v, want %v", decoded, tt.input)
			}
		})
	}

	// 测试解码错误情况
	t.Run("invalid_encoding", func(t *testing.T) {
		_, err := UrlDecode("%xx")
		if err == nil {
			t.Error("UrlDecode() expected error for invalid encoding, got nil")
		}
	})
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		separators []rune
		want       int
	}{{
		name:       "default_separator",
		input:      "hello world go",
		separators: nil,
		want:       3,
	}, {
		name:       "custom_separator",
		input:      "hello,world,go",
		separators: []rune{','},
		want:       3,
	}, {
		name:       "multiple_separators",
		input:      "hello world\tgo\nrust",
		separators: []rune{' ', '\t', '\n'},
		want:       4,
	}, {
		name:       "empty_string",
		input:      "",
		separators: nil,
		want:       0,
	}, {
		name:       "single_word",
		input:      "hello",
		separators: []rune{','},
		want:       1,
	}, {
		name:       "leading_trailing_separators",
		input:      ",hello,world,",
		separators: []rune{','},
		want:       2,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordCount(tt.input, tt.separators...); got != tt.want {
				t.Errorf("WordCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharCount(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		caseSensitive bool
		want          map[rune]int
	}{{
		name:          "case_sensitive",
		input:         "Hello",
		caseSensitive: true,
		want:          map[rune]int{'H': 1, 'e': 1, 'l': 2, 'o': 1},
	}, {
		name:          "case_insensitive",
		input:         "Hello",
		caseSensitive: false,
		want:          map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1},
	}, {
		name:          "empty_string",
		input:         "",
		caseSensitive: true,
		want:          map[rune]int{},
	}, {
		name:          "unicode_chars",
		input:         "你好世界",
		caseSensitive: true,
		want:          map[rune]int{'你': 1, '好': 1, '世': 1, '界': 1},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CharCount(tt.input, tt.caseSensitive)
			if len(got) != len(tt.want) {
				t.Fatalf("CharCount() length = %v, want %v", len(got), len(tt.want))
			}

			for char, count := range tt.want {
				if got[char] != count {
					t.Errorf("CharCount() count for '%c' = %v, want %v", char, got[char], count)
				}
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{{
		name:  "palindrome_even",
		input: "abba",
		want:  true,
	}, {
		name:  "palindrome_odd",
		input: "abcba",
		want:  true,
	}, {
		name:  "not_palindrome",
		input: "hello",
		want:  false,
	}, {
		name:  "empty_string",
		input: "",
		want:  true,
	}, {
		name:  "single_char",
		input: "a",
		want:  true,
	}, {
		name:  "unicode_palindrome",
		input: "上海自来水来自海上",
		want:  true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.input); got != tt.want {
				t.Errorf("IsPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
		wantErr bool
	}{{
		name:    "match",
		pattern: `^[a-zA-Z0-9]+$`,
		input:   "hello123",
		want:    true,
		wantErr: false,
	}, {
		name:    "no_match",
		pattern: `^[0-9]+$`,
		input:   "abc123",
		want:    false,
		wantErr: false,
	}, {
		name:    "invalid_pattern",
		pattern: `[a-z`,
		input:   "test",
		want:    false,
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexMatch(tt.pattern, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegexMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexExtract(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    string
		wantErr bool
	}{{
		name:    "extract_number",
		pattern: `\d+`,
		input:   "age: 25, score: 90",
		want:    "25",
		wantErr: false,
	}, {
		name:    "no_match",
		pattern: `\d+`,
		input:   "no numbers here",
		want:    "",
		wantErr: false,
	}, {
		name:    "invalid_pattern",
		pattern: `[a-z`,
		input:   "test",
		want:    "",
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexExtract(tt.pattern, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexExtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegexExtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{{
		name: "trim whitespace",
		s:    "  hello world  ",
		want: "hello world",
	}, {
		name: "no whitespace to trim",
		s:    "hello",
		want: "hello",
	}, {
		name: "blank string",
		s:    "   \t\n",
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.s); got != tt.want {
				t.Errorf("Trim(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestTrimAll(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{{
		name: "remove all whitespace",
		s:    "h e l l o",
		want: "hello",
	}, {
		name: "mixed whitespace",
		s:    "h\te\tl\tl\to",
		want: "hello",
	}, {
		name: "no whitespace",
		s:    "hello",
		want: "hello",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimAll(tt.s); got != tt.want {
				t.Errorf("TrimAll(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		start   int
		end     int
		want    string
		wantErr bool
	}{{
		name:    "normal case",
		s:       "hello world",
		start:   0,
		end:     5,
		want:    "hello",
		wantErr: false,
	}, {
		name:    "negative indices",
		s:       "hello world",
		start:   -5,
		end:     -1,
		want:    "worl",
		wantErr: false,
	}, {
		name:    "start greater than end",
		s:       "hello",
		start:   3,
		end:     2,
		want:    "",
		wantErr: true,
	}, {
		name:    "out of bounds",
		s:       "hello",
		start:   3,
		end:     10,
		want:    "lo",
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Substring(tt.s, tt.start, tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Substring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Substring() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		separators []rune
		want       []string
	}{{
		name:       "split by comma",
		s:          "a,b,c",
		separators: []rune{','},
		want:       []string{"a", "b", "c"},
	}, {
		name:       "multiple separators",
		s:          "a,b;c",
		separators: []rune{',', ';'},
		want:       []string{"a", "b", "c"},
	}, {
		name:       "ignore empty results",
		s:          "a,,b",
		separators: []rune{','},
		want:       []string{"a", "b"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.s, tt.separators...); !equalStringSlices(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		sep  string
		want string
	}{{
		name: "join with comma",
		strs: []string{"a", "b", "c"},
		sep:  ",",
		want: "a,b,c",
	}, {
		name: "empty slice",
		strs: []string{},
		sep:  ",",
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.strs, tt.sep); got != tt.want {
				t.Errorf("Join() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{{
		name: "equal strings",
		a:    "hello",
		b:    "hello",
		want: true,
	}, {
		name: "not equal strings",
		a:    "hello",
		b:    "world",
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equals(tt.a, tt.b); got != tt.want {
				t.Errorf("Equals(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDefaultIfEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		def  string
		want string
	}{{
		name: "empty string",
		s:    "",
		def:  "default",
		want: "default",
	}, {
		name: "non-empty string",
		s:    "value",
		def:  "default",
		want: "value",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultIfEmpty(tt.s, tt.def); got != tt.want {
				t.Errorf("DefaultIfEmpty(%q, %q) = %q, want %q", tt.s, tt.def, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{{
		name: "reverse string",
		s:    "hello",
		want: "olleh",
	}, {
		name: "empty string",
		s:    "",
		want: "",
	}, {
		name: "unicode characters",
		s:    "世界",
		want: "界世",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.s); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsNotBlank(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{{
		name: "blank string",
		s:    "   ",
		want: false,
	}, {
		name: "non-blank string",
		s:    "test",
		want: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotBlank(tt.s); got != tt.want {
				t.Errorf("IsNotBlank(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{{
		name: "empty string",
		s:    "",
		want: false,
	}, {
		name: "non-empty string",
		s:    "test",
		want: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmpty(tt.s); got != tt.want {
				t.Errorf("IsNotEmpty(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{{
		name:  "empty string",
		input: "",
		want:  true,
	}, {
		name:  "non-empty string",
		input: "hello",
		want:  false,
	}, {
		name:  "whitespace string",
		input: "   ",
		want:  false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.input); got != tt.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{{
		name:  "empty string",
		input: "",
		want:  true,
	}, {
		name:  "all whitespace",
		input: "   \t\n\r",
		want:  true,
	}, {
		name:  "non-blank string",
		input: " hello ",
		want:  false,
	}, {
		name:  "non-whitespace characters",
		input: "a",
		want:  false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.input); got != tt.want {
				t.Errorf("IsBlank(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{{
		name: "lowercase string",
		s:    "hello",
		want: "Hello",
	}, {
		name: "mixed case string",
		s:    "hElLo",
		want: "Hello",
	}, {
		name: "empty string",
		s:    "",
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Capitalize(tt.s); got != tt.want {
				t.Errorf("Capitalize(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestUncapitalize(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{{
		name: "uppercase string",
		s:    "HELLO",
		want: "hELLO",
	}, {
		name: "empty string",
		s:    "",
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uncapitalize(tt.s); got != tt.want {
				t.Errorf("Uncapitalize(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		length  int
		padChar rune
		want    string
	}{{
		name:    "pad with spaces",
		s:       "123",
		length:  5,
		padChar: ' ',
		want:    "  123",
	}, {
		name:    "no padding needed",
		s:       "123",
		length:  3,
		padChar: '0',
		want:    "123",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadLeft(tt.s, tt.length, tt.padChar); got != tt.want {
				t.Errorf("PadLeft(%q, %d, %q) = %q, want %q", tt.s, tt.length, tt.padChar, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		length  int
		padChar rune
		want    string
	}{{
		name:    "pad with zeros",
		s:       "123",
		length:  5,
		padChar: '0',
		want:    "12300",
	}, {
		name:    "no padding needed",
		s:       "123",
		length:  3,
		padChar: '0',
		want:    "123",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadRight(tt.s, tt.length, tt.padChar); got != tt.want {
				t.Errorf("PadRight(%q, %d, %q) = %q, want %q", tt.s, tt.length, tt.padChar, got, tt.want)
			}
		})
	}
}

func TestRemovePrefix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   string
	}{{
		name:   "remove existing prefix",
		s:      "prefix_test",
		prefix: "prefix_",
		want:   "test",
	}, {
		name:   "no prefix to remove",
		s:      "test",
		prefix: "prefix_",
		want:   "test",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemovePrefix(tt.s, tt.prefix); got != tt.want {
				t.Errorf("RemovePrefix(%q, %q) = %q, want %q", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestRemoveSuffix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		suffix string
		want   string
	}{{
		name:   "remove existing suffix",
		s:      "test_suffix",
		suffix: "_suffix",
		want:   "test",
	}, {
		name:   "no suffix to remove",
		s:      "test",
		suffix: "_suffix",
		want:   "test",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveSuffix(tt.s, tt.suffix); got != tt.want {
				t.Errorf("RemoveSuffix(%q, %q) = %q, want %q", tt.s, tt.suffix, got, tt.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		old  string
		new  string
		want string
	}{{
		name: "replace all occurrences",
		s:    "hello hello",
		old:  "hello",
		new:  "world",
		want: "world world",
	}, {
		name: "no occurrences to replace",
		s:    "test",
		old:  "hello",
		new:  "world",
		want: "test",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Replace(tt.s, tt.old, tt.new); got != tt.want {
				t.Errorf("Replace(%q, %q, %q) = %q, want %q", tt.s, tt.old, tt.new, got, tt.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{{
		name:   "empty string",
		s:      "",
		substr: "a",
		want:   0,
	}, {
		name:   "single occurrence",
		s:      "abc",
		substr: "b",
		want:   1,
	}, {
		name:   "multiple occurrences",
		s:      "ababab",
		substr: "ab",
		want:   3,
	}, {
		name:   "unicode characters",
		s:      "你好，世界！世界你好",
		substr: "世界",
		want:   2,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Count(tt.s, tt.substr); got != tt.want {
				t.Errorf("Count(%q, %q) = %d, want %d", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestIsAllBlank(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want bool
	}{{
		name: "empty input",
		strs: []string{},
		want: true,
	}, {
		name: "all blank strings",
		strs: []string{"", " ", "\t", "\n"},
		want: true,
	}, {
		name: "contains non-blank string",
		strs: []string{"", "123", " "},
		want: false,
	}, {
		name: "all non-blank strings",
		strs: []string{"123", "abc"},
		want: false,
	}, {
		name: "single blank string",
		strs: []string{""},
		want: true,
	}, {
		name: "single whitespace string",
		strs: []string{" "},
		want: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllBlank(tt.strs...); got != tt.want {
				t.Errorf("IsAllBlank(%v) = %v, want %v", tt.strs, got, tt.want)
			}
		})
	}
}

func TestIsAllEmpty(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want bool
	}{{
		name: "empty input",
		strs: []string{},
		want: true,
	}, {
		name: "all empty strings",
		strs: []string{"", ""},
		want: true,
	}, {
		name: "contains non-empty string",
		strs: []string{"", "123"},
		want: false,
	}, {
		name: "all non-empty strings",
		strs: []string{"123", "abc"},
		want: false,
	}, {
		name: "contains whitespace",
		strs: []string{"", " "},
		want: false,
	}, {
		name: "single empty string",
		strs: []string{""},
		want: true,
	}, {
		name: "single whitespace string",
		strs: []string{" "},
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllEmpty(tt.strs...); got != tt.want {
				t.Errorf("IsAllEmpty(%v) = %v, want %v", tt.strs, got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		template string
		params   []string
		want     string
	}{{
		name:     "normal case with multiple placeholders",
		template: "Hello, {}! Today is {}.",
		params:   []string{"Alice", "Monday"},
		want:     "Hello, Alice! Today is Monday.",
	}, {
		name:     "more parameters than placeholders",
		template: "{} and {}",
		params:   []string{"A", "B", "C"},
		want:     "A and B",
	}, {
		name:     "fewer parameters than placeholders",
		template: "Name: {}, Age: {}",
		params:   []string{"Bob"},
		want:     "Name: Bob, Age: {}",
	}, {
		name:     "no parameters",
		template: "Hello, {}!",
		params:   []string{},
		want:     "Hello, {}!",
	}, {
		name:     "empty template",
		template: "",
		params:   []string{"test"},
		want:     "",
	}, {
		name:     "no placeholders",
		template: "Hello, World!",
		params:   []string{"extra"},
		want:     "Hello, World!",
	}, {
		name:     "special characters in parameters",
		template: "User {}: {}",
		params:   []string{"admin", "password123"},
		want:     "User admin: password123",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.template, tt.params...); got != tt.want {
				t.Errorf("Format(%q, %v) = %q, want %q", tt.template, tt.params, got, tt.want)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"empty", "", ""},
		{"lowercase", "hello", "HELLO"},
		{"mixed", "hElLo", "HELLO"},
		{"already_upper", "HELLO", "HELLO"},
		{"unicode", "héllô", "HÉLLÔ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUpper(tt.args); got != tt.want {
				t.Errorf("ToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"empty", "", ""},
		{"uppercase", "HELLO", "hello"},
		{"mixed", "hElLo", "hello"},
		{"already_lower", "hello", "hello"},
		{"unicode", "HÉLLÔ", "héllô"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLower(tt.args); got != tt.want {
				t.Errorf("ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"empty", "", ""},
		{"single_word", "hello", "hello"},
		{"snake_case", "hello_world", "helloWorld"},
		{"multiple_underscores", "hello__world", "helloWorld"},
		{"with_spaces", "hello world", "helloWorld"},
		{"mixed_case", "Hello_World", "HelloWorld"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.args); got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"empty", "", ""},
		{"single_word", "hello", "hello"},
		{"camel_case", "helloWorld", "hello_world"},
		{"pascal_case", "HelloWorld", "hello_world"},
		{"mixed_case", "helloWorldFooBar", "hello_world_foo_bar"},
		{"already_snake", "hello_world", "hello_world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.args); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"empty", "", false},
		{"digits", "12345", true},
		{"with_letters", "123a45", false},
		{"with_symbols", "123.45", false},
		{"unicode_numbers", "１２３４５", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumeric(tt.args); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"empty", "", false},
		{"letters", "abcXYZ", true},
		{"with_digits", "abc123", false},
		{"with_spaces", "a b c", false},
		{"unicode_letters", "αβγδε", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlpha(tt.args); got != tt.want {
				t.Errorf("IsAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAlphanumeric(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"empty", "", false},
		{"letters", "abcXYZ", true},
		{"digits", "12345", true},
		{"mixed", "abc123", true},
		{"with_symbols", "abc123!", false},
		{"with_spaces", "a b c", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlphanumeric(tt.args); got != tt.want {
				t.Errorf("IsAlphanumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{"empty", ""},
		{"simple_text", "hello world"},
		{"special_chars", "!@#$%^&*()"},
		{"unicode", "你好，世界"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Base64Encode(tt.args)
			decoded, err := Base64Decode(encoded)
			if err != nil {
				t.Errorf("Base64Decode() error = %v", err)
				return
			}
			if decoded != tt.args {
				t.Errorf("Base64Encode/Decode() = %v, want %v", decoded, tt.args)
			}
		})
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		name  string
		args  string
		left  int
		right int
		mask  rune
		want  string
	}{
		{"empty", "", 3, 4, '*', ""},
		{"phone_number", "13812345678", 3, 4, '*', "138****5678"},
		{"email", "test@example.com", 2, 3, '*', "te***********com"},
		{"short_string", "abc", 1, 1, '*', "a*c"},
		{"too_short", "ab", 1, 1, '*', "ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.args, tt.left, tt.right, tt.mask); got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomUUID(t *testing.T) {
	uuid1 := RandomUUID()
	uuid2 := RandomUUID()
	if len(uuid1) != 36 {
		t.Errorf("RandomUUID() length = %v, want 36", len(uuid1))
	}
	if uuid1 == uuid2 {
		t.Errorf("RandomUUID() generated duplicate UUIDs")
	}
}

// equalStringSlices checks if two string slices are equal
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
