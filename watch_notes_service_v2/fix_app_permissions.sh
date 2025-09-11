#!/usr/bin/env bash
set -e

APP_PATH="dist/AppleNotesExporter.app"
PLIST_PATH="$APP_PATH/Contents/Info.plist"

echo "🔧 Fixing app permissions for Notes access..."

# 1. 修复 Bundle Identifier
/usr/libexec/PlistBuddy -c "Set :CFBundleIdentifier com.uk0.applenotesexporter" "$PLIST_PATH"

# 2. 添加所有必要的权限描述
/usr/libexec/PlistBuddy -c "Add :NSAppleEventsUsageDescription string 'This app needs to access Apple Notes data.'" "$PLIST_PATH" 2>/dev/null || true
/usr/libexec/PlistBuddy -c "Add :NSSystemAdministrationUsageDescription string 'This app requires Full Disk Access to read Notes database.'" "$PLIST_PATH" 2>/dev/null || true

# 3. 创建 entitlements 文件（关键！）
cat > entitlements.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <!-- 完全禁用沙盒 -->
    <key>com.apple.security.app-sandbox</key>
    <false/>

    <!-- 允许访问文件 -->
    <key>com.apple.security.files.user-selected.read-write</key>
    <true/>

    <!-- 允许访问用户目录 -->
    <key>com.apple.security.temporary-exception.files.home-relative-path.read-write</key>
    <array>
        <string>/Library/</string>
        <string>/Library/Group Containers/</string>
    </array>

    <!-- 禁用库验证 -->
    <key>com.apple.security.cs.disable-library-validation</key>
    <true/>

    <!-- 允许 DYLD 环境变量 -->
    <key>com.apple.security.cs.allow-dyld-environment-variables</key>
    <true/>
</dict>
</plist>
EOF

# 4. 移除所有扩展属性
xattr -cr "$APP_PATH"

# 5. 签名应用（非常重要的步骤顺序）
echo "🔐 Signing app with entitlements..."

# 先签名所有框架和库
find "$APP_PATH/Contents/MacOS" -type f -name "*.dylib" -o -name "*.so" | while read lib; do
    codesign --force --sign - "$lib" 2>/dev/null || true
done

# 签名主程序
codesign --force --deep \
    --sign - \
    --entitlements entitlements.plist \
    --options runtime \
    "$APP_PATH"

echo "✅ App permissions fixed!"

# 6. 验证
echo "📋 Verifying..."
codesign --verify --verbose "$APP_PATH"
codesign -d --entitlements - "$APP_PATH"

echo ""
echo "⚠️  重要："
echo "1. 第一次运行应用时，macOS 会提示安全警告"
echo "2. 在系统设置 > 隐私与安全性 > 完全磁盘访问权限 中添加应用"
echo "3. 重启应用即可正常使用"