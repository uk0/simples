#!/usr/bin/env python3
import sys
import os
import json
import time
import queue
import hashlib
import threading
import subprocess
from pathlib import Path

# 现在导入 PySide6
from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget,
    QVBoxLayout, QHBoxLayout, QLabel,
    QLineEdit, QPushButton, QTextEdit,
    QFileDialog, QMessageBox, QCheckBox, QTabWidget, QSpacerItem, QSizePolicy
)
from PySide6.QtCore import Qt, QFileSystemWatcher, QTimer, QPoint, QSize
from PySide6.QtGui import QIcon, QFont
from PySide6.QtGui import QColor, QPixmap, QPainter

from permission_helper import get_notes_db_with_dialog, check_and_request_access

from export_v2 import AppleNotesExporter  # 假设你已有的导出脚本模块


class StatusLight(QLabel):
    """一个圆形状态灯（红/绿）"""
    def __init__(self, diameter=6, parent=None):
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
        self.resize(600, 580)
        # 文件系统监视器
        self.watcher = QFileSystemWatcher(self)
        self.watcher.fileChanged.connect(self._on_fs_changed)
        self.watcher.directoryChanged.connect(self._on_fs_changed)

        # 导出防抖定时器 (1s)
        self._debounce_timer = QTimer(self)
        self._debounce_timer.setSingleShot(True)
        self._debounce_timer.timeout.connect(lambda: self._run_export(auto=True))

        # 同步防抖定时器 (1s)
        self._rsync_timer = QTimer(self)
        self._rsync_timer.setSingleShot(True)
        self._rsync_timer.timeout.connect(lambda: self._run_rsync(auto=True))

        # 配置采集定时器 (1.2s) —— 主线程采集UI状态，交给后台线程落盘
        self._config_timer = QTimer(self)
        self._config_timer.setInterval(1200)
        self._config_timer.timeout.connect(self._enqueue_config_snapshot)

        # 配置写入线程通信
        self._cfg_queue: "queue.Queue[str]" = queue.Queue()
        self._cfg_writer_thread: threading.Thread | None = None
        self._cfg_last_hash: str | None = None

        # UI 样式
        font = QFont("Helvetica Neue", 11)
        self.setFont(font)
        #TODO 用来存储这一次与上一次的差异
        self.old_note_ids = []

        # 主容器和布局
        container = QWidget()
        main_layout = QVBoxLayout()
        main_layout.setContentsMargins(12, 12, 12, 12)
        main_layout.setSpacing(8)
        container.setLayout(main_layout)
        self.setCentralWidget(container)

        # Tab 控件
        self.tabs = QTabWidget()
        main_layout.addWidget(self.tabs, 1)

        # Export Tab
        self._build_export_tab()

        # Rsync Tab
        self._build_rsync_tab()

        # 日志输出区域（全局）
        self.log_output = QTextEdit()
        self.log_output.setReadOnly(True)
        self.log_output.setStyleSheet("""
            background-color: #2c2c2c;
            color: #ffffff;
            border: none;
            selection-background-color: #3c3c3c;
            font-family: Menlo, Monaco, 'Courier New', monospace;
        """)
        main_layout.addWidget(self.log_output, 1)

        # 默认路径
        default_db = Path.home() / "Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"

        if not default_db.exists():
            self.log("⚠️ Default Notes database not found. Please set the correct path.")
        self.log(f"⚙️ Using default database path {default_db.as_posix()}")
        self.db_input.setText(default_db.as_posix())
        self.out_input.setText(str("/tmp/apple_notes_export"))

        # 启动时尝试加载配置（会覆盖上面默认值）
        self._load_config_from_disk()

        # 根据当前 DB 路径更新 watcher 状态灯
        self.on_db_path_changed(self.db_input.text())

        # 启动配置自动保存线程
        self._start_config_autosave()

        self.check_and_setup_permissions()

    def check_and_setup_permissions(self):
        """启动时检查并设置权限"""
        db_path = check_and_request_access()

        if db_path:
            self.notes_db_path = db_path
            self.db_input.setText(str(db_path))
            self.log(f"✅ 使用数据库: {db_path}")
        else:
            self.log("❌ 无法获取数据库访问权限")

            # 提供最后的选择
            reply = QMessageBox.question(
                self,
                "无法访问",
                "无法访问 Notes 数据库。\n\n是否继续运行（功能受限）？",
                QMessageBox.StandardButton.Yes | QMessageBox.StandardButton.No
            )

            if reply == QMessageBox.StandardButton.No:
                sys.exit(0)

    # ---------------------------- 配置文件 ----------------------------

    def _config_path(self) -> Path:
        # 将配置存于用户目录，避免权限问题
        return Path.home() / ".apple_notes_exporter.json"

    def _snapshot_config(self) -> dict:
        # 仅在主线程调用：读取控件状态
        cfg = {
            "db_path": self.db_input.text().strip(),
            "old_note_ids": ','.join(self.old_note_ids),
            "out_dir": self.out_input.text().strip(),
            "folder_name": self.folder_input.text().strip(),
            "auto_export": bool(self.chk_auto.isChecked()),
            "rsync": {
                "auto": bool(self.chk_rsync_auto.isChecked()),
                "user": self.user_input.text().strip(),
                "ip": self.ip_input.text().strip(),
                "port": self.port_input.text().strip(),
                "remote": self.remote_input.text().strip(),
            },
            "window": {
                "width": int(self.size().width()),
                "height": int(self.size().height()),
                # 位置可能在不同显示器上有差异，尽量保存
                "x": int(self.pos().x()),
                "y": int(self.pos().y()),
            }
        }
        return cfg

    def _apply_config(self, cfg: dict):
        # 设置 UI 值（触发 textChanged 也没关系）
        if "db_path" in cfg and cfg["db_path"]:
            self.db_input.setText(str(cfg["db_path"]))

        if "old_note_ids" in cfg and cfg["old_note_ids"]:
            self.old_note_ids = str(cfg["old_note_ids"]).split(",")

        if "out_dir" in cfg and cfg["out_dir"]:
            self.out_input.setText(str(cfg["out_dir"]))
        if "folder_name" in cfg:
            self.folder_input.setText(str(cfg["folder_name"] or ""))

        if "auto_export" in cfg:
            self.chk_auto.setChecked(bool(cfg["auto_export"]))

        if "rsync" in cfg and isinstance(cfg["rsync"], dict):
            rs = cfg["rsync"]
            if "auto" in rs:
                self.chk_rsync_auto.setChecked(bool(rs["auto"]))
            if "user" in rs and rs["user"] is not None:
                self.user_input.setText(str(rs["user"]))
            if "ip" in rs and rs["ip"] is not None:
                self.ip_input.setText(str(rs["ip"]))
            if "port" in rs and rs["port"] is not None:
                self.port_input.setText(str(rs["port"]))
            if "remote" in rs and rs["remote"] is not None:
                self.remote_input.setText(str(rs["remote"]))

        if "window" in cfg and isinstance(cfg["window"], dict):
            w = cfg["window"]
            try:
                width = int(w.get("width", self.width()))
                height = int(w.get("height", self.height()))
                self.resize(QSize(width, height))
                x = int(w.get("x", self.pos().x()))
                y = int(w.get("y", self.pos().y()))
                # 仅当坐标合理时才移动
                if x >= 0 and y >= 0:
                    self.move(QPoint(x, y))
            except Exception:
                pass

    def _load_config_from_disk(self):
        cfg_file = self._config_path()
        if not cfg_file.exists():
            return
        try:
            data = json.loads(cfg_file.read_text(encoding="utf-8"))
            if isinstance(data, dict):
                self._apply_config(data)
                self.log(f"⚙️ Loaded config from {cfg_file}")
        except Exception as e:
            self.log(f"⚠️ Failed to load config {cfg_file}: {e}")

    def _start_config_autosave(self):
        # 后台写线程（仅负责落盘，不访问任何 Qt 控件）
        def writer():
            last_written_hash = None
            cfg_path = self._config_path()
            while True:
                try:
                    json_text = self._cfg_queue.get()  # 阻塞等待
                    if json_text is None:
                        break  # 预留退出
                    h = hashlib.sha256(json_text.encode("utf-8")).hexdigest()
                    if h == last_written_hash:
                        continue  # 内容未变化，跳过写入
                    # 原子写：先写临时文件，再替换
                    tmp_path = cfg_path.with_suffix(cfg_path.suffix + ".tmp")
                    tmp_path.parent.mkdir(parents=True, exist_ok=True)
                    tmp_path.write_text(json_text, encoding="utf-8")
                    os.replace(tmp_path, cfg_path)
                    last_written_hash = h
                    # 轻量日志节流：不每次都刷屏
                    # self.log(f"💾 Config saved to {cfg_path}")
                except Exception as e:
                    # 写入失败仅记录日志，不崩溃
                    print(f"[config-writer] error: {e}", file=sys.stderr)
                    time.sleep(0.5)

        self._cfg_writer_thread = threading.Thread(target=writer, daemon=True)
        self._cfg_writer_thread.start()

        # 启动主线程定时器，定期采集配置并入队
        self._config_timer.start()

    def _enqueue_config_snapshot(self):
        try:
            cfg = self._snapshot_config()
            json_text = json.dumps(cfg, ensure_ascii=False, sort_keys=True, indent=2)
            h = hashlib.sha256(json_text.encode("utf-8")).hexdigest()
            if h != self._cfg_last_hash:
                self._cfg_last_hash = h
                # 不阻塞 UI；队列满时丢弃旧数据
                try:
                    self._cfg_queue.put_nowait(json_text)
                except queue.Full:
                    pass
        except Exception as e:
            # 捕获异常以免影响 UI 循环
            self.log(f"⚠️ Snapshot config failed: {e}")

    # ---------------------------- UI 构建 ----------------------------

    def _build_export_tab(self):
        """构建导出面板"""
        export_tab = QWidget()
        layout = QVBoxLayout()
        layout.setSpacing(10)
        export_tab.setLayout(layout)

        # App 名称
        title_label = QLabel("Apple Notes Export")
        title_label.setFont(QFont("Helvetica Neue", 16, QFont.Weight.Bold))
        title_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        layout.addWidget(title_label)

        # Database Path
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

        # Output Directory
        out_layout = QHBoxLayout()
        out_layout.setSpacing(8)
        out_layout.addWidget(QLabel("Output Directory:"))
        self.out_input = QLineEdit()
        self.out_input.setPlaceholderText("/tmp/apple_notes_export")
        out_layout.addWidget(self.out_input, 1)
        btn_out_browse = QPushButton("Browse…")
        btn_out_browse.clicked.connect(self.browse_output)
        out_layout.addWidget(btn_out_browse)
        layout.addLayout(out_layout)

        # Folder Name
        folder_layout = QHBoxLayout()
        folder_layout.setSpacing(8)
        folder_layout.addWidget(QLabel("Notes Folder:"))
        self.folder_input = QLineEdit()
        self.folder_input.setPlaceholderText("e.g. blog")
        self.folder_input.setText("blog")
        folder_layout.addWidget(self.folder_input, 1)
        layout.addLayout(folder_layout)

        # Auto-detect Checkbox
        self.chk_auto = QCheckBox("Auto-export on DB change")
        self.chk_auto.setChecked(True)
        layout.addWidget(self.chk_auto)

        # Buttons
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

        self.tabs.addTab(export_tab, "Export")

    def _build_rsync_tab(self):
        """构建Rsync面板"""
        sync_tab = QWidget()
        layout = QVBoxLayout()
        layout.setSpacing(8)
        sync_tab.setLayout(layout)

        title_label = QLabel("Rsync Synchronization")
        title_label.setFont(QFont("Helvetica Neue", 16, QFont.Weight.Bold))
        title_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        layout.addWidget(title_label)

        # SSH User
        user_layout = QHBoxLayout()
        user_layout.addWidget(QLabel("SSH User:"))
        self.user_input = QLineEdit()
        self.user_input.setPlaceholderText("root")
        self.user_input.setText("root")
        user_layout.addWidget(self.user_input, 1)
        layout.addLayout(user_layout)

        # Server IP
        ip_layout = QHBoxLayout()
        ip_layout.addWidget(QLabel("Server IP:"))
        self.ip_input = QLineEdit()
        self.ip_input.setPlaceholderText("e.g. 192.168.1.100")
        ip_layout.addWidget(self.ip_input, 1)
        layout.addLayout(ip_layout)

        # Port
        port_layout = QHBoxLayout()
        port_layout.addWidget(QLabel("SSH Port:"))
        self.port_input = QLineEdit()
        self.port_input.setPlaceholderText("e.g. 22")
        port_layout.addWidget(self.port_input, 1)
        layout.addLayout(port_layout)

        # Remote Path
        remote_layout = QHBoxLayout()
        remote_layout.addWidget(QLabel("Remote Path:"))
        self.remote_input = QLineEdit()
        self.remote_input.setPlaceholderText("e.g. /home/user/notes_export")
        remote_layout.addWidget(self.remote_input, 1)
        layout.addLayout(remote_layout)

        # Auto-Rsync Checkbox
        self.chk_rsync_auto = QCheckBox("Auto-rsync after export")
        self.chk_rsync_auto.setChecked(False)
        layout.addWidget(self.chk_rsync_auto)

        # Rsync Button
        self.btn_rsync = QPushButton("Sync Now")
        self.btn_rsync.clicked.connect(lambda: self._run_rsync(auto=False))
        layout.addWidget(self.btn_rsync)

        self.tabs.addTab(sync_tab, "Sync")

    # ---------------------------- 交互逻辑 ----------------------------

    def browse_db(self):
        path = get_notes_db_with_dialog()
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
        path = os.path.expanduser(path or "")
        exists = os.path.isfile(path)
        if exists:
            self.db_status.set_on(QColor("green"))
            self.log(f"🔗 Connected to database: {path}")
            try:
                self.watcher.removePaths(self.watcher.files())
                self.watcher.removePaths(self.watcher.directories())
            except Exception:
                pass
            # 监听 DB 文件和其所在目录
            self.watcher.addPath(path)
            self.watcher.addPath(os.path.dirname(path))
        else:
            self.db_status.set_off()
            if path:
                self.log(f"⚠️ Database not found: {path}")

    def _on_fs_changed(self, path: str):
        if self.chk_auto.isChecked():
            self.log(f"📁 Change detected in {path}, scheduling export…")
            self._debounce_timer.start(3000)

    def on_list_folders(self):
        db = self.db_input.text().strip()
        if not os.path.exists(db):
            QMessageBox.critical(self, "Error", "Database file not found")
            return
        self.log(f"🔍 Listing folders in {db} …")
        try:
            exporter = AppleNotesExporter(db, self.out_input.text().strip())
            folders = exporter.list_all_folders()
            self.log("📂 Available folders:")
            for folder_id, title, title1, title2, name, identifier, z_ent, total, with_data  in folders:
                display_name = title or title1 or title2 or name or f"Folder_{folder_id}"
                self.log(f"  • ID {folder_id}: '{display_name}' ({with_data} notes)")
        except Exception as e:
            QMessageBox.critical(self, "Error", str(e))

    def on_export_notes(self):
        self._run_export(auto=False)

    def _run_export(self, auto: bool):
        db = self.db_input.text().strip()
        out = self.out_input.text().strip()
        folder_name = self.folder_input.text().strip() or None
        if not os.path.exists(db):
            QMessageBox.critical(self, "Error", "Database file not found")
            return
        if not out:
            QMessageBox.critical(self, "Error", "Output directory is not set")
            return

        self.btn_export.setEnabled(False)
        self.btn_list.setEnabled(False)
        self.btn_rsync.setEnabled(False)

        def task():
            try:
                exporter = AppleNotesExporter(db, out)
                exporter.export_all(folder_name=folder_name)

                if len(self.old_note_ids)==0:
                    self.old_note_ids = exporter.get_curr_note_ids()
                else:
                    new_ids = exporter.get_curr_note_ids()
                    added = [nid for nid in new_ids if nid not in self.old_note_ids]
                    removed = [oid for oid in self.old_note_ids if oid not in new_ids]
                    if added:
                        self.log(f"🆕 New notes added: {added}")
                    if removed:
                        self.log(f"🗑️ Notes removed: {removed}")
                        self.auto_delete_files(removed, out)
                    self.old_note_ids = new_ids
                self.log("✅ Export completed successfully.")
                if auto:
                    self.log("ℹ️ (Auto-export triggered)")
                if self.chk_rsync_auto.isChecked():
                    self.log("🔄 Auto-rsync triggered")
                    self._rsync_timer.start(500)

                self._run_rsync(auto=False)
                self.log("🔄 Auto-rsync triggered")
            except Exception as e:
                self.log(f"✗ Export failed: {e}")
            finally:
                self.btn_export.setEnabled(True)
                self.btn_list.setEnabled(True)
                self.btn_rsync.setEnabled(True)

        threading.Thread(target=task, daemon=True).start()

    def _run_rsync(self, auto: bool):
        user = self.user_input.text().strip() or "root"
        ip = self.ip_input.text().strip()
        port = self.port_input.text().strip() or "22"
        remote = self.remote_input.text().strip()
        out = self.out_input.text().strip()

        if not ip or not remote:
            QMessageBox.critical(self, "Error", "Server IP and Remote Path are required")
            return

        self.btn_rsync.setEnabled(False)
        self.log(f"🔄 Starting rsync to {user}@{ip}:{remote} …")

        def task():
            try:
                cmd = [
                    "rsync","--delete" , "-avz", "-e" ,
                    f"ssh -p {port}",
                    f"{out}",
                    f"{user}@{ip}:{remote}/"
                ]
                self.log(f"· Command: {' '.join(cmd)}")
                result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)
                if result.returncode == 0:
                    self.log("✅ Rsync completed successfully.")
                    if auto:
                        self.log("ℹ️ (Auto-rsync triggered)")
                else:
                    self.log(f"✗ Rsync error (code {result.returncode}):\n{result.stderr.strip()}")
            except Exception as e:
                self.log(f"✗ Rsync exception: {e}")
            finally:
                self.btn_rsync.setEnabled(True)

        threading.Thread(target=task, daemon=True).start()

    def auto_delete_files(self, removed, out):
        """自动删除导出目录中已被删除的笔记文件"""

        if out!="/tmp/apple_notes_export":
            self.log("⚠️ Auto-delete skipped: Output directory is not /tmp/apple_notes_export")
            return
        """
        自动删除导出目录中已被删除的笔记内容。
        顺序：
          1) 先删文件
          2) 再删该目录下的子目录
          3) 最后删除 note_id 那一层的目录

        参数:
          removed: 可迭代对象，元素形如 "p_dir_noteid"
          out    : 导出根目录
        """
        try:
            base_out = Path(out).resolve()

            for note_mx in removed:
                # 仅按第一个下划线分割，允许 p_dir 中包含下划线
                try:
                    p_dir, note_id = note_mx.split("_", 1)
                except ValueError:
                    self.log(f"⚠️ Skip malformed id: {note_mx}")
                    continue

                # 构造期望目录：<out>/<p_dir>/<note_id>
                note_root = (base_out / p_dir / note_id)

                # 必须存在，否则跳过
                if not note_root.exists():
                    continue

                # 解析真实路径（跟随符号链接），并确保严格在 base_out 之下
                try:
                    resolved = note_root.resolve()
                except FileNotFoundError:
                    # 被并发删除或指向不存在目标
                    continue
                except Exception as e:
                    self.log(f"⚠️ Resolve failed: {note_root} ({e})")
                    continue

                # 禁止删除导出根目录本身
                if resolved == base_out:
                    self.log(f"⛔ Refuse to delete base output directory: {resolved}")
                    continue

                # 要求 resolved 必须是 base_out 的子路径（不能是同级/父级/跨越）
                try:
                    rel = resolved.relative_to(base_out)
                except ValueError:
                    self.log(f"⚠️ Unsafe path (outside out): {resolved}")
                    continue

                # 期望相对路径深度恰好为两级：<p_dir>/<note_id>
                if len(rel.parts) != 2 or rel.parts[0] != p_dir or rel.parts[1] != note_id:
                    self.log(f"⚠️ Unexpected path depth for {resolved}, skip")
                    continue

                # 自底向上删除内容：先文件、再子目录、最后根目录
                for root, dirs, files in os.walk(resolved, topdown=False, followlinks=False):
                    # 1) 删除文件
                    for fname in files:
                        fpath = Path(root) / fname
                        try:
                            fpath.unlink(missing_ok=True)
                            self.log(f"🗑️ Deleted exported file: {fpath}")
                        except Exception as e:
                            self.log(f"⚠️ Failed to delete file {fpath}: {e}")

                    # 2) 删除子目录（若是符号链接则只删除链接本身）
                    for dname in dirs:
                        dpath = Path(root) / dname
                        try:
                            if dpath.is_symlink():
                                dpath.unlink(missing_ok=True)
                                self.log(f"🗑️ Deleted symlink dir: {dpath}")
                            else:
                                dpath.rmdir()  # 若为空则成功
                                self.log(f"🗑️ Deleted exported dir: {dpath}")
                        except Exception as e:
                            self.log(f"⚠️ Failed to delete dir {dpath}: {e}")

                # 3) 删除 note_id 根目录
                try:
                    if resolved.is_symlink():
                        resolved.unlink(missing_ok=True)
                        self.log(f"🗑️ Deleted symlink root: {resolved}")
                    else:
                        resolved.rmdir()
                        self.log(f"🗑️ Deleted exported dir: {resolved}")
                except Exception as e:
                    self.log(f"⚠️ Failed to delete root dir {resolved}: {e}")


        except Exception as e:
            self.log(f"⚠️ Auto-delete error: {e}")


def main():
    app = QApplication(sys.argv)
    window = ExportServiceWindow()
    window.setWindowIcon(QIcon.fromTheme("notes"))
    window.show()
    sys.exit(app.exec())


if __name__ == "__main__":
    main()