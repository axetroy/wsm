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

func (s *Service) TransferHostRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer helper.Response(&res, nil, nil, err)

	hostID := c.Param("host_id")
	userID := c.Param("user_id")

	res = s.TransferHost(controller.NewContextFromGinContext(c), hostID, userID)
}

func (s *Service) TransferHost(c controller.Context, hostID string, userID string) (res schema.Response) {
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

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	// 更新 owner ID
	if err = tx.Where(&db.Host{Id: hostInfo.Id}).Update(&db.Host{OwnerID: userID}).Error; err != nil {
		return
	}

	// 删除原有的 Owner Host Record
	if err = tx.Delete(&db.HostRecord{HostID: hostInfo.Id, UserID: c.Uid, Type: db.HostRecordTypeOwner}).Error; err != nil {
		return
	}

	hostRecordInfo := db.HostRecord{HostID: hostID, UserID: userID}

	if err = tx.Where(&hostRecordInfo).First(&hostRecordInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果未存在，则创建
			if err = tx.Create(db.HostRecord{
				HostID: hostID,
				UserID: userID,
				Type:   db.HostRecordTypeOwner,
			}).Error; err != nil {
				return
			}
		}
		return
	}

	// 如果已是协作者，则更新资料
	if err = tx.Where(&db.HostRecord{HostID: hostID, UserID: userID}).Update(&db.HostRecord{Type: db.HostRecordTypeOwner}).Error; err != nil {
		return
	}
	return
}
