package host

import (
	"errors"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (s *Service) AddCollaboratorToHostRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer helper.Response(&res, nil, nil, err)

	hostID := c.Param("host_id")
	userID := c.Param("collaborator_uid")

	res = s.AddCollaboratorToHost(controller.NewContextFromGinContext(c), hostID, userID)
}

func (s *Service) AddCollaboratorToHost(c controller.Context, hostID string, userID string) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		shouldUpdate bool
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, nil, nil, err)
	}()

	tx = db.Db.Begin()

	hostInfo := db.Host{Id: hostID, OwnerID: c.Uid}

	if err := tx.Where(&hostInfo).First(&hostInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	userInfo := db.User{Id: userID}

	if err := tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	hostRecordInfo := db.HostRecord{HostID: hostID, UserID: userID}

	if err := tx.Where(&hostRecordInfo).First(&hostRecordInfo).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
	} else {
		// 如果已是协作者，则返回错误
		err = exception.Duplicate
		return
	}

	// 加入协作者
	if err := tx.Create(db.HostRecord{
		HostID: hostID,
		UserID: userID,
		Type:   db.HostRecordTypeCollaborator,
	}).Error; err != nil {
		return
	}

	return
}
