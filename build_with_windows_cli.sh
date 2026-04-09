#!/usr/bin/env bash
#
# 执行官方 slim 构建（本机 CLI + 多平台 agent bundle），
# 然后额外交叉编译 Windows CLI。
#
# 用法（在仓库根目录执行）：
#   ./build/build_with_windows_cli.sh
#
# 产物（build/）：
#   mutagen                — 本机 CLI
#   mutagen-agents.tar.gz  — 多平台 Agent 包
#   mutagen.exe            — Windows CLI（amd64）
#
# 将 mutagen.exe + mutagen-agents.tar.gz 一起复制到 Windows 即可使用。
#
set -euo pipefail

WIN_ARCH="${WIN_ARCH:-amd64}"

echo "=== 执行官方 slim 构建 ==="
go run scripts/build.go

echo ""
echo "=== 交叉编译 Windows CLI (${WIN_ARCH}) ==="
CGO_ENABLED=0 GOOS=windows GOARCH="$WIN_ARCH" \
	go build -tags mutagencli -o "build/mutagen.exe" ./cmd/mutagen/

echo ""
echo "Build complete:"
ls -lh build/mutagen build/mutagen.exe build/mutagen-agents.tar.gz
echo ""
echo "本机使用: build/mutagen"
echo "Windows:  将 build/mutagen.exe 和 build/mutagen-agents.tar.gz 复制到同一目录"
