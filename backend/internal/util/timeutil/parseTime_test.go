package timeutil

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"纯小时", "2h", 2 * time.Hour, false},
		{"纯天", "7d", 7 * 24 * time.Hour, false},
		{"组合", "1h30m", 90 * time.Minute, false},
		{"全单位", "1d2h3m4s", 24*time.Hour + 2*time.Hour + 3*time.Minute + 4*time.Second, false},
		{"带空白", " 2h ", 2 * time.Hour, false},
		{"空串", "", 0, true},
		{"无单位数字", "100", 0, true},
		{"非法字符", "abc", 0, true},
		{"浮点数", "0.5h", 0, true},
		{"残缺尾缀", "1h30", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseDuration(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseDuration(%q) 应返回错误，得到 %v", tc.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDuration(%q) 意外报错: %v", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("ParseDuration(%q) = %v, 期望 %v", tc.input, got, tc.want)
			}
		})
	}
}
