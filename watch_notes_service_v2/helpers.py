from __future__ import annotations

import re

from markdownify import markdownify as md

# 匹配任意级标题
_HEADING_TAG_RE = re.compile(r"(<h[1-6][^>]*>)(.*?)(</h[1-6]>)",
                             flags=re.I | re.S)

def _collapse_ws_in_heading(match: re.Match) -> str:
    """把标题内部的换行/多空格收成一个空格。"""
    start, inner, end = match.groups()
    inner = re.sub(r"\s+", " ", inner.strip())
    return f"{start}{inner}{end}"


def _merge_consecutive_headings(html: str) -> str:
    """
    1) 将 </hN><hN> 直接替换成空格，实现同级标题合并。
    2) 移除只含 <br> 或空白的“空标题”。
    """
    # 合并同级、相邻的标题
    for level in range(1, 7):
        html = re.sub(
            rf"</h{level}>\s*<h{level}[^>]*>",
            " ",
            html,
            flags=re.I | re.S,
        )
    # 删除空标题
    html = re.sub(
        r"<h[1-6][^>]*>\s*(?:<br\s*/?>)?\s*</h[1-6]>",
        "",
        html,
        flags=re.I,
    )
    return html


def html_to_markdown(html: str) -> str:
    """
    Convert HTML to Markdown, ensuring that logically连续的标题
    不会被解析成多个 Markdown `#` 行。
    """
    # 0) 先合并/清理标题
    html = _merge_consecutive_headings(html)

    # 1) 压缩标题内部的多余空白
    html = _HEADING_TAG_RE.sub(_collapse_ws_in_heading, html)

    # 2) 保持 <ul> 之前的换行（原需求）
    html = html.replace("<ul", "<br><ul")

    # 3) markdownify
    md_text = md(html, heading_style="ATX", newline_style="SPACES")

    # 4) 把换行替换成带两个空格的软换行
    return md_text.replace("\n", "  \n")
