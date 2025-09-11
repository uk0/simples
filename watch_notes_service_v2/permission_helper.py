import os
import sys
import subprocess
import sqlite3
from pathlib import Path
from PySide6.QtWidgets import QMessageBox, QFileDialog, QApplication


class NotesAccessManager:
    """管理 Notes 数据库访问权限"""

    def __init__(self):
        self.notes_db_path = Path.home() / "Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"
        self.backup_db_path = None

    def check_access(self):
        """检查访问权限的多种方式"""
        # 方法1：直接尝试访问
        if self._try_direct_access():
            return True, "direct", str(self.notes_db_path)

        # 方法2：尝试使用备份路径
        if self._try_backup_locations():
            return True, "backup", str(self.backup_db_path)

        # 方法3：返回失败
        return False, "failed", None

    def _try_direct_access(self):
        """尝试直接访问数据库"""
        try:
            if not self.notes_db_path.exists():
                print(f"Database not found at: {self.notes_db_path}")
                return False

            # 尝试以只读方式打开
            conn = sqlite3.connect(f"file:{self.notes_db_path}?mode=ro", uri=True)
            cursor = conn.cursor()
            cursor.execute("SELECT name FROM sqlite_master WHERE type='table' LIMIT 1")
            result = cursor.fetchone()
            conn.close()

            if result:
                print(f"✅ Direct access successful: {self.notes_db_path}")
                return True
        except Exception as e:
            print(f"❌ Direct access failed: {e}")
        return False

    def _try_backup_locations(self):
        """尝试其他可能的位置"""
        backup_paths = [
            Path.home() / "Library/Containers/com.apple.Notes/Data/Library/Notes",
            Path.home() / "Library/Containers/com.apple.Notes.WidgetExtension",
            Path("/tmp/NoteStore.sqlite"),  # 如果用户复制到这里
        ]

        for path in backup_paths:
            if path.exists():
                db_files = list(path.glob("**/NoteStore.sqlite"))
                for db_file in db_files:
                    try:
                        conn = sqlite3.connect(str(db_file))
                        cursor = conn.cursor()
                        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' LIMIT 1")
                        if cursor.fetchone():
                            conn.close()
                            self.backup_db_path = db_file
                            print(f"✅ Found backup database: {db_file}")
                            return True
                        conn.close()
                    except:
                        continue
        return False


def request_full_disk_access_with_copy():
    """增强的权限请求，包含复制数据库选项"""
    msg = QMessageBox()
    msg.setIcon(QMessageBox.Icon.Warning)
    msg.setWindowTitle("需要访问权限")
    msg.setText("无法访问 Apple Notes 数据库")
    msg.setInformativeText(
        "请选择一种方式：\n\n"
        "1. 授予完全磁盘访问权限（推荐）\n"
        "2. 手动复制数据库到临时位置\n"
        "3. 使用文件选择器手动选择"
    )

    grant_btn = msg.addButton("授予权限", QMessageBox.ButtonRole.ActionRole)
    copy_btn = msg.addButton("复制数据库", QMessageBox.ButtonRole.ActionRole)
    select_btn = msg.addButton("手动选择", QMessageBox.ButtonRole.ActionRole)
    msg.addButton("退出", QMessageBox.ButtonRole.RejectRole)

    msg.exec()

    clicked = msg.clickedButton()

    if clicked == grant_btn:
        # 打开系统设置
        subprocess.run([
            "open",
            "x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles"
        ])

        QMessageBox.information(
            None,
            "设置权限",
            "请在系统设置中：\n"
            "1. 点击锁图标解锁\n"
            "2. 点击 + 号添加 AppleNotesExporter\n"
            "3. 重启应用"
        )
        sys.exit(0)

    elif clicked == copy_btn:
        # 提供复制命令
        copy_msg = QMessageBox()
        copy_msg.setWindowTitle("复制数据库")
        copy_msg.setText("请在终端中运行以下命令：")
        copy_msg.setInformativeText(
            'cp ~/Library/Group\\ Containers/group.com.apple.notes/NoteStore.sqlite /tmp/\n\n'
            '然后点击"已复制"继续'
        )

        copy_msg.addButton("已复制", QMessageBox.ButtonRole.AcceptRole)
        copy_msg.addButton("取消", QMessageBox.ButtonRole.RejectRole)

        copy_msg.exec()

        if copy_msg.result() == QMessageBox.StandardButton.Ok:
            # 检查 /tmp/NoteStore.sqlite
            tmp_db = Path("/tmp/NoteStore.sqlite")
            if tmp_db.exists():
                return str(tmp_db)

    elif clicked == select_btn:
        return get_notes_db_with_dialog()

    return None


def get_notes_db_with_dialog():
    """改进的文件选择对话框"""
    # 先显示提示
    instruction = QMessageBox()
    instruction.setWindowTitle("选择数据库文件")
    instruction.setText("请导航到 Notes 数据库位置")
    instruction.setInformativeText(
        "默认位置：\n"
        "~/Library/Group Containers/group.com.apple.notes/\n\n"
        "提示：\n"
        "1. 在打开的对话框中按 Cmd+Shift+G\n"
        "2. 粘贴上述路径\n"
        "3. 选择 NoteStore.sqlite 文件"
    )
    instruction.exec()

    # 尝试不同的起始路径
    start_paths = [
        Path.home() / "Library/Group Containers/group.com.apple.notes",
        Path.home() / "Library/Group Containers",
        Path.home() / "Library",
        Path.home(),
    ]

    start_path = Path.home()
    for path in start_paths:
        if path.exists():
            start_path = path
            break

    file_path, _ = QFileDialog.getOpenFileName(
        None,
        "选择 NoteStore.sqlite",
        str(start_path),
        "SQLite Database (*.sqlite);;All Files (*)"
    )

    if file_path and Path(file_path).exists():
        # 验证是否是有效的 Notes 数据库
        try:
            conn = sqlite3.connect(file_path)
            cursor = conn.cursor()
            cursor.execute(
                "SELECT name FROM sqlite_master WHERE type='table' AND name='ZICCLOUDSYNCINGOBJECT'"
            )
            if cursor.fetchone():
                conn.close()
                return file_path
            conn.close()

            QMessageBox.warning(None, "错误", "所选文件不是有效的 Notes 数据库")
        except Exception as e:
            QMessageBox.warning(None, "错误", f"无法打开数据库：{e}")

    return None


def check_and_request_access():
    """主要的权限检查和请求函数"""
    manager = NotesAccessManager()

    # 检查访问
    success, method, path = manager.check_access()

    if success:
        print(f"✅ 可以访问数据库 (方法: {method})")
        return path

    # 请求权限
    print("❌ 无法访问，请求权限...")
    result = request_full_disk_access_with_copy()

    if result:
        return result

    return None