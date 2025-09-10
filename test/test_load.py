#!/usr/bin/env python3
"""
Complete Apple Notes Parser with Attachments Export - Fixed Version
"""

import sys
import os
import sqlite3
import gzip
import struct
import zlib
import re
import json
import shutil
from datetime import datetime, timedelta
from typing import Optional, List, Dict, Any, Tuple
from pathlib import Path


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


class Attachment:
    """表示一个附件"""

    def __init__(self):
        self.identifier = ""
        self.type_uti = ""
        self.filename = ""
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
        self.modification_date = None
        self.attachments = []
        self.tables = []
        self.checklists = []
        self.links = []
        self.styles = []
        self.raw_data = None

    def parse_from_blob(self, blob_data: bytes) -> bool:
        """从ZDATA blob解析笔记"""
        try:
            # 尝试gzip解压
            try:
                decompressed = gzip.decompress(blob_data)
            except:
                # 可能已经是解压的数据
                decompressed = blob_data

            self.raw_data = decompressed

            # 解析protobuf结构
            return self._parse_protobuf(decompressed)

        except Exception as e:
            print(f"Error parsing blob: {e}")
            # 尝试备用解析方法
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

            # Version (field 2)
            if field_number == 2 and wire_type == AppleNotesProtoParser.VARINT:
                version = parser.read_varint()
            # Note content (field 3)
            elif field_number == 3 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                note_data = parser.read_length_delimited()
                self._parse_note_content(note_data)
            else:
                parser.skip_field(wire_type)

    def _parse_note_content(self, data: bytes):
        """解析Note内容"""
        parser = AppleNotesProtoParser(data)

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
                self._parse_attribute_run(attr_data)

            else:
                parser.skip_field(wire_type)

    def _parse_attribute_run(self, data: bytes):
        """解析属性运行（格式信息）"""
        parser = AppleNotesProtoParser(data)
        attributes = {}

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Length (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.VARINT:
                attributes['length'] = parser.read_varint()

            # Paragraph style (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                style_data = parser.read_length_delimited()
                self._parse_paragraph_style(style_data)

            # Link (field 9)
            elif field_number == 9 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                link_data = parser.read_length_delimited()
                try:
                    link = link_data.decode('utf-8', errors='ignore')
                    self.links.append(link)
                except:
                    pass

            # Attachment (field 12)
            elif field_number == 12 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                attach_data = parser.read_length_delimited()
                attachment = self._parse_attachment(attach_data)
                if attachment:
                    self.attachments.append(attachment)

            else:
                parser.skip_field(wire_type)

        if attributes:
            self.styles.append(attributes)

    def _parse_paragraph_style(self, data: bytes):
        """解析段落样式"""
        parser = AppleNotesProtoParser(data)

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # Checklist (field 5)
            if field_number == 5 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                checklist_data = parser.read_length_delimited()
                self._parse_checklist(checklist_data)
            else:
                parser.skip_field(wire_type)

    def _parse_checklist(self, data: bytes):
        """解析清单项"""
        parser = AppleNotesProtoParser(data)
        checklist = {}

        while parser.pos < len(data):
            tag_info = parser.read_tag()
            if not tag_info:
                break

            field_number, wire_type = tag_info

            # UUID (field 1)
            if field_number == 1 and wire_type == AppleNotesProtoParser.LENGTH_DELIMITED:
                uuid_data = parser.read_length_delimited()
                checklist['uuid'] = uuid_data.hex()

            # Done status (field 2)
            elif field_number == 2 and wire_type == AppleNotesProtoParser.VARINT:
                checklist['done'] = bool(parser.read_varint())

            else:
                parser.skip_field(wire_type)

        if checklist:
            self.checklists.append(checklist)

    def _parse_attachment(self, data: bytes):
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

    def _fallback_parse(self, data: bytes) -> bool:
        """备用解析方法"""
        try:
            # 尝试直接gzip解压
            try:
                decompressed = gzip.decompress(data)
            except:
                decompressed = data

            # 尝试UTF-8解码
            text = decompressed.decode('utf-8', errors='ignore')

            # 跳过头部
            if len(text) > 12:
                text = text[12:]

            # 提取可打印字符
            result = ""
            for ch in text:
                if ch.isprintable() or ch.isspace() or ch == '￼':
                    result += ch
                else:
                    if len(result) > 10:  # 如果已经有一些内容，停止
                        break

            self.text = result
            return bool(self.text)

        except Exception as e:
            print(f"Fallback parse error: {e}")
            return False


class AttachmentExtractor:
    """附件提取器"""

    def __init__(self, db_path: str):
        self.db_path = db_path
        self.base_path = Path(db_path).parent
        self.conn = sqlite3.connect(db_path)

        # Apple Notes 附件可能的存储位置
        self.attachment_paths = [
            self.base_path / "Accounts",
            self.base_path / "Media",
            self.base_path / "FallbackImages",
            self.base_path / "Previews"
        ]

        # 检查数据库表结构
        self._check_database_schema()

    def _check_database_schema(self):
        """检查数据库表结构"""
        cursor = self.conn.cursor()

        # 获取ZICCLOUDSYNCINGOBJECT表的所有列
        cursor.execute("PRAGMA table_info(ZICCLOUDSYNCINGOBJECT)")
        columns = cursor.fetchall()
        self.available_columns = [col[1] for col in columns]

        print(f"    Available columns in database: {len(self.available_columns)} columns")

    def get_attachment_details(self, attachment_id: str) -> Dict[str, Any]:
        """从数据库获取附件详细信息"""
        cursor = self.conn.cursor()

        try:
            # 首先尝试基本查询
            basic_query = """
                          SELECT ZIDENTIFIER, \
                                 ZTYPEUTI, \
                                 ZFILENAME, \
                                 ZMEDIA, \
                                 ZCREATIONDATE, \
                                 ZMODIFICATIONDATE
                          FROM ZICCLOUDSYNCINGOBJECT
                          WHERE ZIDENTIFIER = ? \
                             OR ZIDENTIFIER LIKE ? \
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
                    'modification_date': row[5]
                }

                # 如果有media_id，尝试获取媒体文件信息
                if details['media_id']:
                    try:
                        # 查询媒体记录
                        media_query = """
                                      SELECT ZFILENAME, ZIDENTIFIER
                                      FROM ZICCLOUDSYNCINGOBJECT
                                      WHERE Z_PK = ? \
                                      """
                        media_results = cursor.execute(media_query, (details['media_id'],)).fetchall()
                        if media_results:
                            details['media_filename'] = media_results[0][0]
                            details['media_identifier'] = media_results[0][1]
                    except Exception as e:
                        print(f"      Warning: Could not get media info: {e}")

                return details

        except Exception as e:
            print(f"      Warning: Database query error: {e}")

        return {}

    def find_attachment_file(self, attachment_id: str, filename: str = None) -> Optional[Path]:
        """在文件系统中查找附件文件"""
        # 构建搜索模式
        patterns = []

        if filename:
            patterns.append(filename)
            # 去掉扩展名尝试
            name_without_ext = os.path.splitext(filename)[0]
            patterns.append(f"{name_without_ext}*")

        # 基于ID的搜索模式
        if attachment_id:
            # 清理ID（去掉可能的特殊字符）
            clean_id = re.sub(r'[^a-zA-Z0-9]', '', attachment_id)
            if clean_id:
                patterns.extend([
                    f"*{clean_id}*",
                    f"*{clean_id[:16]}*" if len(clean_id) > 16 else clean_id,
                    f"*{clean_id[:8]}*" if len(clean_id) > 8 else clean_id
                ])

        # 在所有可能的路径中搜索
        for base_path in self.attachment_paths:
            if not base_path.exists():
                continue

            # 搜索文件
            for pattern in patterns:
                try:
                    # 使用glob搜索
                    for file_path in base_path.rglob(pattern):
                        if file_path.is_file():
                            # 验证文件大小
                            if file_path.stat().st_size > 0:
                                return file_path
                except Exception as e:
                    continue

        # 尝试在整个Notes目录下搜索
        try:
            notes_dir = self.base_path
            if attachment_id:
                # 搜索包含ID前8位的文件
                short_id = attachment_id[:8] if len(attachment_id) > 8 else attachment_id
                for file_path in notes_dir.rglob(f"*{short_id}*"):
                    if file_path.is_file() and file_path.stat().st_size > 0:
                        return file_path
        except Exception:
            pass

        return None

    def extract_attachment(self, attachment: Attachment, output_dir: Path) -> bool:
        """提取并保存附件"""
        try:
            details = self.get_attachment_details(attachment.identifier)

            if not details and not attachment.identifier:
                print(f"      ✗ No attachment identifier")
                return False

            # 确定文件名
            filename = None
            if details.get('filename'):
                filename = details['filename']

            # 根据UTI生成文件名
            if not filename:
                ext = self._uti_to_extension(attachment.type_uti)
                base_name = attachment.identifier[:8] if attachment.identifier else "attachment"
                filename = f"{base_name}{ext}"

            # 确保文件名有扩展名
            if '.' not in filename:
                ext = self._uti_to_extension(attachment.type_uti)
                if ext:
                    filename = f"{filename}{ext}"

            # 清理文件名
            filename = self._sanitize_filename(filename)
            output_path = output_dir / filename

            # 尝试查找文件
            source_file = self.find_attachment_file(
                attachment.identifier,
                details.get('filename') or details.get('media_filename')
            )

            if source_file:
                shutil.copy2(source_file, output_path)
                attachment.filename = filename
                print(f"      ✓ Extracted: {filename} (from {source_file.name})")
                return True

            # 如果有media_identifier，尝试用它查找
            if details.get('media_identifier'):
                source_file = self.find_attachment_file(details['media_identifier'])
                if source_file:
                    shutil.copy2(source_file, output_path)
                    attachment.filename = filename
                    print(f"      ✓ Extracted: {filename} (via media ID)")
                    return True

            print(f"      ✗ Could not find file for: {filename}")

            # 保存附件信息到文本文件
            info_file = output_dir / f"{filename}.info.txt"
            with open(info_file, 'w') as f:
                f.write(f"Attachment Information\n")
                f.write(f"=====================\n")
                f.write(f"Identifier: {attachment.identifier}\n")
                f.write(f"Type: {attachment.type_uti}\n")
                if details:
                    for key, value in details.items():
                        f.write(f"{key}: {value}\n")

            return False

        except Exception as e:
            print(f"      ✗ Error extracting attachment: {e}")
            return False

    def _uti_to_extension(self, uti: str) -> str:
        """将UTI转换为文件扩展名"""
        if not uti:
            return ".dat"

        uti_lower = uti.lower()

        uti_map = {
            'public.jpeg': '.jpg',
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

        # 尝试从UTI提取扩展名
        if '.' in uti:
            parts = uti.split('.')
            last_part = parts[-1].lower()
            if len(last_part) <= 5 and last_part.isalnum():
                return f".{last_part}"

        return ".dat"

    def _sanitize_filename(self, filename: str) -> str:
        """清理文件名"""
        # 移除非法字符
        filename = re.sub(r'[<>:"/\\|?*\x00-\x1f]', '_', filename)
        # 限制长度
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

    def export_all(self):
        """导出所有笔记"""
        # 验证数据库
        if not os.path.exists(self.db_path):
            raise FileNotFoundError(f"Database not found: {self.db_path}")

        # 创建输出目录
        os.makedirs(self.output_dir, exist_ok=True)

        # 连接数据库
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # 查询笔记数据
        query = """
                SELECT nd.Z_PK             as note_id,
                       nd.ZDATA            as data,
                       n.ZTITLE            as title,
                       n.ZSNIPPET          as snippet,
                       n.ZCREATIONDATE     as created,
                       n.ZMODIFICATIONDATE as modified,
                       n.ZFOLDER           as folder_id,
                       f.ZTITLE            as folder_name
                FROM ZICNOTEDATA nd
                         LEFT JOIN ZICCLOUDSYNCINGOBJECT n ON nd.ZNOTE = n.Z_PK
                         LEFT JOIN ZICCLOUDSYNCINGOBJECT f ON n.ZFOLDER = f.Z_PK
                WHERE nd.ZDATA IS NOT NULL
                ORDER BY nd.Z_PK
                """

        results = cursor.execute(query).fetchall()
        self.stats['total'] = len(results)

        print(f"Found {self.stats['total']} notes to export")
        print("=" * 60)

        for row in results:
            self._export_note(row)

        conn.close()
        self.attachment_extractor.close()

        # 打印统计
        self._print_stats()

    def _export_note(self, row):
        """导出单个笔记"""
        note_id, data, title, snippet, created, modified, folder_id, folder_name = row

        print(f"\nProcessing Note #{note_id}")
        if title:
            print(f"  Title: {title}")
        if folder_name:
            print(f"  Folder: {folder_name}")

        # 解析笔记
        note = AppleNote()
        note.title = title or ""
        note.snippet = snippet or ""
        note.creation_date = created
        note.modification_date = modified

        if note.parse_from_blob(data):
            self._save_note(note_id, note, folder_name)
            self.stats['success'] += 1

            # 更新统计
            if note.attachments:
                self.stats['with_attachments'] += 1
            if note.tables:
                self.stats['with_tables'] += 1
            if note.checklists:
                self.stats['with_checklists'] += 1

            print(f"  ✓ Successfully exported")

            # 显示预览
            preview = note.text[:150].replace('\n', ' ')
            if len(note.text) > 150:
                preview += "..."
            print(f"  Preview: {preview}")
        else:
            self.stats['failed'] += 1
            print(f"  ✗ Failed to parse")

    def _save_note(self, note_id: int, note: AppleNote, folder_name: Optional[str]):
        """保存笔记到文件"""
        # 创建文件夹结构
        if folder_name:
            folder_path = os.path.join(self.output_dir, self._sanitize_filename(folder_name))
            os.makedirs(folder_path, exist_ok=True)
        else:
            folder_path = self.output_dir

        # 为每个笔记创建独立目录
        if note.title:
            note_dir_name = f"{note_id:04d} - {note.title}"
        elif note.text:
            first_line = note.text.split('\n')[0][:50]
            note_dir_name = f"{note_id:04d} - {first_line}"
        else:
            note_dir_name = f"{note_id:04d} - Note"

        note_dir_name = self._sanitize_filename(note_dir_name)
        note_dir_path = os.path.join(folder_path, note_dir_name)
        os.makedirs(note_dir_path, exist_ok=True)

        # 保存笔记内容
        note_file_path = os.path.join(note_dir_path, "note.txt")

        with open(note_file_path, 'w', encoding='utf-8') as f:
            # 元数据头
            f.write("=" * 60 + "\n")
            f.write(f"Note ID: {note_id}\n")
            if note.title:
                f.write(f"Title: {note.title}\n")
            if note.creation_date:
                f.write(f"Created: {self._format_date(note.creation_date)}\n")
            if note.modification_date:
                f.write(f"Modified: {self._format_date(note.modification_date)}\n")
            if folder_name:
                f.write(f"Folder: {folder_name}\n")
            f.write("=" * 60 + "\n\n")

            # 主要内容
            f.write(note.text)

            # 附加信息
            if note.links:
                f.write("\n\n" + "=" * 60 + "\n")
                f.write("Links:\n")
                for link in note.links:
                    f.write(f"  - {link}\n")

            # 附件信息
            if note.attachments:
                f.write("\n\n" + "=" * 60 + "\n")
                f.write(f"Attachments ({len(note.attachments)}):\n")

                # 创建附件目录
                attachments_dir = os.path.join(note_dir_path, "attachments")
                os.makedirs(attachments_dir, exist_ok=True)

                # 提取附件
                print(f"    Extracting {len(note.attachments)} attachment(s)...")
                for i, att in enumerate(note.attachments, 1):
                    success = self.attachment_extractor.extract_attachment(
                        att, Path(attachments_dir)
                    )

                    if success:
                        self.stats['attachments_extracted'] += 1
                        f.write(f"  {i}. {att.filename or 'attachment'}\n")
                        f.write(f"     Type: {att.type_uti}\n")
                        f.write(f"     ID: {att.identifier[:16] if att.identifier else 'unknown'}...\n")
                    else:
                        self.stats['attachments_failed'] += 1
                        f.write(f"  {i}. [Not found - see .info.txt file]\n")
                        f.write(f"     Type: {att.type_uti}\n")
                        f.write(f"     ID: {att.identifier[:16] if att.identifier else 'unknown'}...\n")

            if note.checklists:
                f.write("\n\n" + "=" * 60 + "\n")
                f.write("Checklists:\n")
                for item in note.checklists:
                    status = "✓" if item.get('done') else "○"
                    f.write(f"  {status} {item.get('uuid', 'item')}\n")

    def _sanitize_filename(self, name: str) -> str:
        """清理文件名"""
        # 移除非法字符
        name = re.sub(r'[<>:"/\\|?*\x00-\x1f]', '', name)
        # 限制长度
        if len(name) > 100:
            name = name[:100]
        return name.strip() or "unnamed"

    def _format_date(self, timestamp):
        """格式化时间戳"""
        if timestamp:
            try:
                # Apple的时间戳是从2001-01-01开始的
                apple_epoch = datetime(2001, 1, 1)
                dt = apple_epoch + timedelta(seconds=timestamp)
                return dt.strftime('%Y-%m-%d %H:%M:%S')
            except:
                return str(timestamp)
        return "unknown"

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


def main():
    """主函数"""
    # 默认路径
    default_db = os.path.expanduser("~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite")

    # 处理命令行参数
    if len(sys.argv) > 1:
        db_path = sys.argv[1]
    else:
        db_path = default_db
        print(f"Using default database: {db_path}")

    if len(sys.argv) > 2:
        output_dir = sys.argv[2]
    else:
        output_dir = "apple_notes_export"

    # 执行导出
    try:
        exporter = AppleNotesExporter(db_path, output_dir)
        exporter.export_all()
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()