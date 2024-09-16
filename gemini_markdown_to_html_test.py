from datetime import datetime
import os
import hashlib
import random
from time import sleep

import google.generativeai as genai

# 设置你的 API 密钥
genai.configure(api_key='xxxxxxxxxxxxxxxxxxxxxxxxxxxx')

# 选择模型
model = genai.GenerativeModel('gemini-1.5-flash-001')


def call_gemini_api(prompt, chat, max_retries=3, delay=5):
    for attempt in range(max_retries):
        try:
            response = chat.send_message(prompt)
            return response
        except Exception as e:
            print(f"API request failed on attempt {attempt + 1}: {e}")
            if attempt < max_retries - 1:
                sleep_time = delay * (2 ** attempt)
                print(f"Retrying in {sleep_time} seconds...")
                sleep(sleep_time)
            else:
                raise Exception(f"API request failed after {max_retries} attempts") from e


def get_complete_response(chat, prompt):
    response = call_gemini_api(prompt, chat)
    full_response = response.text
    if response.usage_metadata.candidates_token_count >= 8192:
        continue_prompt = "请继续提供剩余的HTML内容,按照规则提供不要进行format，不要包含```"
        response = call_gemini_api(continue_prompt, chat)
        full_response += response.text
        # TODO fix
        full_response = full_response.replace("```html", "")
        full_response = full_response.replace("```", "")
    return full_response


def markdown_to_unique_html(content):
    current_date = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    chat = model.start_chat(history=[])

    prompt = f"""
    职责： 你需要使用Markdown内容创建一个HTML页面。
    规则：
        1.  响应式设计，最大宽度820像素,文本行高21px,代码行高24px，整体居中显示,代码不需要进行居中。
        2.  在meta标签中包含标题、描述和关键词等信息，你需要分析整体markdown来进行填充，并且将标题，关键字描述 居中显示。
        3.  需要发挥想象力创造这个页面的背景可以有一部分阴影(4px~8px)。
        4.  所有CSS应内联在<style>标签中。
        5.  在页面中添加一些JavaScript交互效果，被``包裹的内容高亮以及代码高亮使用DlHighlight.js实现，但是需要与主题色进行呼应。
        6.  应用随机种子: 666666661，你可以使用这个种子来生成随机数。
        7.  字体大小需要适当，h1(20px), h2(18px), h3(16px), h4(14px), h5(12px), h6(11px)。
        8.  在html最底部添加 <footer> 标签，包含 Power By Gemini TextGenerate {current_date}。
        9.  背景颜色分为两种一种是白天阅读模式，一种是夜晚阅读模式，你可以通过js获取当前时间来决定使用那种主题背景色。
        10. 图片标签[<img>]需要居中显示，图片大小600x375。
        11. 高亮代码的背景色需要和字体颜色不一样，更方便阅读，也需要将主题色考虑进去。 
        13. body的背景色:[动漫风格，Dos复古风格，Unix风格，像素风格]生成Html的时候需要用任意一个风格去生成。
    Markdown内容：
        {content}

    返回内容规则如下：
        必须只返回从<!DOCTYPE html>到</html>的原始HTML，包括内联的CSS和JavaScript。
        不要使用任何Markdown格式或代码块语法（```html...```）。
        返回的内容应该是有效的HTML，可以直接保存为.html文件并在浏览器中查看。
        不要返回其他无关的内容，只需要返回html，也不要进行总结和声明。
    """

    html_content = get_complete_response(chat, prompt)
    return html_content


def generate_md5(content):
    return hashlib.md5(content.encode()).hexdigest()


def process_markdown_files(input_folder, output_folder):
    if not os.path.exists(output_folder):
        os.makedirs(output_folder)

    for filename in os.listdir(input_folder):
        if filename.endswith('.md'):
            input_path = os.path.join(input_folder, filename)

            with open(input_path, 'r', encoding='utf-8') as file:
                markdown_content = file.read()

            md5_filename = generate_md5(markdown_content) + '.html'
            output_path = os.path.join(output_folder, md5_filename)

            if os.path.exists(output_path):
                print(f"Skipping {filename} as {md5_filename} already exists")
                continue

            try:
                html_content = markdown_to_unique_html(markdown_content)

                with open(output_path, 'w', encoding='utf-8') as file:
                    file.write(html_content)

                print(f"Converted {filename} to {md5_filename}")
                sleep(random.randint(5, 12))

            except Exception as e:
                print(f"Failed to convert {filename} to HTML: {e}")


if __name__ == "__main__":
    input_folder = '/Users/firshme/Desktop/work/simples/book'
    output_folder = '/Users/firshme/Desktop/work/simples/static'
    process_markdown_files(input_folder, output_folder)
