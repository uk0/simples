package main

import (
	"context"
	"fmt"
	"regexp"
	"retroblog/utils"
	"strings"
)

// 可选：注入你自己的 Translator（示意）
type MyTranslator struct{}

func (MyTranslator) Translate(ctx context.Context, text, src, tgt string) (string, error) {
	// 这里接 OpenAI/智谱/DeepL 任意实现，返回英文句子
	return "", nil
}
func pinyin() {
	ctx := context.Background()

	opt := utils.Options{
		MaxLen:      64,
		KeepHash:    true,
		HashLen:     6,
		Lowercase:   true,
		DedupHyphen: true,
		PinyinMode:  utils.PinyinFull, // 或 utils.PinyinInitial
		Translator:  MyTranslator{},   // 若不想用外部翻译，置为 nil
		Glossary: map[string]string{ // 命中则直接替换为英文词
			"解决方案": "solution",
			"教程":   "tutorial",
		},
		GlossaryMaxLen: 8,
	}

	title := "Komari 飞书通知 JavaScript：DNS 可能存在问题（解决方案）"
	en, py := utils.TitleToEnglishAndPinyin(ctx, title, opt)
	// en:  (若 Translator 返回空) 按 Glossary + 拼音兜底生成英文词 slug
	// py:  zhong-wen-quan-pin-slug（示意）
	fmt.Println("EN:", en)
	fmt.Println("PY:", py)
}

func main() {
	pinyin()
	// 更全面的文件名匹配：支持中文、英文、数字、空格、常见符号
	re := regexp.MustCompile(`\[Attachment:\s*([^\[\]]+?)\s*\]`)
	//                                         ^^^^^^^^^ 匹配除了方括号外的所有字符（非贪婪）

	gen := utils.NewOGGenerator()
	meta, info, err := gen.GenerateFromURL("https://firsh.me/blog/0153.html")
	if err != nil { /* handle */
	}
	fmt.Println(meta) // 直接贴到 <head> 里
	fmt.Println(info) // 直接贴到 <head> 里
	// info 是结构化字段

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
