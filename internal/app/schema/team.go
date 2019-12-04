package schema

// 团队信息相关
type TeamPure struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	OwnerID string  `json:"owner_id"`
	Remark  *string `json:"remark"`
}

type Team struct {
	TeamPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
