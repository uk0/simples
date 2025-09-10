#!/usr/bin/env bash
set -e

echo "🧹 Cleaning old builds..."
rm -rf build dist *.build *.dist *.app .build
rm -rf "AppleNotesExporter.app"

echo "🔨 Building macOS App with PySide6..."

echo "🔨 Building with anti-bloat optimization..."

python -m nuitka \
    --standalone \
    --macos-create-app-bundle \
    --macos-app-name="AppleNotesExporter" \
    --macos-app-version="1.0.0" \
    --enable-plugin=pyside6 \
    --enable-plugin=anti-bloat \
    --noinclude-pytest-mode=nofollow \
    --noinclude-setuptools-mode=nofollow \
    --macos-target-arch=arm64 \
    --macos-app-icon=app_logo.icns \
    --include-qt-plugins=sensible,styles,platforms,iconengines,imageformats \
    --include-module=markdownify \
    --include-module=export_v2 \
    --include-module=helpers \
    --assume-yes-for-downloads \
    --show-progress \
    --output-dir=dist \
    AppleNotesExporter.py

echo "✅ Build complete!"

# 查找生成的 App
APP_PATH="dist/AppleNotesExporter.app"

if [ -d "$APP_PATH" ]; then
    echo "📦 App location: $APP_PATH"

    # 移除隔离属性
    xattr -cr "$APP_PATH"

    # 显示信息
    ls -lah "$APP_PATH"
    du -sh "$APP_PATH"

    # 打开 Finder
    open dist/
else
    echo "❌ App not found!"
    exit 1
fi