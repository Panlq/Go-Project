// split/split.go

package unit_test

import (
	"strings"
)

// Split ...
func Split(s, sep string) (result []string) {
	// 提前使用make函数将result初始化为一个容量足够大的切片，而不再像之前一样通过调用append函数来追加
	// 测试优化后的性能
	result = make([]string, 0, strings.Count(s, sep)+1)
	i := strings.Index(s, sep)
	for i > -1 {

		result = append(result, s[:i])
		s = s[i+len(sep):] // 考虑sep多字符的情况
		i = strings.Index(s, sep)
	}
	if len(s) != 0 {
		result = append(result, s)
	}
	return
}
