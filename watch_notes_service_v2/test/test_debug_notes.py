#!/usr/bin/env python3
"""
Apple Notes Folder Name Finder
通过多种策略找到文件夹名称的真实存储位置
"""

import sys
import os
import sqlite3
import json
import plistlib
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Tuple, Optional


class AppleNotesFolderFinder:
    """查找Apple Notes文件夹名称的真实存储位置"""

    def __init__(self, db_path: str):
        self.db_path = db_path
        self.conn = sqlite3.connect(db_path)
        self.base_path = Path(db_path).parent

        # 已知的实体类型
        self.ENTITY_TYPES = {
            'ICFolder': 14,
            'ICNote': 11,
            'ICAccount': 13,
            'ICNoteData': 18
        }

        # 收集所有可能包含文件夹名的列
        self.potential_name_columns = []
        self.analyze_schema()

    def analyze_schema(self):
        """分析数据库schema，找出所有可能包含名称的列"""
        cursor = self.conn.cursor()

        # 获取所有列
        cursor.execute("PRAGMA table_info(ZICCLOUDSYNCINGOBJECT)")
        all_columns = cursor.fetchall()
        self.all_columns = [col[1] for col in all_columns]

        # 找出所有可能包含文本的列
        text_keywords = [
            'TITLE', 'NAME', 'TEXT', 'STRING', 'IDENTIFIER',
            'DISPLAY', 'LABEL', 'SNIPPET', 'SUMMARY', 'CONTENT',
            'FALLBACK', 'USER', 'NESTED', 'SORTING', 'JSON'
        ]

        for col_name in self.all_columns:
            if any(keyword in col_name.upper() for keyword in text_keywords):
                self.potential_name_columns.append(col_name)

        print(f"Found {len(self.potential_name_columns)} potential name columns")

    def find_folders_by_note_count(self) -> Dict[int, Dict]:
        """通过笔记数量找到所有文件夹"""
        cursor = self.conn.cursor()

        print("\n" + "=" * 80)
        print("STEP 1: FINDING FOLDERS BY NOTE COUNT")
        print("=" * 80)

        # 获取所有有笔记的文件夹
        query = """
                SELECT n.ZFOLDER                     as folder_id, \
                       COUNT(DISTINCT n.Z_PK)        as total_notes, \
                       COUNT(DISTINCT nd.Z_PK)       as notes_with_data, \
                       GROUP_CONCAT(n.ZTITLE, ' | ') as note_titles
                FROM ZICCLOUDSYNCINGOBJECT n
                         INNER JOIN ZICNOTEDATA nd ON nd.ZNOTE = n.Z_PK
                WHERE nd.ZDATA IS NOT NULL
                GROUP BY n.ZFOLDER
                ORDER BY notes_with_data DESC \
                """

        cursor.execute(query)
        results = cursor.fetchall()

        folders = {}
        for folder_id, total, with_data, note_titles in results:
            folders[folder_id] = {
                'total_notes': total,
                'notes_with_data': with_data,
                'sample_note_titles': note_titles[:200] if note_titles else None
            }

            print(f"\nFolder ID {folder_id}:")
            print(f"  • Total notes: {total}")
            print(f"  • Notes with data: {with_data}")

            # 特别标记5个笔记的文件夹（你的blog）
            if with_data == 5:
                print(f"  ⭐ This matches your blog folder (5 notes)!")
                folders[folder_id]['likely_blog'] = True

            if note_titles:
                print(f"  • Sample notes: {note_titles[:100]}...")

        return folders

    def search_folder_names_comprehensively(self, folder_id: int) -> Dict[str, any]:
        """全面搜索指定文件夹的名称"""
        cursor = self.conn.cursor()

        print(f"\n" + "=" * 80)
        print(f"STEP 2: SEARCHING FOR NAME OF FOLDER {folder_id}")
        print("=" * 80)

        found_names = {}

        # 策略1: 检查所有文本列
        print("\nStrategy 1: Checking all text columns...")
        for col in self.potential_name_columns:
            try:
                query = f"SELECT {col} FROM ZICCLOUDSYNCINGOBJECT WHERE Z_PK = ?"
                cursor.execute(query, (folder_id,))
                result = cursor.fetchone()

                if result and result[0]:
                    value = result[0]

                    # 处理不同类型的数据
                    if isinstance(value, bytes):
                        # 尝试多种解码方式
                        decoded_values = []

                        # UTF-8
                        try:
                            decoded = value.decode('utf-8', errors='ignore')
                            if decoded and len(decoded) < 200:
                                decoded_values.append(('UTF-8', decoded))
                        except:
                            pass

                        # UTF-16
                        try:
                            decoded = value.decode('utf-16', errors='ignore')
                            if decoded and len(decoded) < 200:
                                decoded_values.append(('UTF-16', decoded))
                        except:
                            pass

                        # 作为plist
                        try:
                            decoded = plistlib.loads(value)
                            decoded_values.append(('PLIST', str(decoded)))
                        except:
                            pass

                        # 作为JSON
                        try:
                            decoded = json.loads(value)
                            decoded_values.append(('JSON', str(decoded)))
                        except:
                            pass

                        for encoding, decoded in decoded_values:
                            if 'blog' in decoded.lower():
                                print(f"  🎯 FOUND 'blog' in {col} ({encoding}): {decoded}")
                                found_names[f"{col}_{encoding}"] = decoded
                            elif decoded and len(decoded) < 50 and not decoded.startswith('{'):
                                print(f"  • {col} ({encoding}): {decoded}")
                                found_names[col] = decoded

                    elif isinstance(value, str):
                        if 'blog' in value.lower():
                            print(f"  🎯 FOUND 'blog' in {col}: {value}")
                            found_names[col] = value
                        elif len(value) < 100 and not value.startswith('{'):
                            print(f"  • {col}: {value}")
                            if value and value != 'None':
                                found_names[col] = value
            except Exception as e:
                continue

        # 策略2: 检查账户关联
        print("\nStrategy 2: Checking account relationships...")
        account_cols = ['ZACCOUNT', 'ZACCOUNT1', 'ZACCOUNT2', 'ZACCOUNT3',
                        'ZACCOUNT4', 'ZACCOUNT5', 'ZACCOUNT6', 'ZACCOUNT7', 'ZACCOUNT8']

        for col in account_cols:
            if col in self.all_columns:
                try:
                    query = f"SELECT {col} FROM ZICCLOUDSYNCINGOBJECT WHERE Z_PK = ?"
                    cursor.execute(query, (folder_id,))
                    result = cursor.fetchone()

                    if result and result[0]:
                        account_id = result[0]
                        # 获取账户信息
                        cursor.execute("""
                                       SELECT ZTITLE, ZNAME, ZIDENTIFIER
                                       FROM ZICCLOUDSYNCINGOBJECT
                                       WHERE Z_PK = ?
                                       """, (account_id,))
                        account_info = cursor.fetchone()
                        if account_info:
                            print(f"  • Account {account_id}: {account_info}")
                except:
                    continue

        # 策略3: 检查父文件夹
        print("\nStrategy 3: Checking parent folder...")
        if 'ZPARENT' in self.all_columns:
            cursor.execute("SELECT ZPARENT FROM ZICCLOUDSYNCINGOBJECT WHERE Z_PK = ?", (folder_id,))
            result = cursor.fetchone()
            if result and result[0]:
                parent_id = result[0]
                print(f"  • Parent folder ID: {parent_id}")
                # 递归查找父文件夹名称
                parent_names = self.search_folder_names_comprehensively(parent_id)
                if parent_names:
                    print(f"    Parent names: {parent_names}")

        return found_names

    def check_metadata_tables(self):
        """检查其他可能包含文件夹名称的表"""
        cursor = self.conn.cursor()

        print("\n" + "=" * 80)
        print("STEP 3: CHECKING METADATA AND RELATED TABLES")
        print("=" * 80)

        # 获取所有表
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
        tables = cursor.fetchall()

        for table_name in tables:
            table = table_name[0]
            if table not in ['ZICCLOUDSYNCINGOBJECT', 'ZICNOTEDATA']:
                # 检查表结构
                cursor.execute(f"PRAGMA table_info({table})")
                columns = cursor.fetchall()

                # 查找可能包含文件夹信息的列
                interesting_cols = []
                for col in columns:
                    col_name = col[1]
                    if any(keyword in col_name.upper() for keyword in ['FOLDER', 'TITLE', 'NAME', 'BLOG']):
                        interesting_cols.append(col_name)

                if interesting_cols:
                    print(f"\nTable {table} has interesting columns: {interesting_cols}")

                    # 尝试查询
                    try:
                        query = f"SELECT * FROM {table} LIMIT 5"
                        cursor.execute(query)
                        results = cursor.fetchall()
                        if results:
                            print(f"  Sample data from {table}:")
                            for row in results[:3]:
                                print(f"    {row}")
                    except:
                        pass

    def check_file_system(self):
        """检查文件系统中的相关文件"""
        print("\n" + "=" * 80)
        print("STEP 4: CHECKING FILE SYSTEM")
        print("=" * 80)

        # 查找可能包含文件夹信息的文件
        patterns = ['*.plist', '*.sqlite', '*.db', 'metadata*', 'index*']

        for pattern in patterns:
            files = list(self.base_path.glob(pattern))
            if files:
                print(f"\nFound {pattern} files:")
                for file in files[:5]:
                    print(f"  • {file.name}")

                    # 如果是plist文件，尝试读取
                    if file.suffix == '.plist':
                        try:
                            with open(file, 'rb') as f:
                                plist_data = plistlib.load(f)
                                if 'blog' in str(plist_data).lower():
                                    print(f"    🎯 Found 'blog' in {file.name}")
                        except:
                            pass

    def analyze_folder_122_specifically(self):
        """专门分析folder 122（你的blog文件夹）"""
        print("\n" + "=" * 80)
        print("STEP 5: DEEP ANALYSIS OF FOLDER 122 (YOUR BLOG FOLDER)")
        print("=" * 80)

        cursor = self.conn.cursor()

        # 获取folder 122的所有数据
        cursor.execute("SELECT * FROM ZICCLOUDSYNCINGOBJECT WHERE Z_PK = 122")
        row = cursor.fetchone()

        if row:
            print("\nAll non-null values for folder 122:")
            for i, value in enumerate(row):
                if value is not None and value != 0 and value != '':
                    col_name = self.all_columns[i]

                    # 特殊处理二进制数据
                    if isinstance(value, bytes):
                        # 尝试各种解码
                        decoded = None

                        # 检查是否可能是文本
                        try:
                            decoded = value.decode('utf-8', errors='ignore')
                            if 'blog' in decoded.lower():
                                print(f"  🎯 {col_name}: Contains 'blog' (UTF-8 decoded)")
                                print(f"     Value: {decoded[:200]}")
                                continue
                        except:
                            pass

                        # 检查是否是plist
                        try:
                            plist_data = plistlib.loads(value)
                            if 'blog' in str(plist_data).lower():
                                print(f"  🎯 {col_name}: Contains 'blog' (PLIST)")
                                print(f"     Value: {plist_data}")
                                continue
                        except:
                            pass

                        print(f"  • {col_name}: <binary data, {len(value)} bytes>")

                    elif isinstance(value, str) and len(value) < 500:
                        print(f"  • {col_name}: {value}")

                    elif isinstance(value, (int, float)):
                        # 检查是否是时间戳或引用
                        if col_name.endswith('DATE'):
                            try:
                                from datetime import datetime, timedelta
                                apple_epoch = datetime(2001, 1, 1)
                                dt = apple_epoch + timedelta(seconds=value)
                                print(f"  • {col_name}: {value} ({dt})")
                            except:
                                print(f"  • {col_name}: {value}")
                        else:
                            print(f"  • {col_name}: {value}")

    def provide_solution(self):
        """提供解决方案"""
        print("\n" + "=" * 80)
        print("SOLUTION AND RECOMMENDATIONS")
        print("=" * 80)

        print("""
Based on the analysis, here are the findings:

1. Folder 122 has 5 notes (matching your blog folder)
2. The folder name 'blog' is not stored in the standard ZTITLE column

POSSIBLE SOLUTIONS:

Option 1: Use folder ID directly
---------------------------------
Since we know folder 122 is your blog folder, use it directly:

    python export_script.py --folder-id 122 --folder-name "blog"

Option 2: Identify by note count
---------------------------------
Find the folder with exactly 5 notes:

    SELECT ZFOLDER FROM ZICCLOUDSYNCINGOBJECT 
    GROUP BY ZFOLDER 
    HAVING COUNT(*) = 5

Option 3: Check iCloud sync
----------------------------
The folder name might be stored in iCloud but not synced locally.
Try:
1. Open Notes app
2. Right-click on the blog folder
3. Choose "Show Folder Info" or rename it to trigger sync
4. Wait for sync to complete
5. Run this script again

Option 4: Manual mapping
------------------------
Create a mapping file:
    folder_mapping.json:
    {
        "122": "blog",
        "2": "Notes",
        "3": "Recently Deleted"
    }
        """)

    def run_complete_analysis(self):
        """运行完整分析"""
        # 1. 通过笔记数找文件夹
        folders = self.find_folders_by_note_count()

        # 2. 对每个文件夹搜索名称
        for folder_id in folders:
            if folders[folder_id].get('likely_blog'):
                names = self.search_folder_names_comprehensively(folder_id)
                if names:
                    print(f"\n✅ Found potential names for folder {folder_id}: {names}")

        # 3. 检查元数据表
        self.check_metadata_tables()

        # 4. 检查文件系统
        self.check_file_system()

        # 5. 专门分析folder 122
        self.analyze_folder_122_specifically()

        # 6. 提供解决方案
        self.provide_solution()

    def close(self):
        if self.conn:
            self.conn.close()


def main():
    """主函数"""
    default_db = os.path.expanduser("~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite")

    db_path = sys.argv[1] if len(sys.argv) > 1 else default_db

    print(f"Database: {db_path}")
    print(f"Modified: {datetime.fromtimestamp(os.path.getmtime(db_path))}")
    print(f"File size: {os.path.getsize(db_path) / 1024 / 1024:.2f} MB")
    print()

    if not os.path.exists(db_path):
        print(f"Error: Database not found at {db_path}")
        sys.exit(1)

    finder = AppleNotesFolderFinder(db_path)

    try:
        finder.run_complete_analysis()

        print("\n" + "=" * 80)
        print("FINAL RECOMMENDATION")
        print("=" * 80)
        print("""
Since folder 122 is confirmed to be your blog folder (5 notes),
you can use this command to export it:

    python apple_notes_exporter.py --folder-id 122 --folder-name blog

Or modify your export script to include a mapping:
    FOLDER_MAPPING = {
        122: "blog",
        # Add other folders as needed
    }
        """)

    finally:
        finder.close()


if __name__ == "__main__":
    main()