from datetime import datetime
import os
import hashlib
import markdown
from pathlib import Path
import re


def extract_yaml_front_matter(content):
    """提取 YAML front matter 并返回剩余内容和元数据"""
    yaml_pattern = re.compile(r'^---\s*\n(.*?)\n---\s*\n', re.DOTALL)
    match = yaml_pattern.match(content)

    if not match:
        return content, {
            'layout': 'post',
            'title': 'Untitled',
            'categories': '',
            'description': '',
            'keywords': ''
        }

    yaml_content = match.group(1)
    remaining_content = content[match.end():]

    # 解析 YAML 内容
    metadata = {}
    for line in yaml_content.split('\n'):
        line = line.strip()
        if line and ':' in line:
            key, value = line.split(':', 1)
            metadata[key.strip()] = value.strip()

    # 设置默认值
    default_metadata = {
        'layout': 'post',
        'title': 'Untitled',
        'categories': '',
        'description': '',
        'keywords': ''
    }

    # 合并默认值和实际值
    for key in default_metadata:
        if key not in metadata:
            metadata[key] = default_metadata[key]

    return remaining_content, metadata


def create_html_template(metadata, content, current_date):
    title = metadata.get('title', 'Untitled')
    description = metadata.get('description', '')
    keywords = metadata.get('keywords', '')
    categories = metadata.get('categories', '')

    return f"""
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{description}">
    <meta name="keywords" content="{keywords}">
    <title>{title}</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/default.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <style>
        :root {{
            --day-bg: #ffffff;
            --day-text: #333333;
            --night-bg: #1a1a1a;
            --night-text: #f0f0f0;
        }}

        body {{
            max-width: 820px;
            margin: 0 auto;
            padding: 20px;
            line-height: 21px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            transition: background-color 0.3s, color 0.3s;
        }}

        .metadata {{
            margin-bottom: 30px;
            padding: 15px;
            background: rgba(0, 0, 0, 0.05);
            border-radius: 5px;
        }}

        .metadata p {{
            margin: 5px 0;
            font-size: 0.9em;
        }}

        .day-mode {{
            background-color: var(--day-bg);
            color: var(--day-text);
        }}

        .night-mode {{
            background-color: var(--night-bg);
            color: var(--night-text);
        }}

        h1 {{ font-size: 19px; text-align: center; }}
        h2 {{ font-size: 18px; }}
        h3 {{ font-size: 16px; }}
        h4 {{ font-size: 14px; }}
        h5 {{ font-size: 12px; }}
        h6 {{ font-size: 11px; }}

        pre {{
            line-height: 24px;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            overflow-x: auto;
        }}

        code {{
            font-family: 'Fira Code', Consolas, Monaco, 'Courier New', monospace;
        }}

        img {{
            display: block;
            max-width: 600px;
            height: auto;
            margin: 20px auto;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }}

        .content {{
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            padding: 20px;
            border-radius: 8px;
            background-color: rgba(255, 255, 255, 0.05);
        }}

        footer {{
            text-align: center;
            margin-top: 40px;
            padding: 20px;
            border-top: 1px solid #ddd;
            font-size: 0.9em;
        }}

        @keyframes fadeIn {{
            from {{ opacity: 0; }}
            to {{ opacity: 1; }}
        }}

        .content {{
            animation: fadeIn 0.5s ease-in;
        }}
    </style>
</head>
<body class="day-mode">
    <div class="metadata">
        <p><strong>Title:</strong> {title}</p>
        <p><strong>Categories:</strong> {categories}</p>
        <p><strong>Description:</strong> {description}</p>
        <p><strong>Keywords:</strong> {keywords}</p>
    </div>
    <div class="content">
        {content}
    </div>
    <footer>
        Power By Python Markdown Generate {current_date}
    </footer>
    <script>
        document.addEventListener('DOMContentLoaded', (event) => {{
            document.querySelectorAll('pre code').forEach((block) => {{
                hljs.highlightBlock(block);
            }});
        }});

        function setThemeMode() {{
            const hour = new Date().getHours();
            const body = document.body;
            if (hour >= 18 || hour < 6) {{
                body.classList.remove('day-mode');
                body.classList.add('night-mode');
            }} else {{
                body.classList.remove('night-mode');
                body.classList.add('day-mode');
            }}
        }}

        setThemeMode();
        setInterval(setThemeMode, 60000);
    </script>
</body>
</html>
"""


def markdown_to_html(content):
    # 首先提取 YAML front matter
    markdown_content, metadata = extract_yaml_front_matter(content)

    extensions = [
        'markdown.extensions.extra',
        'markdown.extensions.codehilite',
        'markdown.extensions.toc',
        'markdown.extensions.tables',
        'markdown.extensions.fenced_code',
    ]

    # 转换 Markdown 为 HTML
    html_content = markdown.markdown(markdown_content, extensions=extensions)

    # 生成完整的 HTML
    current_date = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    return create_html_template(metadata, html_content, current_date)


def generate_md5(content):
    return hashlib.md5(content.encode()).hexdigest()


def process_markdown_files(input_folder, output_folder):
    input_path = Path(input_folder)
    output_path = Path(output_folder)

    # 创建输出目录
    output_path.mkdir(parents=True, exist_ok=True)

    # 处理所有 .md 文件
    for md_file in input_path.glob('**/*.md'):
        try:
            # 读取 Markdown 内容
            markdown_content = md_file.read_text(encoding='utf-8')

            # 生成输出文件名
            md5_filename = generate_md5(markdown_content) + '.html'
            output_file = output_path / md5_filename

            # 检查文件是否已存在
            if output_file.exists():
                print(f"overwrite {md_file.name} as {md5_filename}")

            # 转换为 HTML
            html_content = markdown_to_html(markdown_content)

            # 写入文件
            output_file.write_text(html_content, encoding='utf-8')
            print(f"Converted {md_file.name} to {md5_filename}")

        except Exception as e:
            print(f"Failed to convert {md_file.name}: {e}")


if __name__ == "__main__":
    input_folder = '/Users/firshme/Desktop/work/simples/book'
    output_folder = '/Users/firshme/Desktop/work/simples/static'
    process_markdown_files(input_folder, output_folder)