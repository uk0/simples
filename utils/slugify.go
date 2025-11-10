package utils

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin" // go get github.com/mozillazg/go-pinyin
)

/* ================= 公共选项与接口 ================= */

type RomanizeMode int

const (
	PinyinFull    RomanizeMode = iota // 中文全拼：中文标题 -> zhong-wen-biao-ti
	PinyinInitial                     // 中文首字母：中文标题 -> zwbt
)

type Options struct {
	// 通用
	MaxLen      int  // 最终最长长度（不含扩展名），0 表示不限
	KeepHash    bool // 末尾追加短哈希（基于原始标题），避免重复
	HashLen     int  // 短哈希长度（默认 6）
	Lowercase   bool // 转小写（默认 true）
	DedupHyphen bool // 合并重复连字符（默认 true）
	AllowUnders bool // 是否允许下划线（默认 false）
	WordSep     byte // 词间分隔符（默认 '-'）
	PinyinMode  RomanizeMode
	PinyinSep   byte // 拼音片段间分隔符（默认 '-'）
	EnglishSep  byte // 英文词间分隔符（默认 '-'）

	// 英文版相关
	Translator     Translator        // 可选：外部翻译器（OpenAI/智谱等）
	Glossary       map[string]string // 可选：自定义中文→英文词汇表（优先最长匹配）
	GlossaryMaxLen int               // Glossary 最大词长（rune 数），默认 8
}

type Translator interface {
	Translate(ctx context.Context, text, source, target string) (string, error)
}

/* ================= 对外主函数 ================= */

// TitleToEnglishAndPinyin 返回两个 slug：英文单词版、拼音版（均仅 ASCII）
func TitleToEnglishAndPinyin(ctx context.Context, title string, opt Options) (englishSlug, pinyinSlug string) {
	normalized := strings.TrimSpace(title)
	if normalized == "" {
		normalized = "untitled"
	}

	// 缺省
	if opt.HashLen <= 0 {
		opt.HashLen = 6
	}
	if !opt.Lowercase {
		opt.Lowercase = true
	}
	if !opt.DedupHyphen {
		opt.DedupHyphen = true
	}
	if opt.WordSep == 0 {
		opt.WordSep = '-'
	}
	if opt.PinyinSep == 0 {
		opt.PinyinSep = '-'
	}
	if opt.EnglishSep == 0 {
		opt.EnglishSep = '-'
	}
	if opt.GlossaryMaxLen <= 0 {
		opt.GlossaryMaxLen = 8
	}

	// 1) 拼音版（可全拼/首字母）
	pinyinSlug = slugifyPinyin(normalized, opt)

	// 2) 英文版
	// 2.1 优先走外部翻译（更准确）
	var en string
	if opt.Translator != nil {
		if tr, err := opt.Translator.Translate(ctx, normalized, "zh", "en"); err == nil && strings.TrimSpace(tr) != "" {
			en = tr
		}
	}

	// 2.2 词典替换（最长匹配），把中文段替成英文词，再与原有英文单词合并
	if en == "" && len(opt.Glossary) > 0 && hasHan(normalized) {
		en = glossaryEnglish(normalized, opt)
	}

	// 2.3 兜底：把剩余中文转为“拼音英文”（可选首字母/全拼），保证可读且为 ASCII
	if en == "" {
		en = hanToRomanizedEnglish(normalized, opt)
	}

	englishSlug = slugifyAsciiWords(en, opt.EnglishSep, opt)

	// 3) 全局截断/哈希（对两个 slug 各自生效）
	englishSlug = postTrimAndHash(englishSlug, title, opt)
	pinyinSlug = postTrimAndHash(pinyinSlug, title, opt)

	return englishSlug, pinyinSlug
}

/* ================= 拼音版实现 ================= */

func slugifyPinyin(title string, opt Options) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	args.Heteronym = false

	var b strings.Builder
	b.Grow(len(title) * 2)

	pushSep := func() {
		if b.Len() == 0 {
			return
		}
		last := b.String()[b.Len()-1]
		if opt.DedupHyphen && (last == '-' || (opt.AllowUnders && last == '_')) {
			return
		}
		b.WriteByte(opt.WordSep)
	}

	for _, r := range title {
		switch {
		case unicode.Is(unicode.Han, r):
			ps := pinyin.SinglePinyin(r, args)
			if len(ps) > 0 && len(ps[0]) > 0 {
				if opt.PinyinMode == PinyinInitial {
					b.WriteByte(ps[0][0])
				} else {
					b.WriteString(ps[0])
				}
				b.WriteByte(opt.PinyinSep)
			} else {
				pushSep()
			}
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case r == ' ' || r == '-' || (opt.AllowUnders && r == '_'):
			b.WriteByte(opt.WordSep)
		default:
			// 丢弃其它符号
		}
	}

	s := normalizeSlug(b.String(), opt)
	return s
}

