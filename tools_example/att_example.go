package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`\[Attachment:\s*([\p{Han}A-Za-z0-9_.-]+)\]`)

	tests := []string{
		"[Attachment: kof97.zip]",
		"[Attachment: 漢字文件名.zip]",
		"[Attachment: 混合Mix123.zip]",
	}

	for _, t := range tests {
		if m := re.FindStringSubmatch(t); m != nil {
			fmt.Printf("匹配成功：%s -> %s\n", t, m[1])
		} else {
			fmt.Printf("匹配失败：%s\n", t)
		}
	}
}
