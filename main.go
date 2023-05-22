package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/blackfriday/v2"
	"github.com/shirou/gopsutil/process"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var Totalaccept = sync.Map{}

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
type ReMsg struct {
	Id       int
	SendTime template.HTML
	Name     template.HTML
	Context  template.HTML
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

// 获取ip
func GetRequestIP(c *gin.Context) {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	var accept any
	var ok bool
	if accept, ok = Totalaccept.Load(reqIP); ok {
		c := accept.(int64)
		atomic.AddInt64(&c, 1)
		Totalaccept.Store(reqIP, c)
	} else {
		var c int64
		atomic.AddInt64(&c, 1)
		Totalaccept.Store(reqIP, c)
	}
	log.Printf("accept ip %s accept count %d \n", reqIP, accept)
	c.Next()
}

type Comment struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	ParentId  int       `json:"parentId"`
	RootId    int       `json:"rootId"`
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
	Website   string    `json:"website"`
	Page      string    `json:"page"`
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./comments.db")
	if err != nil {
		log.Fatal(err)
	}

	/**
	  CREATE TABLE IF NOT EXISTS comments(
	        id TEXT NOT NULL PRIMARY KEY,
	        rootId TEXT,
	        parentId TEXT,
	        username TEXT NOT NULL,
	        website TEXT NOT NULL,
	        content TEXT NOT NULL
	    )
	*/

	statement, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comments (
	    id INTEGER PRIMARY KEY,
		content TEXT,
		parentId INTEGER,
		rootId INTEGER,
		username TEXT,
		page TEXT,
		website TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)
`)
	statement.Exec()
}

func getComments(page string) ([]Comment, error) {
	//rows, err := db.Query("SELECT id, content, parentId, rootId, username, website, timestamp FROM comments WHERE rootId = ?", rootid)
	rows, err := db.Query("SELECT id, content, parentId, rootId, username, website, timestamp FROM comments  WHERE page = ?", page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Id, &comment.Content, &comment.ParentId, &comment.RootId, &comment.Username, &comment.Website, &comment.Timestamp)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func insertComment(comment Comment) error {
	statement, err := db.Prepare(`
		INSERT INTO comments (page,content, parentId, rootId, username, website) VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	_, err = statement.Exec(comment.Page, comment.Content, comment.ParentId, comment.RootId, comment.Username, comment.Website)
	if err != nil {
		return err
	}
	log.Println("statement.Insert")
	return nil
}

func StartServer() {
	r := gin.Default()
	r.Use(GetRequestIP)
	r.Use(static.Serve("/", static.LocalFile("static", true)))
	r.StaticFS("/static", http.Dir("./static"))
	r.GET("/comments", func(c *gin.Context) {
		page := c.Query("page")
		comments, err := getComments(page)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, comments)
	})

	r.POST("/comments", func(c *gin.Context) {
		var comment Comment
		err := c.BindJSON(&comment)
		if err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = insertComment(comment)
		if err != nil {
			log.Println(err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(201)
	})

	err := r.Run("0.0.0.0:80")
	if err != nil {
		return
	}
}

func gitpull() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	log.Println("pull ", exPath)
	cmd := exec.Command("git", "-C", exPath, "pull")
	cmd.Dir = ""
	err = cmd.Run()
	if err != nil {
		log.Println("error ")
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

var DB *sql.DB

func main() {
	InitDB()
	go func() {
		for {
			gitpull()
			tmp := GenPage()
			allPage = tmp
			GenIndex()
			time.Sleep(time.Second * 60)
			log.Println("update")
		}
	}()
	StartServer()
}
