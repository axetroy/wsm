package db

import (
	"log"
	"time"

	"github.com/axetroy/terminal/internal/app/config"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

type TeamRole string

const (
	TeamRoleOwner   TeamRole = "owner"         // 拥有者, 获得所有权限
	TeamRoleAdmin   TeamRole = "administrator" // 管理员, 具有 增删改查服务器和连接服务器的权限
	TeamRoleMember  TeamRole = "member"        // 成员, 具有 查询和连接服务器的权限
	TeamRoleVisitor TeamRole = "visitor"       // 游客, 仅有查询服务器的权限, 连接服务器的权限都没有
)

var (
	teamMemberID *snowflake.Node
)

func init() {
	node, err := snowflake.NewNode(config.Common.MachineId)

	if err != nil {
		log.Panicln(err)
	}

	teamMemberID = node
}

type TeamMember struct {
	Id     string   `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // ID
	TeamID string   `gorm:"not null;index;type:varchar(32);" json:"team_id"`               // 团队 ID
	Team   Team     `gorm:"foreignkey:TeamID" json:"team"`                                 // ** 关联外键 **
	UserID string   `gorm:"not null;index;type:varchar(32);" json:"user_id"`               // 成员 ID
	User   User     `gorm:"foreignkey:UserID" json:"user"`                                 // ** 关联外键 **
	Role   TeamRole `gorm:"not null;type:varchar(32);" json:"role"`                        // 在团队中的角色
	Remark *string  `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (t *TeamMember) TableName() string {
	return "team_member"
}

func (t *TeamMember) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	id := teamMemberID.Generate().String()
	if err := scope.SetColumn("id", id); err != nil {
		return err
	}
	return nil
}

// 检查是否是可用的角色
func (t *TeamMember) AvailableRole() bool {
	return t.Role == TeamRoleOwner || t.Role == TeamRoleAdmin || t.Role == TeamRoleMember || t.Role == TeamRoleVisitor
}
