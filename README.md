# Qianfan CLI

Qianfan CLI 是千帆管控 API 与推理 API 的 Go 命令行客户端，同时内嵌一组 Agent Skills，
可通过 `qianfan +connect` 安装到本机 Agent（Comate / Claude Code 等），让 Agent 用自然语言调用千帆能力。

## 安装

```bash
npm i -g @baiducloud/qianfan-cli@latest
qianfan --version
```

### 直接下载二进制

也可以从 [GitHub Releases](https://github.com/baidubce/qianfan-cli/releases) 直接下载对应平台的二进制文件。

**推荐使用 `curl` 下载**（浏览器下载的文件在 macOS 上会被系统隔离，导致无法执行）：

```bash
# macOS (Apple Silicon)
curl -L -o qianfan "https://github.com/baidubce/qianfan-cli/releases/latest/download/qianfan-latest-darwin-arm64"
chmod +x qianfan && sudo mv qianfan /usr/local/bin/

# macOS (Intel)
curl -L -o qianfan "https://github.com/baidubce/qianfan-cli/releases/latest/download/qianfan-latest-darwin-amd64"
chmod +x qianfan && sudo mv qianfan /usr/local/bin/

# Linux (amd64)
curl -L -o qianfan "https://github.com/baidubce/qianfan-cli/releases/latest/download/qianfan-latest-linux-amd64"
chmod +x qianfan && sudo mv qianfan /usr/local/bin/
```

如果已通过浏览器下载，执行前需先移除系统隔离标记：

```bash
xattr -d com.apple.quarantine qianfan-*
chmod +x qianfan-*
```

> Windows 用户下载 `.exe` 文件后可直接运行，无需额外操作。

## 快速开始

```bash
qianfan auth login                                     # 浏览器 OAuth 登录
qianfan models list --format json                      # 查看可绑定模型
qianfan profile set-default --model-type chat <modelId> # 绑定默认 chat 模型
qianfan +chat --format json "你好"                      # 文本推理
qianfan service usage --startTime "2026-08-01 00:00:00" --endTime "2026-08-21 00:00:00"
qianfan doctor --format json                           # 体检
qianfan +connect                                       # 安装 Agent Skills
```

除 `profile use` 与 `profile set-default`（只输出一行确认文本）外，所有命令支持 `--format text|json`
（默认 `text`）；`auth logout`、`+connect`、`+connect uninstall` 会改动本地状态，支持 `--dry-run` 预演与
`--yes` 确认。

## 认证

`auth login` 通过系统浏览器完成 OAuth2/OIDC + PKCE 授权（Public Client、强制 PKCE、loopback 回调），
再用 ID Token 向 STS 换取临时凭证：Access Key 写入 `config.yaml`，Secret Key 与 session/refresh/id token
分别写入 `0600` 权限的凭证与会话文件。控制面请求用官方
[`github.com/baidubce/bce-sdk-go`](https://github.com/baidubce/bce-sdk-go) 的 BCE V1 signer 签名；
临时凭证临期或失效时由 refresh token 自动续期，无需重新登录。

CLI 同时只保存一个登录账号：已登录时再执行 `auth login` 返回 `ALREADY_LOGGED_IN`，换账号需先 `auth logout`。
本地账号名取自 ID Token 的 `principal_id`，凭证、会话与 API Key 缓存都按它分目录存放。

## 套餐与 API Key

- 套餐 profile：`platform`（平台预置服务，当前唯一可用）、`tokenplan_personal`、`tokenplan_enterprise`
  （已预置但尚未接入 API Key，切换会返回 `INVALID_ARGUMENT`）。
- 推理与模型列表使用 API Key（`Authorization: Bearer`）。platform 的 Key 由 CLI 经 openapi-iam 按固定
  名称 `APIKey-PermanentForQianfanCli` 查找，不存在则以全部权限创建，再 decrypt 取明文，登录时导入缓存。
- API Key 按 `account/profile` 缓存于 `~/.qianfan/api_keys`（`0600`），不设过期时间，生命周期与登录态一致；
  推理与管控请求遇 401/403 时在线重取一次并只重试一次——重取成功才覆盖缓存，失败保留原记录，
  避免一次网络或 IAM 故障连带作废本地仍可用的 Key。CLI 不会创建或伪造管控侧不存在的凭证。
- 沙盒环境的管控面请求会带 `run_mode: qa`（`qianfan.sandbox.baidu-int.com` 靠它区分 QA 流量）；
  生产环境与推理链路都不带这个头。

## 本地状态与环境变量

配置默认位于 `~/.qianfan`（`config.yaml`、`credentials/`、`sessions/`、`api_keys/`），敏感目录 `0700`、
文件 `0600`。

home 只能在登录时设置或重置：`auth login` 用 `QIANFAN_HOME` 指定状态目录（不设则沿用上一次登录的 home），
登录把最终路径记进固定位置的指针 `~/.qianfan/home.json`。其余指令不读该环境变量，只按 指针 → `~/.qianfan`
解析，因此永远落在登录态所在的那份 home 上。`auth logout` 不重置指针；要换回默认位置，用
`QIANFAN_HOME=~/.qianfan qianfan auth login` 重新登录一次。

可用环境变量：

- `QIANFAN_HOME`：仅 `auth login` 读取，用于指定本地状态根目录
- `QIANFAN_ENV`：`prod`（默认）/ `sandbox`，仅 `auth login` 读取。登录时固化进 config 的账号条目，此后所有命令按登录时记录的环境选域名；换环境需先 `auth logout` 再带新值登录
- `QIANFAN_CONTROL_PLANE_ENDPOINT` / `QIANFAN_INFERENCE_ENDPOINT` / `QIANFAN_CONSOLE_IAM_ENDPOINT` / `QIANFAN_OPENAPI_IAM_ENDPOINT` / `QIANFAN_STS_ENDPOINT`：分别覆盖管控面、推理、console-iam、openapi-iam 与 STS 域名（推理与管控面不同域，沙盒推理走独立调度服务）
- `QIANFAN_DEBUG=1`：向 stderr 输出 openapi-iam（取 API Key）请求的诊断信息，签名末段与令牌均脱敏

## Agent Skills

`skills/` 下的 `qianfan-*` 技能包通过 `go:embed` 打进二进制，`qianfan +connect` 会原子安装到探测到的
Agent 技能目录（`~/.comate/skills` → `~/.claude/skills` → `~/.qianfan/skills`）。技能包只编排 CLI 命令，
不读取本地凭证、不直接调用千帆 HTTP 接口。

## 开发

```bash
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go run ./cmd/qianfan --help
```

模块与扩展边界见 [docs/architecture.md](docs/architecture.md)。

## 发布

GitHub Release 由 `goreleaser.yaml` 构建关闭 CGO 的 darwin/linux/windows 二进制。npm 主包
`@baiducloud/qianfan-cli` 通过平台 optional dependency 启动对应二进制：GoReleaser 生成并解压各平台工件后，
执行 `scripts/package-npm.sh <version> <release-directory>` 生成主包与平台包。npm 安装过程不执行登录、
不访问千帆服务。
