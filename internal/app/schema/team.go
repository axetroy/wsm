package schema

import "github.com/axetroy/terminal/internal/app/db"

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

type TeamWithMember struct {
	TeamPure
	Owner     ProfilePublic `json:"owner"`   // 团队拥有者的基本信息
	UserID    string        `json:"user_id"` // 用户 ID
	Role      db.TeamRole   `json:"role"`    // 用户在团队中扮演的角色
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

type TeamStat struct {
	Team
	MemberNum int `json:"member_num"` // 团队的成员数量
	HostNum   int `json:"host_num"`   // 拥有的服务器数量
}

// 团队成员信息
type TeamMember struct {
	ProfilePublic
	Role      db.TeamRole `json:"role"`       // 用户在团队的角色
	CreatedAt string      `json:"created_at"` // 用户加入团队的时间
}

type TeamMemberInvitePure struct {
	ID     string         `json:"id"`
	Role   db.TeamRole    `json:"role"`
	State  db.InviteState `json:"state"`
	Remark *string        `json:"remark"`
}

type TeamMemberInvite struct {
	TeamMemberInvitePure
	User      ProfilePublic `json:"user"`
	Team      TeamPure      `json:"team"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}
