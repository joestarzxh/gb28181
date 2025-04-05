// Code generated by godddx, DO AVOID EDIT.
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gowvp/gb28181/internal/conf"
	"github.com/gowvp/gb28181/internal/core/config"
	"github.com/gowvp/gb28181/internal/core/config/store/configdb"
	"github.com/ixugo/goddd/pkg/orm"
	"github.com/ixugo/goddd/pkg/web"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ConfigAPI struct {
	configCore config.Core
	conf       *conf.Bootstrap
}

func NewConfigAPI(db *gorm.DB, conf *conf.Bootstrap) ConfigAPI {
	core := config.NewCore(configdb.NewDB(db).AutoMigrate(orm.EnabledAutoMigrate))
	return ConfigAPI{configCore: core, conf: conf}
}

func registerConfig(g gin.IRouter, api ConfigAPI, handler ...gin.HandlerFunc) {
	{
		group := g.Group("/configs", handler...)
		// group.GET("", web.WarpH(api.findConfig))
		// group.GET("/:id", web.WarpH(api.getConfig))
		// group.PUT("/:id", web.WarpH(api.editConfig))
		// group.POST("", web.WarpH(api.addConfig))
		// group.DELETE("/:id", web.WarpH(api.delConfig))

		group.GET("/info", web.WarpH(api.getConfigInfo))
		group.PUT("/info/sip", web.WarpH(api.editSIP))
	}
}

// >>> config >>>>>>>>>>>>>>>>>>>>

func (a ConfigAPI) findConfig(c *gin.Context, in *config.FindConfigInput) (any, error) {
	items, total, err := a.configCore.FindConfig(c.Request.Context(), in)
	return gin.H{"items": items, "total": total}, err
}

func (a ConfigAPI) getConfig(c *gin.Context, _ *struct{}) (any, error) {
	configID, _ := strconv.Atoi(c.Param("id"))
	return a.configCore.GetConfig(c.Request.Context(), configID)
}

func (a ConfigAPI) editConfig(c *gin.Context, in *config.EditConfigInput) (any, error) {
	configID, _ := strconv.Atoi(c.Param("id"))
	return a.configCore.EditConfig(c.Request.Context(), in, configID)
}

func (a ConfigAPI) addConfig(c *gin.Context, in *config.AddConfigInput) (any, error) {
	return a.configCore.AddConfig(c.Request.Context(), in)
}

func (a ConfigAPI) delConfig(c *gin.Context, _ *struct{}) (any, error) {
	configID, _ := strconv.Atoi(c.Param("id"))
	return a.configCore.DelConfig(c.Request.Context(), configID)
}

type getConfigInfoOutput struct {
	SIP conf.SIP `json:"sip"`
}

func (a ConfigAPI) getConfigInfo(c *gin.Context, _ *struct{}) (*getConfigInfoOutput, error) {
	return &getConfigInfoOutput{
		SIP: a.conf.Sip,
	}, nil
}

func (a ConfigAPI) editSIP(_ *gin.Context, in *conf.SIP) (gin.H, error) {
	sip := a.conf.Sip
	if err := copier.Copy(&sip, in); err != nil {
		return nil, web.ErrServer.Msg(err.Error())
	}
	if err := conf.WriteConfig(a.conf, a.conf.ConfigPath); err != nil {
		return nil, web.ErrServer.Msg(err.Error())
	}
	return gin.H{"msg": "ok"}, nil
}
