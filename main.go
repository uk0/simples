package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"github.com/shirou/gopsutil/process"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

var allPage []Page

type SystemInfo struct {
	//Content string
	CPU       template.HTML
	Mem       template.HTML
	Thread    template.HTML
	Goroutine template.HTML
}
type MK struct {
	//Content string
	Content template.HTML
	Title   template.HTML
}

type Page struct {
	Name template.HTML `json:"name"`
	Url  template.HTML `json:"url"`
}

func MarkdownAutoNewline(str string) string {
	re, _ := regexp.Compile("\\ *\\n")
	str = re.ReplaceAllLiteralString(str, "  \n")
	//m.Content=strings.Replace(m.Content, "\n", "  \n", -1)
	reg := regexp.MustCompile("```([\\s\\S]*)```")
	//返回str中第一个匹配reg的字符串
	data := reg.Find([]byte(str))
	strs := strings.Replace(string(data), "  \n", "\n", -1)
	re, _ = regexp.Compile("```([\\s\\S]*)```")
	return re.ReplaceAllLiteralString(str, strs)
}

type Index struct {
	Urls       []Page     `json:"urls"`
	SystemInfo SystemInfo `json:"systemInfo"`
}

func StartServer() {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("static", true)))
	r.StaticFS("/static", http.Dir("./static"))
	err := r.Run()
	if err != nil {
		return
	}
}

func visit(paths *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Failed to access path", path, "error:", err)
			return err
		}
		if !info.IsDir() {
			*paths = append(*paths, path)
		}
		return nil
	}
}
func md5Str(origin string) string {
	m := md5.New()
	m.Write([]byte(origin))
	return hex.EncodeToString(m.Sum(nil))
}

func GenIndex() {
	index := &Index{
		Urls:       allPage,
		SystemInfo: getInfo(),
	}
	td, _ := template.ParseFiles("index_template.html")
	file, err := os.Create("static/index.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(index)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	err = td.Execute(file, index)
	if err != nil {
		return
	}
	err = file.Sync()
	if err != nil {
		return
	}
}

func GenPage() []Page {
	var paths []string
	var urls []Page
	root := "book"
	err := filepath.Walk(root, visit(&paths))
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", root, err)
	}
	var data []byte
	for _, path := range paths {
		data, err = ioutil.ReadFile(path)
		if err != nil {
			log.Println("error ")
		}
		fName := strings.Split(path, "/")
		output := blackfriday.Run(data, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak))
		content := template.HTML(MarkdownAutoNewline(string(output)))
		title := template.HTML(fName[len(fName)-1])
		td, _ := template.ParseFiles("template.html")
		mk := MK{Content: content, Title: title}
		saveFileName := fmt.Sprintf("static/%s.html", md5Str(fName[len(fName)-1]))
		urls = append(urls, Page{Url: template.HTML(saveFileName), Name: title})
		file, err := os.Create(saveFileName)
		if err != nil {
			log.Fatal(err)
		}
		err = td.Execute(file, mk)
		if err != nil {
			return nil
		}
		err = file.Close()
		if err != nil {
			return nil
		}
	}
	return urls
}

func getInfo() SystemInfo {
	p, _ := process.NewProcess(int32(os.Getpid()))
	cpuPercent, _ := p.Percent(time.Second)
	cp := cpuPercent / float64(runtime.NumCPU())
	// 获取进程占用内存的比例
	mp, _ := p.MemoryPercent()
	// 创建的线程数
	threadCount := pprof.Lookup("threadcreate").Count()
	// Goroutine数
	gNum := runtime.NumGoroutine()

	return SystemInfo{
		CPU:       template.HTML(fmt.Sprintf("%v", cp)),
		Mem:       template.HTML(fmt.Sprintf("%v", mp)),
		Thread:    template.HTML(fmt.Sprintf("%v", threadCount)),
		Goroutine: template.HTML(fmt.Sprintf("%v", gNum)),
	}
}
func main() {
	go func() {
		for {
			tmp := GenPage()
			allPage = tmp
			GenIndex()
			time.Sleep(time.Second * 60)
			log.Println("update")
		}
	}()
	StartServer()
}
