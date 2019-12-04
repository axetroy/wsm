package db

import (
	"time"

	"github.com/axetroy/terminal/internal/library/util"
	"github.com/jinzhu/gorm"
)

type Team struct {
	Id      string  `gorm:"primary_key;not null;unique;index;type:varchar(32);" json:"id"` // 团队ID
	OwnerID string  `gorm:"not null;index;type:varchar(32);" json:"owner_id"`              // 拥有者
	Name    string  `gorm:"not null;type:varchar(32);" json:"name"`                        // 团队名
	Remark  *string `gorm:"null;type:varchar(36);" json:"remark"`                          // 备注

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (t *Team) TableName() string {
	return "team"
}

func (t *Team) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := util.GenerateId()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}
	return nil
}
