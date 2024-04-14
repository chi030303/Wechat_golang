package utils

import (
	"regexp"
	"strings"
)

// 检查用户输入是否包含数字
func containsNumber(input string) bool {
	// 定义包含数字的字符集合
	numbers := "0123456789"

	// 遍历输入字符串的每个字符
	for _, char := range input {
		// 将字符转换为字符串，并在数字集合中查找
		if strings.ContainsAny(string(char), numbers) {
			return true // 如果找到数字，返回true
		}
	}

	return false // 如果未找到数字，返回false
}

// 提取用户输入的数字
func extractNumbers(input string) []string {
	// 定义匹配数字的正则表达式
	re := regexp.MustCompile("[0-9]+")

	// 在输入字符串中查找所有匹配的数字
	matches := re.FindAllString(input, -1)

	return matches
}
