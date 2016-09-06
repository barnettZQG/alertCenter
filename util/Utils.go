package util

import "sort"

func GetLabelString(labels map[string]string) string {
	var data []string
	for k, _ := range labels {
		data = append(data, k)
	}
	sort.Strings(data)
	var result string
	for _, k := range data {
		result += k + labels[k]
	}
	return result
}

func FormatTime(oldTime string) (newTime string) {
	newTime = Substr(oldTime, 0, 19)
	return
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
