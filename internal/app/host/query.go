package host

import (
	"errors"
	"net/http"
	"time"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type QueryList struct {
	schema.Query
}

func (s *Service) QueryMyHostByID(c controller.Context, hostID string) (res schema.Response) {
	var (
		err  error
		data = schema.Host{}
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

		helper.Response(&res, data, nil, err)
	}()

	hostRecordInfo := db.HostRecord{
		HostID: hostID,
		UserID: c.Uid,
	}

	if err = db.Db.Where(&hostRecordInfo).Preload("Host", "host.owner_type = ?", string(db.HostOwnerTypeUser)).First(&hostRecordInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(hostRecordInfo.Host, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostRecordInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostRecordInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func (s *Service) QueryMyHostByIDRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.QueryMyHostByID(controller.NewContextFromGinContext(c), c.Param("host_id")))
}

func (s *Service) QueryMyOperationalServer(c controller.Context, input QueryList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.Host, 0)   // 输出到外部的结果
		list  = make([]db.HostRecord, 0) // 数据库查询出来的原始结果
		total int64
		meta  = &schema.Meta{}
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

		helper.Response(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	filter := db.HostRecord{
		UserID: c.Uid,
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Host").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Host{}
		if err = mapstructure.Decode(v.Host, &d.HostPure); err != nil {
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

func (s *Service) QueryMyOperationalServerRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input QueryList
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.QueryMyOperationalServer(controller.NewContextFromGinContext(c), input)
}

func (s *Service) QueryMyHostByTeam(c controller.Context, teamID string, hostID string) (res schema.Response) {
	var (
		err  error
		data = schema.Host{}
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

		helper.Response(&res, data, nil, err)
	}()

	// 查询用户是否在团队中
	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&db.TeamMember{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	hostInfo := db.Host{
		Id:        hostID,
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	if err = db.Db.Where(&hostInfo).First(&hostInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func (s *Service) QueryMyHostByTeamRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.QueryMyHostByTeam(controller.NewContextFromGinContext(c), c.Param("team_id"), c.Param("host_id")))
}

// 获取一个团队可操作的服务器列表
func (s *Service) QueryHostsByTeam(c controller.Context, teamID string, input QueryList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.Host, 0) // 输出到外部的结果
		list  = make([]db.Host, 0)     // 数据库查询出来的原始结果
		total int64
		meta  = &schema.Meta{}
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

		helper.Response(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	filter := db.Host{
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	teamMemberInfo := db.TeamMember{}

	// 查询用户是否在团队中
	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Team").Find(&list).Error; err != nil {
		err = exception.DataBase.New(err.Error())
		return
	}

	for _, v := range list {
		d := schema.Host{}
		if err = mapstructure.Decode(v, &d.HostPure); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

func (s *Service) QueryHostByTeamRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input QueryList
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.QueryHostsByTeam(controller.NewContextFromGinContext(c), c.Param("team_id"), input)
}
