# qianfan +connect 行为细节

## 命令

```
qianfan +connect [--target <dir>] [--dry-run] [--yes] [--format text|json]   # 安装/更新
qianfan +connect list [--target <dir>] [--format text|json]                  # 枚举（只读）
qianfan +connect uninstall [--target <dir>] [--dry-run] [--yes] [--format text|json]  # 卸载
```

## 目标目录探测

未显式 `--target` 时，按以下顺序选择目标技能目录：

1. `~/.comate/skills`（存在则使用）
2. `~/.claude/skills`（存在则使用）
3. 回退 `~/.qianfan/skills`

## 安装语义

- 内嵌技能包（全部 `qianfan-*`）先写入目标目录下的临时目录，再通过 rename **原子替换**同名目录。
- 幂等：重复安装等价于更新；同名目录被整体替换，旧内容不会残留。
- 替换失败时保留原目录，不会产生半安装状态。

## 预演与确认

- `--dry-run --format json` 输出执行计划：每个技能包对应 `create`（新增）或 `replace`（更新）动作及目标路径，
  不写入任何文件。
- 确认后加 `--yes` 执行。
- `--dry-run` 与 `--yes` 互斥（`INVALID_ARGUMENT`）；非交互环境未传 `--yes` 返回 `CONFIRMATION_REQUIRED`。

## 输出示例（--format json）

```json
{ "target": "/home/u/.comate/skills", "installed": ["qianfan-auth", "qianfan-chat", "qianfan-shared", "..."] }
```

`list` 返回 `installed`，`uninstall` 返回 `removed`。

## 安全约束

- `+connect list` 不联网、不需要凭证。
- 安装的技能包只编排 `qianfan` CLI，不读取本地凭证、不直接构造千帆 HTTP 请求、不输出敏感信息。
