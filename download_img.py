import re
import os
from time import sleep

import requests
from urllib.parse import urlparse, unquote

# 支持的图片扩展名
IMAGE_EXTENSIONS = ('.png', '.jpg', '.jpeg', '.gif', '.webp', '.bmp', '.ico')


def is_image_url(url):
    """检查是否为图片 URL"""
    clean_url = clean_url_and_params(url)
    return clean_url.lower().endswith(IMAGE_EXTENSIONS)


def clean_url_and_params(url):
    """清理 URL，移除 ! 后面的参数"""
    # URL 解码
    decoded_url = unquote(url)
    # 分割并获取主要 URL 部分
    base_url = decoded_url.split('!')[0]
    # 确保 URL 以支持的图片扩展名结尾
    for ext in IMAGE_EXTENSIONS:
        if base_url.lower().endswith(ext):
            # 获取扩展名的位置
            ext_position = base_url.lower().rindex(ext) + len(ext)
            return base_url[:ext_position]
    return base_url


def get_path_from_url(url):
    """从 URL 中提取路径部分"""
    parsed_url = urlparse(url)
    path = parsed_url.path.lstrip('/')  # 移除开头的斜杠
    return path


def download_image(url, save_dir):
    try:
        sleep(1.2)
        # 使用清理后的 URL 下载
        clean = clean_url_and_params(url)
        print(f"下载图片: {clean}")

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

            print(f"保存到: {full_save_path}")
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

    # 使用正则表达式匹配完整的图片标记，包括所有 http/https URL
    pattern = r'!\[(.*?)\]\((http[s]?://[^)]+)\)'

    def replace_match(match):
        alt_text = match.group(1)
        url = match.group(2)

        # 检查是否为图片 URL
        if not is_image_url(url):
            return match.group(0)

        relative_path = download_image(url, static_dir)
        if relative_path:
            return f'![{alt_text}](/static/{relative_path})'
        return match.group(0)  # 如果下载失败，保持原样

    # 替换所有匹配项
    new_content = re.sub(pattern, replace_match, content)

    # 检查内容是否发生变化
    if new_content != content:
        # 保存修改后的文件
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"已更新文件: {file_path}")


def process_directory(directory_path, static_dir):
    # 处理目录下所有 .md 文件
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(root, file)
                print(f"\n处理文件: {file_path}")
                process_markdown_file(file_path, static_dir)


if __name__ == "__main__":
    # 设置目录路径
    markdown_dir = "./book"  # 当前目录
    static_dir = "./static"  # 静态文件存储目录

    # 处理所有文件
    process_directory(markdown_dir, static_dir)
    print("\n处理完成！")