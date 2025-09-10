#!/usr/bin/env python3
import sys
import os
import threading
from pathlib import Path

from PyQt6.QtWidgets import (
    QApplication, QMainWindow, QWidget,
    QVBoxLayout, QHBoxLayout, QLabel,
    QLineEdit, QPushButton, QTextEdit,
    QFileDialog, QMessageBox, QSpacerItem,
    QSizePolicy, QCheckBox
)
from PyQt6.QtGui import QIcon, QColor, QFont, QPixmap, QPainter
from PyQt6.QtCore import Qt, QFileSystemWatcher, QTimer

from export import AppleNotesExporter  # 假设你已有的导出脚本模块


class StatusLight(QLabel):
    """一个圆形状态灯（红/绿）"""
    def __init__(self, diameter=12, parent=None):
        super().__init__(parent)
        self.diameter = diameter
        self.setFixedSize(diameter, diameter)
        self.set_off()

    def set_on(self, color: QColor = QColor("green")):
        pix = QPixmap(self.diameter, self.diameter)
        pix.fill(Qt.GlobalColor.transparent)
        painter = QPainter(pix)
        painter.setRenderHint(QPainter.RenderHint.Antialiasing)
        painter.setBrush(color)
        painter.setPen(Qt.PenStyle.NoPen)
        painter.drawEllipse(0, 0, self.diameter, self.diameter)
        painter.end()
        self.setPixmap(pix)

    def set_off(self):
        pix = QPixmap(self.diameter, self.diameter)
        pix.fill(Qt.GlobalColor.transparent)
        painter = QPainter(pix)
        painter.setRenderHint(QPainter.RenderHint.Antialiasing)
        painter.setBrush(QColor("red"))
        painter.setPen(Qt.PenStyle.NoPen)
        painter.drawEllipse(0, 0, self.diameter, self.diameter)
        painter.end()
        self.setPixmap(pix)


class ExportServiceWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("🍏 Apple Notes Exporter")
        self.resize(500, 480)

        # 文件系统监视器
        self.watcher = QFileSystemWatcher(self)
        self.watcher.fileChanged.connect(self._on_fs_changed)
        self.watcher.directoryChanged.connect(self._on_fs_changed)

        # 防抖定时器 (1s)
        self._debounce_timer = QTimer(self)
        self._debounce_timer.setSingleShot(True)
        self._debounce_timer.timeout.connect(lambda: self._run_export(auto=True))

        # UI 样式
        font = QFont("Helvetica Neue", 11)
        self.setFont(font)

        # 主容器
        container = QWidget()
        layout = QVBoxLayout()
        layout.setContentsMargins(14, 14, 14, 14)
        layout.setSpacing(12)
        container.setLayout(layout)
        self.setCentralWidget(container)

        # ========== App Name ==========
        title_label = QLabel("Apple Notes Exporter")
        title_label.setFont(QFont("Helvetica Neue", 16, QFont.Weight.Bold))
        title_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        layout.addWidget(title_label)

        # ========== Database Path ==========
        db_layout = QHBoxLayout()
        db_layout.setSpacing(8)
        db_layout.addWidget(QLabel("Database Path:"))
        self.db_input = QLineEdit()
        self.db_input.setPlaceholderText("~/Library/Group Containers/.../NoteStore.sqlite")
        self.db_input.textChanged.connect(self.on_db_path_changed)
        db_layout.addWidget(self.db_input, 1)

        self.db_status = StatusLight(14)
        db_layout.addWidget(self.db_status)

        btn_db_browse = QPushButton("Browse…")
        btn_db_browse.clicked.connect(self.browse_db)
        db_layout.addWidget(btn_db_browse)
        layout.addLayout(db_layout)

        # ========== Output Directory ==========
        out_layout = QHBoxLayout()
        out_layout.setSpacing(8)
        out_layout.addWidget(QLabel("Output Directory:"))
        self.out_input = QLineEdit()
        self.out_input.setPlaceholderText("apple_notes_export")
        out_layout.addWidget(self.out_input, 1)
        btn_out_browse = QPushButton("Browse…")
        btn_out_browse.clicked.connect(self.browse_output)
        out_layout.addWidget(btn_out_browse)
        layout.addLayout(out_layout)

        # ========== Folder Name ==========
        folder_layout = QHBoxLayout()
        folder_layout.setSpacing(8)
        folder_layout.addWidget(QLabel("Notes Folder:"))
        self.folder_input = QLineEdit()
        self.folder_input.setPlaceholderText("e.g. blog")
        self.folder_input.setText("Test")
        folder_layout.addWidget(self.folder_input, 1)
        layout.addLayout(folder_layout)

        # ========== Auto-detect Checkbox ==========
        self.chk_auto = QCheckBox("Auto-detect changes")
        self.chk_auto.setChecked(True)
        layout.addWidget(self.chk_auto)

        # ========== Buttons ==========
        btn_layout = QHBoxLayout()
        btn_layout.setSpacing(12)
        self.btn_list = QPushButton("List Folders")
        self.btn_list.clicked.connect(self.on_list_folders)
        btn_layout.addWidget(self.btn_list)
        btn_layout.addItem(QSpacerItem(20, 0, QSizePolicy.Policy.Expanding, QSizePolicy.Policy.Minimum))
        self.btn_export = QPushButton("Export Notes")
        self.btn_export.clicked.connect(self.on_export_notes)
        btn_layout.addWidget(self.btn_export)
        layout.addLayout(btn_layout)

        # ========== Log Output ==========
        self.log_output = QTextEdit()
        self.log_output.setReadOnly(True)
        self.log_output.setStyleSheet("""
            background-color: #2c2c2c;
            color: #ffffff;
            border: none;
            selection-background-color: #3c3c3c;
            font-family: Menlo, Monaco, 'Courier New', monospace;
        """)
        layout.addWidget(self.log_output, 1)

        # 默认值
        default_db = os.path.expanduser(
            "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"
        )
        self.db_input.setText(default_db)
        self.out_input.setText(str(Path.cwd() / "apple_notes_export"))

    def browse_db(self):
        path, _ = QFileDialog.getOpenFileName(
            self, "Select NoteStore.sqlite", str(Path.home()), "SQLite Files (*.sqlite *.db)"
        )
        if path:
            self.db_input.setText(path)

    def browse_output(self):
        path = QFileDialog.getExistingDirectory(
            self, "Select Output Directory", str(Path.cwd())
        )
        if path:
            self.out_input.setText(path)

    def log(self, message: str):
        """在日志框里输出信息，并滚动到最底部"""
        self.log_output.append(message)
        self.log_output.verticalScrollBar().setValue(
            self.log_output.verticalScrollBar().maximum()
        )

    def on_db_path_changed(self, path: str):
        """当数据库路径更新时，检测文件是否存在并更新状态灯、文件监视"""
        path = os.path.expanduser(path)
        exists = os.path.isfile(path)
        if exists:
            self.db_status.set_on(QColor("green"))
            self.log(f"🔗 Connected to database: {path}")
            # 重新监视文件及其父目录
            try:
                self.watcher.removePaths(self.watcher.files())
                self.watcher.removePaths(self.watcher.directories())
            except:
                pass
            self.watcher.addPath(path)
            self.watcher.addPath(os.path.dirname(path))
        else:
            self.db_status.set_off()
            self.log(f"⚠️ Database not found: {path}")

    def _on_fs_changed(self, path: str):
        """文件或目录改变时，防抖后自动执行导出"""
        if self.chk_auto.isChecked():
            self.log(f"📁 Change detected in {path}, scheduling export…")
            # 重启防抖定时器
            self._debounce_timer.start(1000)

    def on_list_folders(self):
        """列出所有文件夹"""
        db = self.db_input.text().strip()
        if not os.path.exists(db):
            QMessageBox.critical(self, "Error", "Database file not found")
            return
        self.log(f"🔍 Listing folders in {db} …")
        try:
            exporter = AppleNotesExporter(db, self.out_input.text().strip())
            folders = exporter.list_all_folders()
            self.log("📂 Available folders:")
            for fid, title, title1, title2, name, total, with_data in folders:
                display_name = title or title1 or title2 or name or f"Folder_{fid}"
                self.log(f"  • ID {fid}: '{display_name}' ({with_data} notes)")
        except Exception as e:
            QMessageBox.critical(self, "Error", str(e))

    def on_export_notes(self):
        """手动执行导出"""
        self._run_export(auto=False)

    def _run_export(self, auto: bool):
        """导出逻辑抽取到后台线程执行"""
        db = self.db_input.text().strip()
        out = self.out_input.text().strip()
        folder_name = self.folder_input.text().strip() or None

        if not os.path.exists(db):
            QMessageBox.critical(self, "Error", "Database file not found")
            return
        if not out:
            QMessageBox.critical(self, "Error", "Output directory is not set")
            return

        # 禁用按钮以防重复
        self.btn_export.setEnabled(False)
        self.btn_list.setEnabled(False)

        def task():
            try:
                exporter = AppleNotesExporter(db, out)
                exporter.export_all(folder_name=folder_name)
                self.log("✅ Export completed successfully.")
                if auto:
                    self.log("ℹ️ (Auto-export triggered)")
            except Exception as e:
                self.log(f"✗ Export failed: {e}")
            finally:
                self.btn_export.setEnabled(True)
                self.btn_list.setEnabled(True)

        threading.Thread(target=task, daemon=True).start()


def main():
    app = QApplication(sys.argv)
    window = ExportServiceWindow()
    window.setWindowIcon(QIcon.fromTheme("notes"))
    window.show()
    sys.exit(app.exec())


if __name__ == "__main__":
    main()