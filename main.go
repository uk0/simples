package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Article struct {
	Title string
	URL   string
}

type GitHubUser struct {
	AvatarURL string `json:"avatar_url"`
}

func main() {
	markdownDir := "./book"
	staticDir := "./static"
	articles := processMarkdownFiles(markdownDir, staticDir)
	avatarURL := getGitHubAvatar("uk0")
	generateIndexHTML(articles, avatarURL)

	// 指定静态文件目录
	fs := http.FileServer(http.Dir("./static"))

	// 使用 http.Handle() 函数将文件服务器注册到 /static 路径
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 使用 http.HandleFunc() 函数将 index.html 文件注册到根路径
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	// 启动服务器，监听在 8080 端口
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func processMarkdownFiles(markdownDir, staticDir string) []Article {
	var articles []Article

	err := filepath.Walk(markdownDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			title := extractTitle(string(content))
			md5sum := calculateMD5(content)
			url := fmt.Sprintf("%s/%s.html", staticDir, md5sum)

			articles = append(articles, Article{Title: title, URL: url})
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through markdown directory: %v\n", err)
	}

	return articles
}

func extractTitle(content string) string {
	re := regexp.MustCompile(`title:\s*(.+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return "Untitled"
}

func calculateMD5(content []byte) string {
	hash := md5.Sum(content)
	return hex.EncodeToString(hash[:])
}

func getGitHubAvatar(username string) string {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching GitHub user data: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return ""
	}

	var user GitHubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return ""
	}

	return user.AvatarURL
}

func generateIndexHTML(articles []Article, avatarURL string) {
	tmpl := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>打怪兽升级咯～</title>
	<!-- Cloudflare Web Analytics -->
	<script defer src='https://static.cloudflareinsights.com/beacon.min.js' data-cf-beacon='{"token": "38ef94cda57d43a1a9de524646e66805"}'></script>
    <!-- End Cloudflare Web Analytics -->
	<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=VT323&display=swap');
        	
        body {
            font-family: 'VT323', monospace;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #000;
            color: #00ff00;
        }
        header {
            text-align: center;
            margin-bottom: 20px;
        }
        .github-avatar {
            width: 80px;
            height: 80px;
            border-radius: 50%;
            border: 2px solid #00ff00;
            filter: grayscale(100%) brightness(150%) contrast(120%);
        }
        h1 {
            text-align: center;
            color: #00ff00;
            font-size: 2em;
            text-shadow: 0 0 5px #00ff00;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            margin-bottom: 10px;
            padding: 5px;
            border: 1px solid #00ff00;
        }
        a {
            color: #00ff00;
            text-decoration: none;
            display: block;
        }
        a:hover {
            text-decoration: underline;
        }
        .cursor {
            display: inline-block;
            width: 10px;
            height: 1em;
            background-color: #00ff00;
            animation: blink 0.7s steps(1) infinite;
            vertical-align: middle;
            margin-left: 2px;
        }
        @keyframes blink {
            0%, 100% { opacity: 0; }
            50% { opacity: 1; }
        }
        footer {
            text-align: center;
            margin-top: 20px;
            font-size: 0.8em;
            color: #00aa00;
        }
        @media (max-width: 600px) {
            body {
                padding: 10px;
            }
            .github-avatar {
                width: 60px;
                height: 60px;
            }
            h1 {
                font-size: 1.5em;
            }
            li {
                padding: 3px;
            }
        }
    </style>
	<script>
	window.onloadTurnstileCallback = function () {
    turnstile.render('body', {
        sitekey: '0x4AAAAAAAkDvnJ2XJn6mQqA',
        callback: function(token) {
			console.log(token);
			},
		});
	};
	</script>
</head>
<body>
    <header>
        <a href="https://github.com/uk0" target="_blank">
            <img src="{{.AvatarURL}}" alt="GitHub Avatar" class="github-avatar">
        </a>
    </header>
    <h1>君子论心不论迹。</h1>
    <ul id="article-list">
        {{range .Articles}}
        <li><a href="{{.URL}}"></a></li>
        {{end}}
    </ul>
    <footer>
        Power By Gemini TextGenerate
    </footer>

    <script>
        const articleList = document.getElementById('article-list');
        const items = articleList.getElementsByTagName('li');
        const titles = {{.Articles | json}}.map(article => article.Title);
        let currentItemIndex = 0;
        let currentCharIndex = 0;

        function typeNextChar() {
            if (currentItemIndex >= items.length) {
                return;
            }

            const item = items[currentItemIndex];
            const link = item.getElementsByTagName('a')[0];
            const title = titles[currentItemIndex];

            if (currentCharIndex < title.length) {
                link.textContent += title[currentCharIndex];
                currentCharIndex++;
                setTimeout(typeNextChar, Math.random() * 3 + 1);
            } else {
                link.innerHTML += '<span class="cursor"></span>';
                currentItemIndex++;
                currentCharIndex = 0;
                setTimeout(typeNextChar, 2);
            }
        }

        window.onload = typeNextChar;
    </script>
</body>
</html>
`

	t, err := template.New("index").Funcs(template.FuncMap{
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	}).Parse(tmpl)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	file, err := os.Create("index.html")
	if err != nil {
		fmt.Printf("Error creating index.html: %v\n", err)
		return
	}
	defer file.Close()

	data := struct {
		Articles  []Article
		AvatarURL string
	}{
		Articles:  articles,
		AvatarURL: avatarURL,
	}

	err = t.Execute(file, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	fmt.Println("index.html generated successfully.")
}
