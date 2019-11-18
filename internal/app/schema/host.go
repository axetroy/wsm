package schema

// 服务器的相关信息
type HostPure struct {
	Id       string  `json:"id"`
	OwnerID  string  `json:"owner_id"`
	Host     string  `json:"host"`
	Port     uint    `json:"port"`
	Username string  `json:"username"`
	Remark   *string `json:"remark"`
}

type Host struct {
	HostPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
