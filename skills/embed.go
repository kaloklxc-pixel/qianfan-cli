// Package skills 通过 go:embed 将内嵌 Agent 技能包打包进 qianfan 二进制，
// 供 internal/connect 在本地安装。新增技能包目录后需在下方 embed 指令中登记。
package skills

import "embed"

// FS 保存所有随 CLI 分发的技能包目录。每个顶层目录对应一个技能包，
// 目录名以 qianfan- 开头。
//
//go:embed qianfan-common qianfan-shared qianfan-auth qianfan-profile qianfan-models qianfan-chat qianfan-service qianfan-doctor qianfan-connect
var FS embed.FS
