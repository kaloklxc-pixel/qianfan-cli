---
name: qianfan-connect
version: 1.0.0
description: "把 qianfan CLI 内嵌的 qianfan-* Agent 技能包安装到本机 Agent 技能目录（qianfan +connect），并支持枚举与卸载。当用户提到安装/更新千帆技能包、把千帆技能装到 Agent、qianfan +connect、列出或卸载已安装技能包，或希望本地 Agent 具备调用千帆 CLI 的能力时使用。反触发：用户想调用千帆能力本身（推理、用量、套餐）时转对应能力 skill，不必重新安装。"
license: Apache-2.0
metadata:
  requires:
    bins: ["qianfan"]
  cliHelp: "qianfan +connect --help"
  install:
    - kind: npm
      package: "@baiducloud/qianfan-cli"
---

# Qianfan 技能包安装

先读 [qianfan-common](../qianfan-common/SKILL.md) 的公共规则。`+connect` 只操作本地文件系统，
安装的技能包只编排 qianfan CLI，不读取凭证、不联网、不输出敏感信息。

## 何时使用本 skill

- 首次让本机 Agent 获得千帆能力
- CLI 升级后同步更新技能包
- 查看本机装了哪些 `qianfan-*` 技能包
- 清理不再需要的技能包

## 命令一览

| 命令 | 说明 |
| --- | --- |
| `qianfan +connect list [--target <dir>] [--format json]` | 枚举已安装技能包，只读、不联网、不需凭证 |
| `qianfan +connect [--target <dir>] [--dry-run] [--yes] [--format json]` | 安装/更新（同名目录原子替换），受保护 |
| `qianfan +connect uninstall [--target <dir>] [--dry-run] [--yes] [--format json]` | 卸载已安装技能包，受保护 |

## Agent 快速执行顺序

1. `qianfan +connect list --format json` 看目标目录与现状
2. 需要安装/更新 → `qianfan +connect --dry-run --format json` 展示将 `create`/`replace` 的目录
3. 用户确认后 → `qianfan +connect --yes`
4. 卸载同样先 `uninstall --dry-run --format json`，确认后 `--yes`

## 目标目录与安装语义

未显式 `--target` 时按顺序探测：`~/.comate/skills` → `~/.claude/skills` → 回退 `~/.qianfan/skills`。
安装为幂等的原子替换：内嵌技能包先写入临时目录再 rename 替换同名目录，替换失败保留原目录，
不会产生半安装状态。JSON 输出形如 `{"target": "...", "installed": ["qianfan-auth", ...]}`；
`list` 返回 `installed`，`uninstall` 返回 `removed`。更多细节见 [references/connect.md](references/connect.md)。

## 约束

- 安装/卸载属于本地变更：非交互环境未传 `--yes` 返回 `CONFIRMATION_REQUIRED`；
  `--dry-run` 与 `--yes` 不能同时使用（`INVALID_ARGUMENT`）；交互式拒绝确认返回 `OPERATION_CANCELLED`。
- `list` 是唯一无需确认的只读子命令。

## 与其他 skill 的串联

- 装完后要验证能力可用：转 [qianfan-doctor](../qianfan-doctor/SKILL.md)（`skills` 检查项）
- 尚未登录，装完也调不通：转 [qianfan-auth](../qianfan-auth/SKILL.md)

## 自然语言触发词 + 跨技能指引表

| 用户怎么说 | 走哪个命令/skill |
| --- | --- |
| "把千帆技能装到我的 Agent / 同步 skills" | `qianfan +connect`（先 dry-run） |
| "装了哪些千帆技能" | `qianfan +connect list --format json` |
| "卸载千帆技能包" | `qianfan +connect uninstall`（先 dry-run） |
| "装到别的目录" | 加 `--target <dir>` |
| "装完还是不能用" | 转 `qianfan-doctor` / `qianfan-auth` |

## 安全与边界

- 只写技能目录：不改凭证、不改 `~/.qianfan` 下的敏感文件，不联网。
- 覆盖前必须确认：`replace` 会整体替换同名目录，先用 `--dry-run` 让用户看清受影响路径。
- 不擅自换目录：`--target` 只在用户明确要求时使用，默认走探测结果。

## 参考

- `references/connect.md` — 目标目录探测、原子替换与幂等行为
- [qianfan-common/references/error-codes.md](../qianfan-common/references/error-codes.md) — 错误码、退出码与处置
