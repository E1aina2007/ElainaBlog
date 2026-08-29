package verifycode

import "testing"

func TestGenerateCode(t *testing.T) {
	code := GenerateCode(6)
	if len(code) != 6 {
		t.Fatalf("长度应为 6，得到 %d", len(code))
	}
	for _, ch := range code {
		if ch < '0' || ch > '9' {
			t.Errorf("验证码应仅含数字，出现 %q", ch)
		}
	}
}

func TestGenerateCodeDistribution(t *testing.T) {
	seen := map[rune]bool{}
	for i := 0; i < 100; i++ {
		for _, ch := range GenerateCode(6) {
			seen[ch] = true
		}
	}
	if len(seen) < 5 {
		t.Errorf("600 位随机数字应覆盖至少 5 种数字，实际 %d 种", len(seen))
	}
}
