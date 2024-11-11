import re
import os
import requests
from urllib.parse import urlparse, urlsplit


def clean_url(url):
    """清理 URL，移除 ! 后面的参数"""
    return url.split('!')[0]


def get_path_from_url(url):
    """从 URL 中提取路径部分"""
    # 首先清理 URL
    clean = clean_url(url)
    parsed_url = urlparse(clean)
    path = parsed_url.path.lstrip('/')  # 移除开头的斜杠
    return path


def download_image(url, save_dir):
    try:
        # 使用清理后的 URL 下载
        clean = clean_url(url)
        response = requests.get(clean, timeout=10)
        if response.status_code == 200:
            # 从清理后的 URL 获取相对路径
            relative_path = get_path_from_url(clean)
            # 完整的保存路径
            full_save_path = os.path.join(save_dir, relative_path)

            # 创建必要的子目录
            os.makedirs(os.path.dirname(full_save_path), exist_ok=True)

            # 保存文件
            with open(full_save_path, 'wb') as f:
                f.write(response.content)

            return relative_path
    except Exception as e:
        print(f"下载失败 {url}: {str(e)}")
    return None


def process_markdown_file(file_path, static_dir):
    # 确保静态文件夹存在
    os.makedirs(static_dir, exist_ok=True)

    # 读取 Markdown 文件
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # 使用正则表达式匹配完整的图片标记
    pattern = r'!\[(.*?)\]\((http[s]?://[^)]*?112firshme11224[^)]*?)\)'

    def replace_match(match):
        alt_text = match.group(1)
        url = match.group(2)
        relative_path = download_image(url, static_dir)
        if relative_path:
            return f'![{alt_text}](/static/{relative_path})'
        return match.group(0)  # 如果下载失败，保持原样

    # 替换所有匹配项
    new_content = re.sub(pattern, replace_match, content)

    # 保存修改后的文件
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)


def process_directory(directory_path, static_dir):
    # 处理目录下所有 .md 文件
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(root, file)
                print(f"处理文件: {file_path}")
                process_markdown_file(file_path, static_dir)


if __name__ == "__main__":
    # 设置目录路径
    markdown_dir = "./book"  # 当前目录
    static_dir = "./static"  # 静态文件存储目录

    # 处理所有文件
    process_directory(markdown_dir, static_dir)
    print("处理完成！")