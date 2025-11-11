import sys
import os
import sqlite3
import gzip
import struct
import re
import shutil
import json
import html
from datetime import datetime, timedelta, timezone
from typing import Optional, List, Dict, Any, Tuple
from pathlib import Path
from enum import IntEnum

import helpers

DEBUG = False


class ProtoVarInt:
    """处理Protocol Buffer变长整数"""

    @staticmethod
    def decode(data: bytes, pos: int = 0) -> Tuple[int, int]:
        """解码varint，返回(值, 新位置)"""
        result = 0
        shift = 0
        while pos < len(data):
            byte = data[pos]
            pos += 1
            result |= (byte & 0x7F) << shift
            if not (byte & 0x80):
                return result, pos
            shift += 7
        raise ValueError("Invalid varint")


class AppleNoteConstants(IntEnum):
    """Apple Notes 常量定义"""
    # Style types
    STYLE_TYPE_TITLE = 0
    STYLE_TYPE_HEADING = 1
    STYLE_TYPE_SUBHEADING = 2
    STYLE_TYPE_MONOSPACED = 4
    STYLE_TYPE_NUMBERED_LIST = 100
    STYLE_TYPE_DOTTED_LIST = 101
    STYLE_TYPE_DASHED_LIST = 102
    STYLE_TYPE_CHECKBOX = 103
    STYLE_TYPE_BLOCK_QUOTE = 104

    # Font weights
    FONT_TYPE_DEFAULT = 0
    FONT_TYPE_BOLD = 1
    FONT_TYPE_ITALIC = 2
    FONT_TYPE_BOLD_ITALIC = 3

    # Alignments
    STYLE_ALIGNMENT_LEFT = 0
    STYLE_ALIGNMENT_CENTER = 1
    STYLE_ALIGNMENT_RIGHT = 2
    STYLE_ALIGNMENT_JUSTIFY = 3

    # Emphasis styles
    EMPHASIS_PURPLE = 1
    EMPHASIS_PINK = 2
    EMPHASIS_ORANGE = 3
    EMPHASIS_MINT = 4
    EMPHASIS_BLUE = 5


class AppleNotesProtoParser:
    """Apple Notes Protocol Buffer 解析器"""

    # Wire types
    VARINT = 0
    FIXED64 = 1
    LENGTH_DELIMITED = 2
    START_GROUP = 3
    END_GROUP = 4
    FIXED32 = 5

    def __init__(self, data: bytes):
        self.data = data
        self.pos = 0

    def read_tag(self) -> Optional[Tuple[int, int]]:
        """读取字段标签，返回(field_number, wire_type)"""
        if self.pos >= len(self.data):
            return None

        tag, self.pos = ProtoVarInt.decode(self.data, self.pos)
        field_number = tag >> 3
        wire_type = tag & 0x07
        return field_number, wire_type

    def read_varint(self) -> int:
        """读取varint值"""
        value, self.pos = ProtoVarInt.decode(self.data, self.pos)
        return value

    def read_length_delimited(self) -> bytes:
        """读取长度分隔的数据"""
        length = self.read_varint()
        value = self.data[self.pos:self.pos + length]
        self.pos += length
        return value

    def read_fixed32(self) -> int:
        """读取32位固定长度整数"""
        value = struct.unpack('<I', self.data[self.pos:self.pos + 4])[0]
        self.pos += 4
        return value

    def read_fixed64(self) -> int:
        """读取64位固定长度整数"""
        value = struct.unpack('<Q', self.data[self.pos:self.pos + 8])[0]
        self.pos += 8
        return value

    def read_float(self) -> float:
        """读取32位浮点数"""
        value = struct.unpack('<f', self.data[self.pos:self.pos + 4])[0]
        self.pos += 4
        return value

    def skip_field(self, wire_type: int):
        """跳过不需要的字段"""
        if wire_type == self.VARINT:
            self.read_varint()
        elif wire_type == self.FIXED64:
            self.pos += 8
        elif wire_type == self.LENGTH_DELIMITED:
            length = self.read_varint()
            self.pos += length
        elif wire_type == self.FIXED32:
            self.pos += 4


class ParagraphStyle:
    """段落样式"""

    def __init__(self):
        self.style_type = None
        self.alignment = None
        self.indent_amount = 0
        self.checklist = None
        self.block_quote = None


class Font:
    """字体信息"""

    def __init__(self):
        self.font_name = None
        self.point_size = None


class Color:
    """颜色信息"""

    def __init__(self):
        self.red = 0
        self.green = 0
        self.blue = 0
        self.alpha = 1

    def to_hex(self) -> str:
        """转换为十六进制颜色"""
        r = int(self.red * 255)
        g = int(self.green * 255)
        b = int(self.blue * 255)
        return f"#{r:02X}{g:02X}{b:02X}"


class Checklist:
    """清单信息"""

    def __init__(self):
        self.uuid = None
        self.done = 0


class AttributeRun:
    """属性运行 - 代表一段具有相同格式的文本"""

    def __init__(self):
        self.length = 0
        self.paragraph_style = None
        self.font = None
        self.font_weight = None
        self.underlined = 0
        self.strikethrough = 0
        self.superscript = 0
        self.link = None
        self.color = None
        self.emphasis_style = None
        self.attachment_info = None
        self.position = 0


class Attachment:
    """表示一个附件"""

    def __init__(self):
        self.identifier = ""
        self.type_uti = ""
        self.filename = ""
        self.original_filename = ""
        self.source_filename = ""
        self.position = -1
        self.length = 0
        self.media_path = ""
        self.size = 0
        self.data = None


