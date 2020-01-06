package schema

import "github.com/axetroy/wsm/internal/app/db"

// 服务器的相关信息
type HostPure struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	OwnerID     string             `json:"owner_id"`
	OwnerType   db.HostOwnerType   `json:"owner_type"`
	Host        string             `json:"host"`
	Port        uint               `json:"port"`
	Username    string             `json:"username"`
	ConnectType db.HostConnectType `json:"connect_type"`
	Remark      *string            `json:"remark"`
}

type Host struct {
	HostPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
