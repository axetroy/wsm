package db

import (
	"time"

	"github.com/axetroy/terminal/internal/library/util"
	"github.com/jinzhu/gorm"
)

type TeamMemberInvite struct {
	Id        string   `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // ID
	TeamID    string   `gorm:"not null;index;type:varchar(32);" json:"team_id"`               // 团队 ID
	Team      Team     `gorm:"foreignkey:TeamID" json:"team"`                                 // 外建关联
	UserID    string   `gorm:"not null;index;type:varchar(32);" json:"user_id"`               // 成员 ID
	Role      TeamRole `gorm:"not null;index;type:varchar(32);" json:"role"`                  // 在团队中的角色
	Available bool     `gorm:"not null;index;" json:"available"`                              // 当前邀请记录是否可用，当拒绝/接收邀请之后，之后就变成不可用
	Remark    *string  `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (t *TeamMemberInvite) TableName() string {
	return "team_member_invite"
}

func (t *TeamMemberInvite) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := util.GenerateId()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}