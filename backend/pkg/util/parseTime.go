package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(t string) (time.Duration, error) {
	t = strings.TrimSpace(t)
	if len(t) == 0 {
		return 0, fmt.Errorf("empty duration string!\n")
	}

	unitPattern := map[string]time.Duration{
		"d": time.Hour * 24,
		"h": time.Hour,
		"m": time.Minute,
		"s": time.Second,
	}

	var totalDuration time.Duration
	for _, unit := range []string{"d", "h", "m", "s"} {
		if strings.Contains(t, unit) {
			unitIndex := strings.Index(t, unit)
			part := t[:unitIndex]
			if part == "" {
				part = "0"
			}
			val, err := strconv.Atoi(part)
			if err != nil {
				return 0, fmt.Errorf("invalid duration format: %v", err)
			}
			totalDuration += time.Duration(val) * unitPattern[unit]
			t = t[unitIndex+1:] // 截断已消费的部分
		}
	}

	if len(t) > 0 {
		return 0, fmt.Errorf("invalid duration format.")
	}

	return totalDuration, nil
}
