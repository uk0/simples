import html2text

def html_to_md_with_html2text(html: str) -> str:
    h = html2text.HTML2Text()
    h.body_width = 0               # 不自动折行
    h.ignore_images = True         # 如果不想要图片标签
    h.ignore_tables = False        # 如果想保留表格
    h.inline_links = False         # 用 reference-style 链接
    # 核心：html2text 会自动保留 ```…``` 代码块原样
    return h.handle(html)

if __name__ == "__main__":
    sample = """
    <p>下面是一段原始脚本：</p>
    ```bash
    export MVN_HOME=/opt/soft/apache-maven-3.6.0
    export PATH=$MVN_HOME/bin:$PATH
    ```
    <p>结束。</p>
    """
    print(html_to_md_with_html2text(sample))