/* ================= 英文版实现 ================= */

// 2.2 Glossary 最长匹配替换：把中文串按最长词典项替换为英文词，保留原始英文
func glossaryEnglish(title string, opt Options) string {
	// 将标题按“中文片段/非中文片段”切分
	parts := splitHanAndNonHan(title)

	var out []string
	for _, p := range parts {
		if p.isHan {
			out = append(out, longestGlossaryReplace(p.text, opt.Glossary, opt.GlossaryMaxLen))
		} else {
			// 非中文段直接保留
			out = append(out, p.text)
		}
	}
	// 合并后输出为自然英文句子
	en := strings.Join(out, " ")
	en = spaceNormalize(en)
	return en
}

// 2.3 兜底：中文 -> 拼音英文（可首字母/全拼），非中文保留
func hanToRomanizedEnglish(title string, opt Options) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	args.Heteronym = false

	var out []string
	parts := splitHanAndNonHan(title)
	for _, p := range parts {
		if !p.isHan {
			out = append(out, p.text)
			continue
		}
		var b strings.Builder
		for _, r := range p.text {
			if unicode.Is(unicode.Han, r) {
				ps := pinyin.SinglePinyin(r, args)
				if len(ps) > 0 && len(ps[0]) > 0 {
					if opt.PinyinMode == PinyinInitial {
						b.WriteByte(ps[0][0])
					} else {
						b.WriteString(ps[0])
					}
					b.WriteByte(' ')
				}
			}
		}
		out = append(out, strings.TrimSpace(b.String()))
	}
	return spaceNormalize(strings.Join(out, " "))
}

/* ================= 工具函数 ================= */

func postTrimAndHash(s, title string, opt Options) string {
	if opt.Lowercase {
		s = strings.ToLower(s)
	}
	if opt.MaxLen > 0 && len(s) > opt.MaxLen {
		s = s[:opt.MaxLen]
		s = strings.Trim(s, "-_ ")
		if s == "" {
			s = "untitled"
		}
	}
	if opt.KeepHash {
		h := sha1.Sum([]byte(title))
		hex := fmt.Sprintf("%x", h[:])
		if opt.HashLen > len(hex) {
			opt.HashLen = len(hex)
		}
		if s != "" {
			s += "-"
		}
		s += hex[:opt.HashLen]
	}
	return s
}

func slugifyAsciiWords(s string, sep byte, opt Options) string {
	// 只保留字母数字，把其它分隔为 sep
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_' || r == '.' || r == '/' || r == '+' || r == '#' || r == ':' || r == ',':
			b.WriteByte(sep)
		default:
			// 其它统统丢弃（或转 sep）
			b.WriteByte(sep)
		}
	}
	return normalizeSlug(b.String(), opt)
}

func normalizeSlug(s string, opt Options) string {
	if opt.DedupHyphen {
		re := regexp.MustCompile(`[-_]+`)
		s = re.ReplaceAllString(s, "-")
	}
	s = strings.Trim(s, "-_ ")
	if opt.Lowercase {
		s = strings.ToLower(s)
	}
	if s == "" {
		s = "untitled"
	}
	return s
}

func hasHan(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

type piece struct {
	text  string
	isHan bool
}

// 把字符串切成 [中文片段]/[非中文片段]
func splitHanAndNonHan(s string) []piece {
	var res []piece
	var b strings.Builder
	curHan := false
	started := false

	flush := func() {
		if !started {
			return
		}
		res = append(res, piece{text: b.String(), isHan: curHan})
		b.Reset()
		started = false
	}

	for _, r := range s {
		isHan := unicode.Is(unicode.Han, r)
		if !started {
			started, curHan = true, isHan
			b.WriteRune(r)
			continue
		}
		if isHan != curHan {
			flush()
			started, curHan = true, isHan
			b.WriteRune(r)
			continue
		}
		b.WriteRune(r)
	}
	flush()
	return res
}

// 最长优先的 Glossary 替换：中文词 -> 英文词
func longestGlossaryReplace(han string, glossary map[string]string, maxLen int) string {
	rs := []rune(han)
	n := len(rs)
	var out []string
	for i := 0; i < n; {
		bestLen, bestWord := 0, ""
		upper := i + maxLen
		if upper > n {
			upper = n
		}
		for l := upper; l > i; l-- {
			k := string(rs[i:l])
			if v, ok := glossary[k]; ok {
				bestLen = l - i
				bestWord = v
				break
			}
		}
		if bestLen > 0 {
			out = append(out, bestWord)
			i += bestLen
		} else {
			// 没命中就跳过一个汉字（后续用拼音兜底时会处理）
			i++
		}
	}
	// 注意：此函数只做词典替换，不转拼音；最终会与非中文段一起进入 slugify 处理
	return strings.Join(out, " ")
}

func spaceNormalize(s string) string {
	s = strings.ReplaceAll(s, "　", " ")
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}