class AppleNote:
    """表示一个Apple Note"""

    def __init__(self):
        self.text = ""
        self.title = ""
        self.snippet = ""
        self.creation_date = None
        self.modified_iso = None # iso 时间
        self.created_iso = None  # iso 时间
        self.modification_date = None
        self.attachments = []
        self.attachment_positions = {}
        self.attribute_runs = []  # 存储所有的属性运行
        self.tables = []
        self.checklists = []
        self.links = []
        self.styles = []
        self.raw_data = None
        self.raw_protobuf = None
        self.attachment_extractor = None # 父执行器
        self.formatted_text = ""
        self.html_content = ""  # HTML格式的内容
        self.folder_name = ""

    def parse_from_blob(self, blob_data: bytes) -> bool:
        """从ZDATA blob解析笔记"""
        try:
            self.raw_protobuf = blob_data

            # 尝试gzip解压
            try:
                decompressed = gzip.decompress(blob_data)
            except:
                decompressed = blob_data

            self.raw_data = decompressed

            # 解析protobuf结构
            result = self._parse_protobuf(decompressed)

            if result:
                # 生成HTML内容
                self._generate_html_content()
                # 生成格式化文本
                self._generate_formatted_text()

            return result

        except Exception as e:
            print(f"Error parsing blob: {e}")
            return self._fallback_parse(blob_data)

    def _parse_protobuf(self, data: bytes) -> bool:
        """解析Protocol Buffer格式的笔记数据"""
        parser = AppleNotesProtoParser(data)

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Document wrapper (field 2)
            if field_number == 2 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                doc_data = parser.read_length_delimited()
                self._parse_document(doc_data)
            else:
                parser.skip_field(wire_type)

        return bool(self.text)

    def _parse_document(self, data: bytes):
        """解析Document结构"""
        parser = AppleNotesProtoParser(data)

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            if field_number == 2 and wire_type == AppleNotesProtoParser.VARINT:
                version = parser.read_varint()
            elif field_number == 3 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                note_data = parser.read_length_delimited()
                self._parse_note_content(note_data)
            else:
                parser.skip_field(wire_type)

    def _parse_note_content(self, data: bytes):
        """解析Note内容，包括文本和属性运行"""
        parser = AppleNotesProtoParser(data)
        current_position = 0

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Note text (field 2)
            if field_number == 2 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                text_data = parser.read_length_delimited()
                try:
                    self.text = text_data.decode('utf-8', errors='replace')
                except:
                    self.text = str(text_data)

            # Attribute runs (field 5)
            elif field_number == 5 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                attr_data = parser.read_length_delimited()
                attr_run = self._parse_attribute_run(attr_data, current_position)
                if attr_run:
                    self.attribute_runs.append(attr_run)
                    current_position += attr_run.length

            else:
                parser.skip_field(wire_type)

    def _parse_attribute_run(self, data: bytes, position: int) -> AttributeRun:
        """解析属性运行（样式信息）"""
        parser = AppleNotesProtoParser(data)
        attr_run = AttributeRun()
        attr_run.position = position

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Length (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.VARINT:
                attr_run.length = parser.read_varint()

            # Paragraph style (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                style_data = parser.read_length_delimited()
                attr_run.paragraph_style = self._parse_paragraph_style(style_data)

            # Font (field 3)
            elif field_number == 3 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                font_data = parser.read_length_delimited()
                attr_run.font = self._parse_font(font_data)

            # Font weight (field 5)
            elif field_number == 5 and wire_type == AppleNotesProtoParser.VARINT:
                attr_run.font_weight = parser.read_varint()

            # Underlined (field 6)
            elif field_number == 6 and wire_type == AppleNotesProtoParser.VARINT:
                attr_run.underlined = parser.read_varint()

            # Strikethrough (field 7)
            elif field_number == 7 and wire_type == AppleNotesProtoParser.VARINT:
                attr_run.strikethrough = parser.read_varint()

            # Superscript (field 8)
            elif field_number == 8 and wire_type == AppleNotesProtoParser.VARINT:
                attr_run.superscript = parser.read_varint()

            # Link (field 9)
            elif field_number == 9 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                link_data = parser.read_length_delimited()
                try:
                    attr_run.link = link_data.decode('utf-8', errors='ignore')
                    self.links.append({'url': attr_run.link, 'position': position})
                except:
                    pass

            # Color (field 10)
            elif field_number == 10 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                color_data = parser.read_length_delimited()
                attr_run.color = self._parse_color(color_data)

            # Emphasis style (field 11)
            elif field_number == 11 and wire_type == AppleNotesProtoParser.VARINT:
                attr_run.emphasis_style = parser.read_varint()

            # Attachment (field 12)
            elif field_number == 12 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                attach_data = parser.read_length_delimited()
                attachment = self._parse_attachment(attach_data)
                if attachment:
                    attachment.position = position
                    attachment.length = attr_run.length
                    self.attachments.append(attachment)
                    self.attachment_positions[position] = attachment
                    attr_run.attachment_info = attachment

            else:
                parser.skip_field(wire_type)

        return attr_run

    def _parse_paragraph_style(self, data: bytes) -> ParagraphStyle:
        """解析段落样式"""
        parser = AppleNotesProtoParser(data)
        style = ParagraphStyle()

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Style type (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.VARINT:
                style.style_type = parser.read_varint()

            # Alignment (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.VARINT:
                style.alignment = parser.read_varint()

            # Indent amount (field 4)
            elif field_number == 4 and wire_type == AppleNotesProtoParser.VARINT:
                style.indent_amount = parser.read_varint()

            # Checklist (field 5)
            elif field_number == 5 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                checklist_data = parser.read_length_delimited()
                style.checklist = self._parse_checklist(checklist_data)
                if style.checklist:
                    self.checklists.append(style.checklist)

            # Block quote (field 6)
            elif field_number == 6 and wire_type == AppleNotesProtoParser.VARINT:
                style.block_quote = parser.read_varint()

            else:
                parser.skip_field(wire_type)

        return style

    def _parse_checklist(self, data: bytes) -> Checklist:
        """解析清单项"""
        parser = AppleNotesProtoParser(data)
        checklist = Checklist()

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # UUID (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                uuid_data = parser.read_length_delimited()
                checklist.uuid = uuid_data.hex()

            # Done status (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.VARINT:
                checklist.done = parser.read_varint()

            else:
                parser.skip_field(wire_type)

        return checklist

    def _parse_font(self, data: bytes) -> Font:
        """解析字体信息"""
        parser = AppleNotesProtoParser(data)
        font = Font()

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Font name (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                font_name = parser.read_length_delimited()
                font.font_name = font_name.decode('utf-8', errors='ignore')

            # Point size (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.FIXED32:
                font.point_size = parser.read_float()

            else:
                parser.skip_field(wire_type)

        return font

    def _parse_color(self, data: bytes) -> Color:
        """解析颜色信息"""
        parser = AppleNotesProtoParser(data)
        color = Color()

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Red (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.FIXED32:
                color.red = parser.read_float()

            # Green (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.FIXED32:
                color.green = parser.read_float()

            # Blue (field 3)
            elif field_number == 3 and wire_type == AppleNotesProtoParser.FIXED32:
                color.blue = parser.read_float()

            # Alpha (field 4)
            elif field_number == 4 and wire_type == AppleNotesProtoParser.FIXED32:
                color.alpha = parser.read_float()

            else:
                parser.skip_field(wire_type)

        return color

    def _parse_attachment(self, data: bytes) -> Attachment:
        """解析附件"""
        parser = AppleNotesProtoParser(data)
        attachment = Attachment()

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Attachment identifier (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                id_data = parser.read_length_delimited()
                try:
                    attachment.identifier = id_data.decode('utf-8', errors='ignore')
                except:
                    attachment.identifier = id_data.hex()

            # Type UTI (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                type_data = parser.read_length_delimited()
                try:
                    attachment.type_uti = type_data.decode('utf-8', errors='ignore')
                except:
                    attachment.type_uti = str(type_data)

            else:
                parser.skip_field(wire_type)

        return attachment if attachment.identifier else None

    def _generate_html_content(self):
        """生成HTML格式的内容"""
        if not self.text or not self.attribute_runs:
            self.html_content = html.escape(self.text) if self.text else ""
            return

        html_parts = []
        html_parts.append('<!DOCTYPE html>\n<html>\n<head>\n<meta charset="UTF-8">\n')
        html_parts.append('<style>\n')
        html_parts.append('''
            body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
            h1 { font-size: 28px; font-weight: bold; }
            h2 { font-size: 24px; font-weight: bold; }
            h3 { font-size: 20px; font-weight: bold; }
            pre { font-family: "SF Mono", Monaco, "Courier New", monospace; background: #f5f5f5; padding: 10px; }
            ul.checklist { list-style: none; }
            ul.checklist li.checked::before { content: "☑ "; }
            ul.checklist li.unchecked::before { content: "☐ "; }
            ul.dotted { list-style-type: disc; }
            ul.dashed { list-style-type: circle; }
            blockquote.block-quote { border-left: 3px solid #ccc; padding-left: 10px; margin-left: 0; }
        ''')
        html_parts.append('</style>\n</head>\n<body>\n')

        # 处理每个属性运行
        for i, attr_run in enumerate(self.attribute_runs):
            start = attr_run.position
            end = start + attr_run.length
            text_segment = self.text[start:end] if start < len(self.text) else ""

            if not text_segment:
                continue

            # 应用样式
            styled_html = self._apply_styles_to_html(text_segment, attr_run)
            html_parts.append(styled_html)

        html_parts.append('\n</body>\n</html>')
        self.html_content = ''.join(html_parts)

    def _apply_styles_to_html(self, text: str, attr_run: AttributeRun) -> str:
        """将样式应用到HTML"""
        result = []

        # 处理段落样式
        if attr_run.paragraph_style:
            style = attr_run.paragraph_style

            # 标题和列表
            if style.style_type == AppleNoteConstants.STYLE_TYPE_TITLE:
                result.append('<h1>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_HEADING:
                result.append('<h2>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_SUBHEADING:
                result.append('<h3>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_MONOSPACED:
                result.append('<pre>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_NUMBERED_LIST:
                result.append('<ol><li>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_DOTTED_LIST:
                result.append('<ul class="dotted"><li>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_DASHED_LIST:
                result.append('<ul class="dashed"><li>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_CHECKBOX:
                checked = "checked" if style.checklist and style.checklist.done else "unchecked"
                result.append(f'<ul class="checklist"><li class="{checked}">')

            # 缩进
            if style.indent_amount:
                for _ in range(style.indent_amount):
                    result.append('<blockquote>')

            # 对齐
            if style.alignment == AppleNoteConstants.STYLE_ALIGNMENT_CENTER:
                result.append('<div style="text-align: center">')
            elif style.alignment == AppleNoteConstants.STYLE_ALIGNMENT_RIGHT:
                result.append('<div style="text-align: right">')
            elif style.alignment == AppleNoteConstants.STYLE_ALIGNMENT_JUSTIFY:
                result.append('<div style="text-align: justify">')

        # 处理文本样式
        style_attrs = []

        # 字体
        if attr_run.font:
            if attr_run.font.font_name:
                style_attrs.append(f"font-family: '{attr_run.font.font_name}'")
            if attr_run.font.point_size:
                style_attrs.append(f"font-size: {attr_run.font.point_size}px")

        # 颜色
        if attr_run.color:
            style_attrs.append(f"color: {attr_run.color.to_hex()}")

        # 强调样式
        if attr_run.emphasis_style:
            if attr_run.emphasis_style == AppleNoteConstants.EMPHASIS_PURPLE:
                style_attrs.append("color: #FF00FF; background-color: #BA55D333")
            elif attr_run.emphasis_style == AppleNoteConstants.EMPHASIS_PINK:
                style_attrs.append("color: #FF4081; background-color: #D5000044")
            elif attr_run.emphasis_style == AppleNoteConstants.EMPHASIS_ORANGE:
                style_attrs.append("color: #FBC02D; background-color: #FF6F0022")
            elif attr_run.emphasis_style == AppleNoteConstants.EMPHASIS_MINT:
                style_attrs.append("color: #8DE5DB; background-color: #289C8ECC")
            elif attr_run.emphasis_style == AppleNoteConstants.EMPHASIS_BLUE:
                style_attrs.append("color: #BBDEFB; background-color: #2196F3")

        if style_attrs:
            result.append(f'<span style="{"; ".join(style_attrs)}">')

        # 字体样式
        if attr_run.font_weight == AppleNoteConstants.FONT_TYPE_BOLD:
            result.append('<b>')
        elif attr_run.font_weight == AppleNoteConstants.FONT_TYPE_ITALIC:
            result.append('<i>')
        elif attr_run.font_weight == AppleNoteConstants.FONT_TYPE_BOLD_ITALIC:
            result.append('<b><i>')

        if attr_run.underlined:
            result.append('<u>')

        if attr_run.strikethrough:
            result.append('<s>')

        if attr_run.superscript == 1:
            result.append('<sup>')
        elif attr_run.superscript == -1:
            result.append('<sub>')

        # 链接
        if attr_run.link:
            result.append(f'<a href="{html.escape(attr_run.link)}">')

        # 附件
        if attr_run.attachment_info:
            att = attr_run.attachment_info
            details = self.attachment_extractor.get_attachment_details(att.identifier)
            print("Processing attachment in HTML:", att.identifier, details)
            result.append(f'[Attachment: {html.escape(details.get("original_filename", details.get("media_filename", "unknown")))}]')
        else:
            # 添加文本内容
            escaped_text = html.escape(text)
            escaped_text = escaped_text.replace('\n', '<br>\n')
            escaped_text = escaped_text.replace('\u2028', '<br>')
            escaped_text = escaped_text.replace('\u0000', '\u2400')
            result.append(escaped_text)

        # 关闭标签
        if attr_run.link:
            result.append('</a>')

        if attr_run.superscript == 1:
            result.append('</sup>')
        elif attr_run.superscript == -1:
            result.append('</sub>')

        if attr_run.strikethrough:
            result.append('</s>')

        if attr_run.underlined:
            result.append('</u>')

        if attr_run.font_weight == AppleNoteConstants.FONT_TYPE_BOLD:
            result.append('</b>')
        elif attr_run.font_weight == AppleNoteConstants.FONT_TYPE_ITALIC:
            result.append('</i>')
        elif attr_run.font_weight == AppleNoteConstants.FONT_TYPE_BOLD_ITALIC:
            result.append('</i></b>')

        if style_attrs:
            result.append('</span>')

        # 关闭段落标签
        if attr_run.paragraph_style:
            style = attr_run.paragraph_style

            if style.alignment in [AppleNoteConstants.STYLE_ALIGNMENT_CENTER,
                                   AppleNoteConstants.STYLE_ALIGNMENT_RIGHT,
                                   AppleNoteConstants.STYLE_ALIGNMENT_JUSTIFY]:
                result.append('</div>')

            if style.indent_amount:
                for _ in range(style.indent_amount):
                    result.append('</blockquote>')

            if style.style_type == AppleNoteConstants.STYLE_TYPE_TITLE:
                result.append('</h1>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_HEADING:
                result.append('</h2>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_SUBHEADING:
                result.append('</h3>')
            elif style.style_type == AppleNoteConstants.STYLE_TYPE_MONOSPACED:
                result.append('</pre>')
            elif style.style_type in [AppleNoteConstants.STYLE_TYPE_NUMBERED_LIST,
                                      AppleNoteConstants.STYLE_TYPE_DOTTED_LIST,
                                      AppleNoteConstants.STYLE_TYPE_DASHED_LIST,
                                      AppleNoteConstants.STYLE_TYPE_CHECKBOX]:
                result.append(
                    '</li></ul>' if style.style_type != AppleNoteConstants.STYLE_TYPE_NUMBERED_LIST else '</li></ol>')

        return ''.join(result)

    def _generate_formatted_text(self):
        """生成包含附件标记的格式化文本"""
        if not self.text:
            self.formatted_text = ""
            return

        sorted_attachments = sorted(self.attachment_positions.items(), key=lambda x: x[0], reverse=True)

        formatted = self.text

        for position, attachment in sorted_attachments:
            if attachment.filename:
                marker = f"[file:attachments/{attachment.filename}]"
            else:
                marker = f"[file:attachment_{attachment.identifier[:8]}]"

            if position < len(formatted):
                if position < len(formatted) and formatted[position] in ['￼', '\ufffc']:
                    formatted = formatted[:position] + marker + formatted[position + 1:]
                else:
                    formatted = formatted[:position] + marker + formatted[position:]

        self.formatted_text = formatted

    def _fallback_parse(self, data: bytes) -> bool:
        """备用解析方法"""
        try:
            try:
                decompressed = gzip.decompress(data)
            except:
                decompressed = data

            text = decompressed.decode('utf-8', errors='ignore')

            if len(text) > 12:
                text = text[12:]

            result = ""
            for ch in text:
                if ch.isprintable() or ch.isspace() or ch == '￼':
                    result += ch
                else:
                    if len(result) > 10:
                        break

            self.text = result
            self.formatted_text = result
            return bool(self.text)

        except Exception as e:
            print(f"Fallback parse error: {e}")
            return False


# [保留原有的 AttachmentExtractor 类，不变]
class AttachmentExtractor:
    """附件提取器"""

    def __init__(self, db_path: str):
        self.db_path = db_path
        self.base_path = Path(db_path).parent
        self.conn = sqlite3.connect(db_path)

        # Apple Notes 附件可能的存储位置
        self.attachment_paths = [
            self.base_path / "Accounts",
            '''
            cd /Users/firshme/Library/Group Containers/group.com.apple.notes/Accounts/XXXXX-XXXX-4F9F-A9AB-XXXX/Media/
            ❯ du -sh *
              0B	FallbackImages
            1.7G	Media
              0B	Paper
            276M	Previews
             24K	Thumbnails
            '''
        ]

        self._check_database_schema()

    def _check_database_schema(self):
        """检查数据库表结构"""
        cursor = self.conn.cursor()
        cursor.execute("PRAGMA table_info(ZICCLOUDSYNCINGOBJECT)")
        columns = cursor.fetchall()
        self.available_columns = [col[1] for col in columns]

    def get_attachment_details(self, attachment_id: str) -> Dict[str, Any]:
        """从数据库获取附件详细信息"""
        cursor = self.conn.cursor()

        try:
            basic_query = """
                          SELECT ZIDENTIFIER,
                                 ZTYPEUTI,
                                 ZFILENAME,
                                 ZMEDIA,
                                 ZCREATIONDATE,
                                 ZMODIFICATIONDATE,
                                 ZTITLE,
                                 ZTITLE1,
                                 ZNAME
                          FROM ZICCLOUDSYNCINGOBJECT
                          WHERE ZIDENTIFIER = ?
                             OR ZIDENTIFIER LIKE ?
                          """

            results = cursor.execute(basic_query, (attachment_id, f"%{attachment_id}%")).fetchall()

            if results:
                row = results[0]
                details = {
                    'identifier': row[0],
                    'type_uti': row[1],
                    'filename': row[2],
                    'media_id': row[3],
                    'creation_date': row[4],
                    'modification_date': row[5],
                    'title': row[6],
                    'title1': row[7],
                    'name': row[8],
                    'original_filename': row[2] or row[6] or row[7] or row[8]
                }

                if details['media_id']:
                    try:
                        media_query = """
                                      SELECT ZFILENAME, ZIDENTIFIER, ZGENERATION1
                                      FROM ZICCLOUDSYNCINGOBJECT
                                      WHERE Z_PK = ?
                                      """
                        media_results = cursor.execute(media_query, (details['media_id'],)).fetchall()
                        if media_results:
                            details['media_filename'] = media_results[0][0]
                            details['media_identifier'] = media_results[0][1]
                            details['media_layer1_dir'] = media_results[0][2]
                            if media_results[0][0]:
                                details['original_filename'] = media_results[0][0]
                    except Exception:
                        pass

                return details

        except Exception:
            pass

        return {}

    def find_attachment_file(self, attachment_id: str,media_layer1_dir:str = None, filename: str = None,) -> Optional[Path]:
        """在文件系统中查找附件文件"""
        patterns = []

        if filename:
            patterns.append(filename)
            name_without_ext = os.path.splitext(filename)[0]
            patterns.append(f"{name_without_ext}*")

        if attachment_id:
            clean_id = re.sub(r'[^a-zA-Z0-9]', '', attachment_id)
            if clean_id:
                patterns.extend([
                    f"*{clean_id}*",
                    f"*{clean_id[:16]}*" if len(clean_id) > 16 else clean_id,
                    f"*{clean_id[:8]}*" if len(clean_id) > 8 else clean_id
                ])

        for base_path in self.attachment_paths:
            if not base_path.exists():
                continue

            for pattern in patterns:
                try:
                    for file_path in base_path.rglob(pattern):
                        if file_path.is_file() and file_path.stat().st_size > 0:
                            return file_path
                except Exception:
                    continue

        try:
            notes_dir = self.base_path
            if attachment_id:
                short_id = attachment_id[:8] if len(attachment_id) > 8 else attachment_id
                for file_path in notes_dir.rglob(f"*{short_id}*"):
                    if file_path.is_file() and file_path.stat().st_size > 0:
                        return file_path
        except Exception:
            pass

        return None

    def extract_attachment(self, attachment: Attachment, output_dir: Path) -> Tuple[bool,bool]:
        """提取并保存附件"""
        try:
            details = self.get_attachment_details(attachment.identifier)

            if not details and not attachment.identifier:
                return False,True

            attachment.original_filename = details.get('original_filename', '')

            print("      - Extracting attachment:", details)

            filename = None
            if (details.get('media_filename') == details['original_filename']) and details['media_id'] is not None:
                filename = details['media_filename']
            else:
                return False,True
            print("&" * 60)

            if not filename:
                ext = self._uti_to_extension(attachment.type_uti)
                base_name = attachment.identifier[:8] if attachment.identifier else "attachment"
                filename = f"{base_name}{ext}"

            if '.' not in filename:
                ext = self._uti_to_extension(attachment.type_uti)
                if ext:
                    filename = f"{filename}{ext}"

            filename = self._sanitize_filename(filename)
            attachment.filename = filename
            output_path = output_dir / filename

            source_file = self.find_attachment_file(
                attachment.identifier,
                details.get('media_layer1_dir'),
                details.get('filename') or details.get('media_filename')
            )

            if source_file:
                shutil.copy2(source_file, output_path)
                print(f"      ✓ Extracted: {filename}")
                attachment.source_filename = source_file.name
                return True,False

            if details.get('media_identifier') and details.get('media_layer1_dir') and details.get('media_filename'):
                #  通过sql 获取 上层文件夹，再组合路径查找
                source_file = self.find_attachment_file(details.get('media_identifier'),details.get('media_layer1_dir'),details.get('media_filename'))
                if source_file:
                    shutil.copy2(source_file, output_path)
                    print(f"      ✓ Extracted: {filename}")
                    attachment.source_filename = source_file.name
                    return True,False

            print(f"      ✗ Could not find file: {filename}")

            info_file = output_dir / f"{filename}.info.txt"
            with open(info_file, 'w') as f:
                f.write(f"Attachment Information\n")
                f.write(f"=====================\n")
                f.write(f"Identifier: {attachment.identifier}\n")
                f.write(f"Type: {attachment.type_uti}\n")
                f.write(f"Position in text: {attachment.position}\n")
                if details:
                    for key, value in details.items():
                        if value:
                            f.write(f"{key}: {value}\n")

            return False,True

        except Exception as e:
            print(f"      ✗ Error: {e}")
            return False,True

    def _uti_to_extension(self, uti: str) -> str:
        """将UTI转换为文件扩展名"""
        if not uti:
            return ".dat"

        uti_lower = uti.lower()


        uti_map = {
            'public.jpeg': '.jpg',
            'com.sun.java-archive': '.jar',
            'public.python-script': '.py',
            'org.webmproject.webp': '.webp',
            'public.zip-archive': '.zip',
            'public.png': '.png',
            'public.heic': '.heic',
            'public.heif': '.heif',
            'public.pdf': '.pdf',
            'public.plain-text': '.txt',
            'public.rtf': '.rtf',
            'public.html': '.html',
            'public.xml': '.xml',
            'public.json': '.json',
            'public.mpeg-4': '.mp4',
            'public.mpeg-4-audio': '.m4a',
            'public.mp3': '.mp3',
            'public.movie': '.mov',
            'com.apple.quicktime-movie': '.mov',
            'public.avi': '.avi',
            'public.data': '.dat',
            'com.microsoft.word.doc': '.doc',
            'com.microsoft.excel.xls': '.xls',
            'com.microsoft.powerpoint.ppt': '.ppt',
            'org.openxmlformats.wordprocessingml.document': '.docx',
            'org.openxmlformats.spreadsheetml.sheet': '.xlsx',
            'org.openxmlformats.presentationml.presentation': '.pptx',
            'com.apple.webarchive': '.webarchive',
            'public.vcard': '.vcf',
            'public.image': '.img'
        }

        for key, ext in uti_map.items():
            if key in uti_lower:
                return ext

        if '.' in uti:
            parts = uti.split('.')
            last_part = parts[-1].lower()
            if len(last_part) <= 5 and last_part.isalnum():
                return f".{last_part}"

        return ".dat"

    def _sanitize_filename(self, filename: str) -> str:
        """清理文件名"""
        filename = re.sub(r'[<>:"/\\|?*\x00-\x1f]', '_', filename)
        if len(filename) > 200:
            name, ext = os.path.splitext(filename)
            filename = name[:200 - len(ext)] + ext
        return filename or "unnamed"

    def close(self):
        """关闭数据库连接"""
        if self.conn:
            self.conn.close()


class AppleNotesExporter:
    """Apple Notes 导出器"""

    def __init__(self, db_path: str, output_dir: str):

        if Path(db_path).exists():
            db_path = str(Path(db_path).resolve())
        else:
            raise FileNotFoundError(f"Database not found: {db_path}")
        self._current_notes_ids = []
        self.db_path = db_path
        self.output_dir = output_dir
        self.attachment_extractor = AttachmentExtractor(db_path)
        self.stats = {
            'total': 0,
            'success': 0,
            'failed': 0,
            'with_attachments': 0,
            'attachments_extracted': 0,
            'attachments_failed': 0,
            'with_tables': 0,
            'with_checklists': 0
        }

    def export_all(self, folder_name: str = None):
        """导出所有笔记或指定文件夹的笔记（兼容新版本）"""
        if not os.path.exists(self.db_path):
            raise FileNotFoundError(f"Database not found: {self.db_path}")

        os.makedirs(self.output_dir, exist_ok=True)

        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        folder_id = None
        if folder_name:
            print(f"Searching for folder: '{folder_name}'")

            # 新版本：检查 ZTITLE2 字段
            folder_search_query = """
                SELECT DISTINCT Z_PK, ZTITLE, ZTITLE1, ZTITLE2, ZNAME, Z_ENT
                FROM ZICCLOUDSYNCINGOBJECT
                WHERE Z_ENT IN (14, 15)  -- 兼容新旧版本
                  AND (
                    LOWER(ZTITLE2) = LOWER(?) OR  -- 新版本主要在这里
                    LOWER(ZTITLE) = LOWER(?) OR
                    LOWER(ZTITLE1) = LOWER(?) OR
                    LOWER(ZNAME) = LOWER(?)
                  )
            """

            cursor.execute(folder_search_query, (folder_name, folder_name, folder_name, folder_name))
            folder_results = cursor.fetchall()

            if folder_results:
                folder_id = folder_results[0][0]
                # 优先使用 ZTITLE2
                actual_folder_name = folder_results[0][3] or folder_results[0][1] or folder_results[0][2] or \
                                     folder_results[0][4] or folder_name
                z_ent = folder_results[0][5]
                print(f"✓ Found folder '{actual_folder_name}' (ID: {folder_id}, Z_ENT: {z_ent})")

                # 验证文件夹中的笔记数量
                cursor.execute("""
                    SELECT COUNT(*)
                    FROM ZICCLOUDSYNCINGOBJECT n
                    INNER JOIN ZICNOTEDATA nd ON nd.ZNOTE = n.Z_PK
                    WHERE n.ZFOLDER = ?
                      AND nd.ZDATA IS NOT NULL
                """, (folder_id,))
                note_count = cursor.fetchone()[0]
                print(f"  Folder contains {note_count} notes")

            else:
                # 如果没找到，列出所有文件夹供选择
                print(f"⚠️ Folder '{folder_name}' not found.")
                print("\nAvailable folders:")

                cursor.execute("""
                    SELECT f.Z_PK,
                           f.ZTITLE,
                           f.ZTITLE1,
                           f.ZTITLE2,
                           COUNT(DISTINCT n.Z_PK) as note_count
                    FROM ZICCLOUDSYNCINGOBJECT f
                    INNER JOIN ZICCLOUDSYNCINGOBJECT n ON n.ZFOLDER = f.Z_PK
                    INNER JOIN ZICNOTEDATA nd ON nd.ZNOTE = n.Z_PK
                    WHERE f.Z_ENT IN (14, 15)
                      AND nd.ZDATA IS NOT NULL
                    GROUP BY f.Z_PK
                    ORDER BY note_count DESC
                """)

                all_folders = cursor.fetchall()

                if all_folders:
                    for fid, title, title1, title2, count in all_folders:
                        # 优先使用 ZTITLE2
                        display_name = title2 or title or title1 or f"Folder_{fid}"

                        if folder_name.lower() in display_name.lower():
                            print(f"  ✓ ID {fid}: '{display_name}' ({count} notes) <- POSSIBLE MATCH")
                            folder_id = fid
                            folder_name = display_name
                        else:
                            print(f"    ID {fid}: '{display_name}' ({count} notes)")

                    # 特殊处理：对于 blog 文件夹
                    if not folder_id and folder_name.lower() == 'blog':
                        # 直接使用 ID 122
                        print(f"\n⭐ Using known folder ID 122 for blog")
                        folder_id = 122

                if not folder_id:
                    print("\nNo matching folder found. Exporting all notes instead.")

        # 构建查询（使用 ZTITLE2）
        if folder_id:
            print(f"\nUse folder_id Exporting notes from folder ID {folder_id} ('{folder_name}')")

            query = """
                WITH base AS (
                  SELECT
                    nd.Z_PK                                            AS note_id,
                    nd.ZDATA                                           AS data,
                    COALESCE(n.ZTITLE, n.ZTITLE1, n.ZTITLE2)           AS title,
                    n.ZSNIPPET                                         AS snippet,
                
                    /* 可能分布在不同列：先合并出原始创建/修改时间 */
                    COALESCE(n.ZCREATIONDATE, n.ZCREATIONDATE3, n.ZCREATIONDATE1) AS raw_created,
                    COALESCE(n.ZMODIFICATIONDATE, n.ZMODIFICATIONDATE1) AS raw_modified,
                
                    n.ZFOLDER                                          AS folder_id,
                    COALESCE(f.ZTITLE2, f.ZTITLE, f.ZTITLE1, f.ZNAME, 'Unknown') AS folder_name
                  FROM ZICNOTEDATA nd
                  JOIN ZICCLOUDSYNCINGOBJECT n ON nd.ZNOTE = n.Z_PK
                  LEFT JOIN ZICCLOUDSYNCINGOBJECT f ON n.ZFOLDER = f.Z_PK
                  WHERE nd.ZDATA IS NOT NULL
                )
                SELECT
                  note_id, data, title, snippet,
                  raw_created, raw_modified,              -- 保留原始数值以便调试
                  folder_id, folder_name,
                
                  /* 将 CFAbsoluteTime（2001-01-01 起）转 ISO。单位可能是秒/毫秒/微秒，分段判断 */
                  CASE
                    WHEN raw_created IS NULL THEN NULL
                    WHEN raw_created > 10000000000000    THEN datetime((raw_created/1000000.0) + 978307200, 'unixepoch')
                    WHEN raw_created > 10000000000       THEN datetime((raw_created/1000.0)   + 978307200, 'unixepoch')
                    ELSE                                       datetime(raw_created           + 978307200, 'unixepoch')
                  END AS created_iso,
                
                  CASE
                    WHEN raw_modified IS NULL THEN NULL
                    WHEN raw_modified > 10000000000000   THEN datetime((raw_modified/1000000.0) + 978307200, 'unixepoch')
                    WHEN raw_modified > 10000000000      THEN datetime((raw_modified/1000.0)   + 978307200, 'unixepoch')
                    ELSE                                       datetime(raw_modified          + 978307200, 'unixepoch')
                  END AS modified_iso
                
                FROM base where folder_name = ? and folder_id = ?
                ORDER BY folder_id, raw_modified DESC NULLS LAST;
            """

            results = cursor.execute(query, (folder_name, folder_id)).fetchall()

        else:
            print("\nExporting all notes... not found folder_id ")

            query = """
                 WITH base AS (
                  SELECT
                    nd.Z_PK                                            AS note_id,
                    nd.ZDATA                                           AS data,
                    COALESCE(n.ZTITLE, n.ZTITLE1, n.ZTITLE2)           AS title,
                    n.ZSNIPPET                                         AS snippet,
                
                    /* 可能分布在不同列：先合并出原始创建/修改时间 */
                    COALESCE(n.ZCREATIONDATE, n.ZCREATIONDATE3, n.ZCREATIONDATE1) AS raw_created,
                    COALESCE(n.ZMODIFICATIONDATE, n.ZMODIFICATIONDATE1) AS raw_modified,
                
                    n.ZFOLDER                                          AS folder_id,
                    COALESCE(f.ZTITLE2, f.ZTITLE, f.ZTITLE1, f.ZNAME, 'Unknown') AS folder_name
                  FROM ZICNOTEDATA nd
                  JOIN ZICCLOUDSYNCINGOBJECT n ON nd.ZNOTE = n.Z_PK
                  LEFT JOIN ZICCLOUDSYNCINGOBJECT f ON n.ZFOLDER = f.Z_PK
                  WHERE nd.ZDATA IS NOT NULL
                )
                SELECT
                  note_id, data, title, snippet,
                  raw_created, raw_modified,              -- 保留原始数值以便调试
                  folder_id, folder_name,
                
                  /* 将 CFAbsoluteTime（2001-01-01 起）转 ISO。单位可能是秒/毫秒/微秒，分段判断 */
                  CASE
                    WHEN raw_created IS NULL THEN NULL
                    WHEN raw_created > 10000000000000    THEN datetime((raw_created/1000000.0) + 978307200, 'unixepoch')
                    WHEN raw_created > 10000000000       THEN datetime((raw_created/1000.0)   + 978307200, 'unixepoch')
                    ELSE                                       datetime(raw_created           + 978307200, 'unixepoch')
                  END AS created_iso,
                
                  CASE
                    WHEN raw_modified IS NULL THEN NULL
                    WHEN raw_modified > 10000000000000   THEN datetime((raw_modified/1000000.0) + 978307200, 'unixepoch')
                    WHEN raw_modified > 10000000000      THEN datetime((raw_modified/1000.0)   + 978307200, 'unixepoch')
                    ELSE                                       datetime(raw_modified          + 978307200, 'unixepoch')
                  END AS modified_iso
                
                FROM base where folder_name = '?'
                ORDER BY folder_id, raw_modified DESC NULLS LAST;
            """

            results = cursor.execute(query,(folder_name,)).fetchall()

        self.stats['total'] = len(results)

        print(f"Found {self.stats['total']} notes to export")
        print("=" * 60)

        for row in results:
            self._export_note(row)
            print(f"{row} row notes found")


        conn.close()
        self.attachment_extractor.close()

        self._print_stats()

     # 2001-01-01 00:00:00  到 1970-01-01 的秒数

    def cf_to_utc(self,cf_val):
        CF_OFFSET = 978307200
        """将 Apple CFAbsoluteTime（可能是秒/毫秒/微秒）转换为 UTC datetime；None -> None"""
        if cf_val is None:
            return None
        try:
            x = float(cf_val)
        except Exception:
            return None
        # 单位判断
        if x > 1e14:  # 微秒
            x = x / 1_000_000.0
        elif x > 1e10:  # 毫秒
            x = x / 1_000.0
        # 剩下的按秒
        ts = x + CF_OFFSET
        return datetime.fromtimestamp(ts, tz=timezone.utc)

    def _export_note(self, row):
        """导出单个笔记"""
        note_id, data, title, snippet, created, modified, folder_id, folder_name,created_iso,modified_iso = row
        self._current_notes_ids.append(f"{folder_name}_{note_id:04d}")
        print(f"\nProcessing Note #{note_id}")
        if title:
            print(f"  Title: {title}")
        if folder_name:
            print(f"  Folder: {folder_name}")

        note = AppleNote()
        note.title = title or ""
        note.folder_name = folder_name or ""
        note.snippet = snippet or ""
        note.creation_date = created
        note.modification_date = modified
        note.created_iso = created_iso
        note.modified_iso = modified_iso
        note.modification_date = modified
        note.attachment_extractor = self.attachment_extractor

        if note.parse_from_blob(data):
            #TODO 先倒出文件
            self._save_note(note_id, note, folder_name, data)
            self.stats['success'] += 1

            if note.attachments:
                self.stats['with_attachments'] += 1
            if note.tables:
                self.stats['with_tables'] += 1
            if note.checklists:
                self.stats['with_checklists'] += 1

            print(f"  ✓ Successfully exported")

            preview = note.text[:150].replace('\n', ' ')
            if len(note.text) > 150:
                preview += "..."
            print(f"  Preview: {preview}")
        else:
            self.stats['failed'] += 1
            print(f"  ✗ Failed to parse")

    def _save_note(self, note_id: int, note: AppleNote, folder_name: Optional[str], raw_blob: bytes):
        """保存笔记到文件"""
        # 创建文件夹结构
        if folder_name:
            folder_path = os.path.join(self.output_dir, self._sanitize_filename(folder_name))
            os.makedirs(folder_path, exist_ok=True)
        else:
            folder_path = self.output_dir

        if not note.title:
            note.title = note.text.split('\n')[0][:50]

        note_dir_name = f"{note_id:04d}"

        note_dir_name = self._sanitize_filename(note_dir_name)
        note_dir_path = os.path.join(folder_path, note_dir_name)
        os.makedirs(note_dir_path, exist_ok=True)

        # 保存原始二进制数据
        if DEBUG:
            raw_file_path = os.path.join(note_dir_path, "raw_data.bin")
            with open(raw_file_path, 'wb') as f:
                f.write(raw_blob)

        # 保存解压后的原始数据
        if note.raw_data and DEBUG:
            decompressed_file_path = os.path.join(note_dir_path, "decompressed_data.bin")
            with open(decompressed_file_path, 'wb') as f:
                f.write(note.raw_data)

        # 保存HTML版本
        if note.html_content:
            html_file_path = os.path.join(note_dir_path, "note.html.meta")
            with open(html_file_path, 'w', encoding='utf-8') as f:
                f.write(note.html_content)
            md = helpers.html_to_markdown(note.html_content)
            md_file_path = os.path.join(note_dir_path, "note.md")
            with open(md_file_path, 'w', encoding='utf-8') as fmd:
                fmd.write(md)




        # 保存样式信息为JSON
        if note.attribute_runs and DEBUG:
            styles_file_path = os.path.join(note_dir_path, "styles.json")
            styles_data = {
                'attribute_runs': [],
                'links': note.links,
                'checklists': []
            }

            for attr_run in note.attribute_runs:
                run_data = {
                    'position': attr_run.position,
                    'length': attr_run.length
                }

                if attr_run.paragraph_style:
                    run_data['paragraph_style'] = {
                        'style_type': attr_run.paragraph_style.style_type,
                        'alignment': attr_run.paragraph_style.alignment,
                        'indent_amount': attr_run.paragraph_style.indent_amount
                    }

                if attr_run.font:
                    run_data['font'] = {
                        'name': attr_run.font.font_name,
                        'size': attr_run.font.point_size
                    }

                if attr_run.font_weight:
                    run_data['font_weight'] = attr_run.font_weight

                if attr_run.color:
                    run_data['color'] = attr_run.color.to_hex()

                if attr_run.link:
                    run_data['link'] = attr_run.link

                styles_data['attribute_runs'].append(run_data)

            for checklist in note.checklists:
                styles_data['checklists'].append({
                    'uuid': checklist.uuid,
                    'done': checklist.done
                })

            with open(styles_file_path, 'w', encoding='utf-8') as f:
                json.dump(styles_data, f, indent=2, ensure_ascii=False)

        # 如果有附件，先提取附件
        if note.attachments:
            attachments_dir = os.path.join(note_dir_path, "attachments")
            os.makedirs(attachments_dir, exist_ok=True)

            print(f"    Extracting {len(note.attachments)} attachment(s)...")

            for att in note.attachments:
                print(f"export attachment Info {att.type_uti}")
                success,has_skip = self.attachment_extractor.extract_attachment(
                    att, Path(attachments_dir)
                )
                if success:
                    self.stats['attachments_extracted'] += 1
                elif has_skip :
                    print(f"\033[33m Export Attachments Skipped [del for attachments set. ]: {note.title}\033[0m")
                    # 在attachments内删除这个元素，说明不是标准的文件，可能是hash tag
                    note.attachments.remove(att)
                else:
                    self.stats['attachments_failed'] += 1
                    print(f"\033[31m Export Attachments Failed: {note.title}\033[0m")

        # 保存纯文本版本
        note_file_path = os.path.join(note_dir_path, "meta.txt")

        with open(note_file_path, 'w', encoding='utf-8') as f:
            # 元数据头
            f.write("=" * 60 + "\n")
            f.write(f"Note ID: {note_id}\n")
            print("*" * 30)
            print(note.creation_date)
            print(note.modification_date)
            print("*" * 30)
            if note.title:
                f.write(f"Title: {note.title}\n")
            if note.creation_date:
                f.write(f"Created: {self._format_date(note.creation_date)}\n")
            if note.modification_date:
                f.write(f"Modified: {self._format_date(note.modification_date)}\n")
            if note.created_iso:
                f.write(f"Created_ISO: {self._format_date(note.created_iso)}\n")
            if note.modified_iso:
                f.write(f"Modified_ISO: {self._format_date(note.modified_iso)}\n")
            if folder_name:
                f.write(f"Folder: {folder_name}\n")
            if note.attachments:
                f.write(f"Attachments: {len(note.attachments)} file(s)\n")
            f.write("=" * 60 + "\n\n")

            f.write("NOTE: Original formatted content is saved in note.html\n")
            f.write("      Raw data is in raw_data.bin and decompressed_data.bin\n")
            f.write("      Styles information is in styles.json\n\n")

            # f.write("---START_NOTE---\n")

            # 主要内容（使用包含附件标记的格式化文本）
            # if note.formatted_text:
            #     f.write(note.formatted_text)
            # else:
            #     f.write(note.text)
            #
            # f.write("\n---END_NOTE---\n")

            # 附件详细信息
            if note.attachments:
                f.write("\n\n" + "=" * 60 + "\n")
                f.write("Attachment Details:\n")
                for i, att in enumerate(note.attachments, 1):
                    f.write(f"\n{i}. Saved as: {att.filename or 'attachment'}\n")

                    if hasattr(att, 'original_filename') and att.original_filename:
                        f.write(f"   Original filename: {att.original_filename}\n")

                    if hasattr(att, 'source_filename') and att.source_filename:
                        f.write(f"   Source file: {att.source_filename}\n")

                    f.write(f"   Type: {att.type_uti}\n")
                    f.write(f"   Position: character {att.position}\n")
                    f.write(f"   Path: attachments/{att.filename}\n")
                    f.write(f"   ID: {att.identifier.split('-')[0] if att.identifier else 'unknown'}\n")

    def _sanitize_filename(self, name: str) -> str:
        """清理文件名"""
        name = re.sub(r'[<>:"/\\|?*\x00-\x1f]', '', name)
        if len(name) > 100:
            name = name[:100]
        return name.strip() or "unnamed"

    def _format_date(self, timestamp):
        """格式化时间戳"""
        if timestamp:
            try:
                apple_epoch = datetime(2001, 1, 1)
                dt = apple_epoch + timedelta(seconds=timestamp)
                return dt.strftime('%Y-%m-%d %H:%M:%S')
            except:
                return str(timestamp)
        return "unknown"

    def auto_delete_empty_dirs(self):
        """自动删除空目录"""
        for root, dirs, files in os.walk(self.output_dir, topdown=False):
            if not dirs and not files:
                try:
                    os.rmdir(root)
                    print(f"Deleted empty directory: {root}")
                except Exception as e:
                    print(f"Could not delete directory {root}: {e}")


    def get_curr_note_ids(self):
        return self._current_notes_ids

    def list_all_folders(self):
        """列出所有可用的文件夹（适配新版本 macOS 15+）"""
        conn = sqlite3.connect(self.db_path)
        try:
            cursor = conn.cursor()
            print("\nAnalyzing all folders in database...")
            print("=" * 60)

            # 新版本使用 Z_ENT=15 表示文件夹，名称在 ZTITLE2
            cursor.execute("""
                SELECT f.Z_PK as folder_id,
                       f.ZTITLE,
                       f.ZTITLE1,
                       f.ZTITLE2,  -- 新版本的文件夹名在这里！
                       f.ZNAME,
                       f.ZIDENTIFIER,
                       f.Z_ENT,
                       COUNT(DISTINCT n.Z_PK) as total_notes,
                       COUNT(DISTINCT nd.Z_PK) as notes_with_data
                FROM ZICCLOUDSYNCINGOBJECT f
                LEFT JOIN ZICCLOUDSYNCINGOBJECT n ON n.ZFOLDER = f.Z_PK
                LEFT JOIN ZICNOTEDATA nd ON nd.ZNOTE = n.Z_PK AND nd.ZDATA IS NOT NULL
                WHERE f.Z_ENT IN (14, 15)  -- 兼容旧版(14)和新版(15)
                GROUP BY f.Z_PK
                ORDER BY notes_with_data DESC
            """)

            folders = cursor.fetchall()

            if folders:
                print("Available folders:")
                for folder_id, title, title1, title2, name, identifier, z_ent, total, with_data in folders:
                    # 新版本优先使用 ZTITLE2
                    folder_name = title2 or title or title1 or name or f"Folder_{folder_id}"

                    print(f"\n  Folder ID {folder_id}: '{folder_name}'")
                    print(f"    • Total items: {total}")
                    print(f"    • Notes with data: {with_data}")
                    print(f"    • Entity type: {z_ent}")

                    # 特殊标记
                    if 'blog' in folder_name.lower():
                        print(f"    ⭐ This is your blog folder!")

                return folders
            else:
                print("No folders found.")
                return []

        except Exception as e:
            print(f"\nError: {e}")
            import traceback
            traceback.print_exc()
        finally:
            conn.close()


    def _print_stats(self):
        """打印统计信息"""
        print("\n" + "=" * 60)
        print("Export Statistics:")
        print(f"  Total notes: {self.stats['total']}")
        print(f"  Successfully exported: {self.stats['success']}")
        print(f"  Failed: {self.stats['failed']}")
        print(f"  Notes with attachments: {self.stats['with_attachments']}")
        print(f"  Attachments extracted: {self.stats['attachments_extracted']}")
        print(f"  Attachments not found: {self.stats['attachments_failed']}")
        print(f"  Notes with tables: {self.stats['with_tables']}")
        print(f"  Notes with checklists: {self.stats['with_checklists']}")
        print(f"\nOutput directory: {self.output_dir}")

#
def main():
    """主函数"""
    default_db = os.path.expanduser("~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite")

    if len(sys.argv) > 1:
        db_path = sys.argv[1]
    else:
        db_path = default_db
        print(f"Using default database: {db_path}")

    if len(sys.argv) > 2:
        output_dir = sys.argv[2]
    else:
        output_dir = "/tmp/apple_notes_export"

    try:
        exporter = AppleNotesExporter(db_path, output_dir)

        print(exporter.list_all_folders())

        exporter.export_all(folder_name='blog')
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()