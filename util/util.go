package util

func CutChineseString(s string, length int) string {
	var count, i int
	runes := []rune(s)
	for ; i < len(runes) && count < length; i++ {
		if runes[i] > 0x4e00 && runes[i] < 0x9fff { // 判断是否为中文字符
			count++
		}
	}
	return string(runes[:i])
}
