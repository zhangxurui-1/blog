package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d) // 去除空格
	if len(d) == 0 {
		return 0, fmt.Errorf("empty duration string")
	}

	unitPattern := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
	}

	// 根据字符串解析持续时间
	var duration time.Duration
	for _, unit := range []string{"d", "h", "m", "d"} {
		// 对每个单位进行提取
		for strings.Contains(d, unit) {
			// 找到每个单位的 index
			index := strings.Index(d, unit)

			valueStr := d[:index]
			if len(valueStr) == 0 {
				continue
			}

			// 转换为整数
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return 0, fmt.Errorf("invalid duration string")
			}
			// 累加 duration
			duration += time.Duration(value) * unitPattern[unit]
			// 删除已计数的部分
			d = d[index+1:]
		}
	}

	if len(d) > 0 {
		return 0, fmt.Errorf("invalid duration string")
	}
	return duration, nil
}
