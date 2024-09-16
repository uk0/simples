## simples
simples 是一个使用 Gemini API构建博客页面的一个工具，目标就是简化博客的搭建，只需要将文章放到指定的目录下，就可以自动构建博客页面。


#### 一个基于Markdown的简易博客


#### quick 

* 将文章放到book
* 运行 [gemini_markdown_to_html_test.py](gemini_markdown_to_html_test.py)
* 运行go build .
* ./simples 


```bash
go build .
nohup ./simples  >blog.log 2>&1 &
```

![img.png](img.png)


#### 简单干净

* https://firsh.me