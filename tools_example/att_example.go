package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// 更全面的文件名匹配：支持中文、英文、数字、空格、常见符号
	re := regexp.MustCompile(`\[Attachment:\s*([^\[\]]+?)\s*\]`)
	//                                         ^^^^^^^^^ 匹配除了方括号外的所有字符（非贪婪）

	tests := []string{
		"[Attachment: kof97.zip]",
		"[Attachment: 漢字文件名.zip]",
		"[Attachment: 混合Mix123.zip]",
		"[Attachment: 爱你就像爱生命 王小波,李银河.pdf]",
		"[Attachment: My File (2024) - Copy #1.docx]",
		"[Attachment:   test.txt   ]", // 测试多余空格
		"[Attachment: 项目报告_2024年Q1季度.xlsx]",
	}

	for _, t := range tests {
		if m := re.FindStringSubmatch(t); m != nil {
			// 去除首尾空格
			filename := strings.TrimSpace(m[1])
			fmt.Printf("匹配成功：%s -> %s\n", t, filename)
		} else {
			fmt.Printf("匹配失败：%s\n", t)
		}
	}
}
