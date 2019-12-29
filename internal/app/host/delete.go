package host

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
)

// 删除服务器，该操作不可恢复
func (s *Service) DeleteHostByID(c controller.Context, hostID string) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, nil, nil, err)
	}()

	tx = db.Db.Begin()

	hostInfo := db.Host{
		Id:        hostID,
		OwnerID:   c.Uid,
		OwnerType: db.HostOwnerTypeUser,
	}

	hostRecordInfo := db.HostRecord{
		HostID: hostID,
	}

	if err := tx.Where(&hostInfo).Delete(db.Host{}).Error; err != nil {
		return
	}

	if err := tx.Where(&hostRecordInfo).Delete(db.HostRecord{}).Error; err != nil {
		return
	}

	return
}

func (s *Service) DeleteHostByIDRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = s.DeleteHostByID(controller.NewContextFromGinContext(c), c.Param("host_id"))
}

func (s *Service) DeleteHostByTeam(c controller.Context, teamID string, hostID string) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, nil, nil, err)
	}()

	tx = db.Db.Begin()

	teamMemberInfo := db.TeamMember{}

	// 查询用户是否在团队中
	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 只有拥有者和管理员才能删除
	if teamMemberInfo.Role != db.TeamRoleOwner && teamMemberInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	hostInfo := db.Host{
		Id:        hostID,
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	if err := tx.Where(&hostInfo).Delete(db.Host{}).Error; err != nil {
		return
	}

	return
}

func (s *Service) DeleteHostByTeamRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = s.DeleteHostByTeam(controller.NewContextFromGinContext(c), c.Param("team_id"), c.Param("host_id"))
}
