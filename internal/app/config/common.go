// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package config

import (
	"crypto/aes"

	"github.com/axetroy/wsm/internal/library/dotenv"
)

var (
	ModeProduction = "production"
)

type common struct {
	MachineId int64  `json:"machine_id"` // 机器 ID, 用于分布式生成 ID，每个节点的 ID 都应该不一样，并且最大值为 1024
	Mode      string `json:"mode"`       // 运行模式, 开发模式还是生产模式
	Exiting   bool   `json:"exiting"`    // 进程是否出于正在退出的状态，用户优雅的退出进程
	Secret    string `json:"secret"`     // 加密密钥
}

var Common *common

func init() {
	Common = &common{}
	Common.Mode = dotenv.GetByDefault("GO_MOD", ModeProduction)
	Common.MachineId = dotenv.GetInt64ByDefault("MACHINE_ID", 1)
	Common.Secret = dotenv.GetByDefault("SECRET", "astaxie12798akljzmknm.ahkjkljl;k")

	k := len(Common.Secret)
	switch k {
	default:
		panic(aes.KeySizeError(k))
	case 32:
		break
	}
}
