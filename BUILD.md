# 构建指南

## 快速构建

```bash
# CLI（日常开发用这条即可）
go build -tags mutagencli -o build/mutagen ./cmd/mutagen/

# Agent
go build -tags mutagenagent -o build/mutagen-agent ./cmd/mutagen-agent/
```

构建产物统一放在 `build/` 目录（已被 `.gitignore` 忽略）。

## 构建标签（Build Tags）

| 标签 | 作用 | 必需场景 |
|------|------|----------|
| `mutagencli` | 注册内置传输协议（SSH/Docker/Local），解除 tagcheck panic 保护 | 构建独立 CLI 二进制 |
| `mutagenagent` | 解除 agent 的 tagcheck panic 保护 | 构建 agent 二进制 |
| `mutagensspl` | 启用 SSPL 许可的增强功能（xxh128 哈希、zstandard 压缩、Linux fanotify 监听） | 可选，需要额外功能时添加 |
| `mutagenfanotify` | 启用 Linux fanotify 文件系统监听（需配合 `mutagensspl`） | 可选，仅 Linux |
| `mutagensidecar` | 构建 sidecar 二进制 | 构建 sidecar |

**不加必需标签会怎样？** 每个二进制入口都有一个 `tagcheck.go` 守卫文件（构建约束为 `!mutagencli` 或 `!mutagenagent`），缺少标签时该文件会被编入，程序启动即 panic。这是有意设计——mutagen 的命令行代码可被第三方工具嵌入，此时第三方工具会注册自己的协议处理器而非内置的，tagcheck 防止嵌入场景下的二进制被误当作独立 CLI 运行。

## 带 SSPL 增强的构建

```bash
# CLI + SSPL 增强（xxh128 哈希、zstandard 压缩）
go build -tags mutagencli,mutagensspl -o build/mutagen ./cmd/mutagen/

# Agent + SSPL + fanotify（Linux）
go build -tags mutagenagent,mutagensspl,mutagenfanotify -o build/mutagen-agent ./cmd/mutagen-agent/
```

## 已知警告

macOS 上编译会出现以下警告，可安全忽略：

```
'FSEventStreamScheduleWithRunLoop' is deprecated: first deprecated in macOS 13.0
```

来自上游依赖 `github.com/mutagen-io/fsevents` 使用了已弃用的 macOS API，不影响功能。

## Windows 交叉编译

```bash
# 默认构建 amd64
./build/windows_cross_compile.sh

# 构建 arm64
GOARCH=arm64 ./build/windows_cross_compile.sh

# 自定义 agent bundle 中包含的平台
AGENT_TARGETS="linux/amd64 windows/amd64" ./build/windows_cross_compile.sh
```

产物在 `build/windows_<arch>/` 目录中：
- `mutagen.exe` — Windows CLI
- `mutagen-agents.tar.gz` — Agent 包（**必须**与 CLI 放在同一目录）

将整个目录复制到 Windows 即可使用。CLI 运行时会从同目录的 `mutagen-agents.tar.gz` 中按 `{goos}_{goarch}` 条目名提取对应平台的 agent 部署到远端主机。

## 官方构建脚本

项目自带的完整构建脚本位于 `scripts/build.go`，支持交叉编译、agent 打包、release 产物生成等。通过 `go run scripts/build.go` 执行，具体参数参见该文件。
