package db

import (
	"log"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

type HostRecordType string

const (
	HostRecordTypeOwner        HostRecordType = "owner"        // 拥有者
	HostRecordTypeCollaborator HostRecordType = "collaborator" // 协作者，可查看/连接服务器
	HostRecordTypeVisitor      HostRecordType = "visitor"      // 访客，仅能查看服务器
)

var (
	hostRecordID *snowflake.Node
)

func init() {
	node, err := snowflake.NewNode(config.Common.MachineId)

	if err != nil {
		log.Panicln(err)
	}

	hostRecordID = node
}

type HostRecord struct {
	Id     string         `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // 记录 ID
	UserID string         `gorm:"not null;index;type:varchar(32);" json:"user_id"`               // 对应的用户 ID
	HostID string         `gorm:"not null;index;type:varchar(32);" json:"host_id"`               // 对应的服务器 ID
	Host   Host           `gorm:"foreignkey:HostID" json:"host"`                                 // **关联外键**
	Type   HostRecordType `gorm:"not null;index;type:varchar(32);" json:"type"`                  // 类型

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *HostRecord) TableName() string {
	return "host_record"
}

func (u *HostRecord) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := hostRecordID.Generate().String()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}

	return nil
}
