// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/terminal/internal/library/dotenv"
)

var (
	ModeProduction = "production"
)

type common struct {
	MachineId string `json:"machine_id"` // 机器 ID
	Mode      string `json:"mode"`       // 运行模式, 开发模式还是生产模式
	Signature string `json:"signature"`  // 签名密钥，主要用户签名数据
	Exiting   bool   `json:"exiting"`    // 进程是否出于正在退出的状态，用户优雅的退出进程
}

var Common *common

func init() {
	Common = &common{}
	Common.Mode = dotenv.GetByDefault("GO_MOD", ModeProduction)
	Common.MachineId = dotenv.GetByDefault("MACHINE_ID", "0")
	Common.Signature = dotenv.GetByDefault("SIGNATURE_KEY", "signature key")
}
