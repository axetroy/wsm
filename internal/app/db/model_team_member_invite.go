package db

import (
	"log"
	"time"

	"github.com/axetroy/terminal/internal/app/config"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

var (
	teamMemberInviteID *snowflake.Node
)

type InviteState string

const (
	InviteStateInit       InviteState = "init"       // 初始的邀请状态，既未接受，也未拒绝
	InviteStateAccept     InviteState = "accept"     // 接受邀请的状态
	InviteStateRefuse     InviteState = "refuse"     // 拒绝邀请的状态
	InviteStateCancel     InviteState = "cancel"     // 由团队管理者取消邀请
	InviteStateDeprecated InviteState = "deprecated" // 被系统弃用
)

func init() {
	node, err := snowflake.NewNode(config.Common.MachineId)

	if err != nil {
		log.Panicln(err)
	}

	teamMemberInviteID = node
}

type TeamMemberInvite struct {
	Id        string      `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // ID
	InvitorID string      `gorm:"not null;index;type:varchar(32);" json:"invitor_id"`            // 邀请人(发起者)
	Invitor   User        `gorm:"foreignkey:InvitorID" json:"invitor"`                           // 外建关联
	TeamID    string      `gorm:"not null;index;type:varchar(32);" json:"team_id"`               // 团队 ID
	Team      Team        `gorm:"foreignkey:TeamID" json:"team"`                                 // 外建关联
	UserID    string      `gorm:"not null;index;type:varchar(32);" json:"user_id"`               // 成员 ID
	User      User        `gorm:"foreignkey:UserID" json:"user"`                                 // 外建关联
	Role      TeamRole    `gorm:"not null;index;type:varchar(32);" json:"role"`                  // 在团队中的角色
	State     InviteState `gorm:"not null;index;type:varchar(32);" json:"state"`                 // 这条邀请记录的状态
	Remark    *string     `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (t *TeamMemberInvite) TableName() string {
	return "team_member_invite"
}

func (t *TeamMemberInvite) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	id := teamMemberInviteID.Generate().String()

	if err := scope.SetColumn("id", id); err != nil {
		return err
	}

	if err := scope.SetColumn("state", InviteStateInit); err != nil {
		return err
	}

	return nil
}
