package progressutil

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewProgressBar(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		width     int
		fill      string
		empty     string
		output    io.Writer
		wantFill  string
		wantEmpty string
	}{{
		name:      "default values",
		total:     100,
		width:     20,
		fill:      "",
		empty:     "",
		output:    nil,
		wantFill:  "=",
		wantEmpty: " ",
	}, {
		name:      "custom values",
		total:     50,
		width:     10,
		fill:      "#",
		empty:     "-",
		output:    &bytes.Buffer{},
		wantFill:  "#",
		wantEmpty: "-",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewProgressBar(tt.total, tt.width, tt.fill, tt.empty, tt.output)
			if pb.fill != tt.wantFill {
				t.Errorf("fill = %v, want %v", pb.fill, tt.wantFill)
			}
			if pb.empty != tt.wantEmpty {
				t.Errorf("empty = %v, want %v", pb.empty, tt.wantEmpty)
			}
			if pb.total != tt.total {
				t.Errorf("total = %v, want %v", pb.total, tt.total)
			}
			if pb.width != tt.width {
				t.Errorf("width = %v, want %v", pb.width, tt.width)
			}
		})
	}
}

func TestSetProgress(t *testing.T) {
	pb := NewProgressBar(100, 20, "=", " ", nil)

	tests := []struct {
		name        string
		current     int
		wantErr     bool
		wantCurrent int
	}{{
		name:        "valid progress",
		current:     50,
		wantErr:     false,
		wantCurrent: 50,
	}, {
		name:        "progress exceeds total",
		current:     150,
		wantErr:     false,
		wantCurrent: 100,
	}, {
		name:        "negative progress",
		current:     -10,
		wantErr:     true,
		wantCurrent: 0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pb.SetProgress(tt.current); (err != nil) != tt.wantErr {
				t.Errorf("SetProgress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && pb.current != tt.wantCurrent {
				t.Errorf("current = %v, want %v", pb.current, tt.wantCurrent)
			}
		})
	}
}

func TestIncrement(t *testing.T) {
	pb := NewProgressBar(5, 20, "=", " ", nil)

	tests := []struct {
		name        string
		initial     int
		wantErr     bool
		wantCurrent int
	}{{
		name:        "increment valid",
		initial:     3,
		wantErr:     false,
		wantCurrent: 4,
	}, {
		name:        "increment complete",
		initial:     5,
		wantErr:     true,
		wantCurrent: 5,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb.SetProgress(tt.initial)
			if err := pb.Increment(); (err != nil) != tt.wantErr {
				t.Errorf("Increment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && pb.current != tt.wantCurrent {
				t.Errorf("current = %v, want %v", pb.current, tt.wantCurrent)
			}
		})
	}
}

func TestShow(t *testing.T) {
	buf := &bytes.Buffer{}
	pb := NewProgressBar(100, 10, "=", " ", buf)

	tests := []struct {
		name       string
		current    int
		wantErr    bool
		wantOutput string
	}{{
		name:       "valid progress",
		current:    50,
		wantErr:    false,
		wantOutput: "\r[=====     ] 50.00%",
	}, {
		name:       "progress exceeds total",
		current:    150,
		wantErr:    false,
		wantOutput: "[==========] 100.00%\n",
	}, {
		name:       "negative progress",
		current:    -10,
		wantErr:    true,
		wantOutput: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			err := pb.Show(tt.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("Show() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				output := strings.ReplaceAll(buf.String(), "\r", "")
				expected := strings.ReplaceAll(tt.wantOutput, "\r", "")
				if output != expected {
					t.Errorf("output = %q, want %q", output, expected)
				}
			}
		})
	}
}

func TestShow2(t *testing.T) {
	// 写个模拟下载的进度条
	pb := NewProgressBar(100, 10, "=", " ", nil)
	for i := 0; i <= 100; i++ {
		pb.Show(i)
		time.Sleep(100 * time.Millisecond)
	}
}
