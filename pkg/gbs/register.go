package gbs

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gowvp/gb28181/internal/conf"
	"github.com/gowvp/gb28181/internal/core/gb28181"
	"github.com/gowvp/gb28181/pkg/gbs/sip"
	"github.com/ixugo/goweb/pkg/orm"
)

const ignorePassword = "#"

type GB28181API struct {
	cfg   *conf.SIP
	store gb28181.GB28181

	catalog *sip.Collector[Channels]
}

func NewGB28181API(cfg *conf.Bootstrap, store gb28181.GB28181) *GB28181API {
	g := GB28181API{
		cfg:   &cfg.Sip,
		store: store,
		catalog: sip.NewCollector[Channels](func(c1, c2 *Channels) bool {
			return c1.ChannelID == c2.ChannelID
		}),
	}
	go g.catalog.Start(func(s string, c []*Channels) {
		out := make([]*gb28181.Channel, len(c))
		for i, ch := range c {
			out[i] = &gb28181.Channel{
				DeviceID:  s,
				ChannelID: ch.ChannelID,
				Name:      ch.Name,
				IsOnline:  ch.Status == "OK",
				Ext: gb28181.DeviceExt{
					Manufacturer: ch.Manufacturer,
					Model:        ch.Model,
				},
			}
		}
		g.store.SaveChannels(out)
	})
	return &g
}

func (g GB28181API) handlerRegister(ctx *sip.Context) {
	if len(ctx.DeviceID) < 18 {
		ctx.String(http.StatusBadRequest, "device id too short")
		return
	}

	dev, err := g.store.GetDeviceByDeviceID(ctx.DeviceID)
	if err != nil {
		ctx.Log.Error("GetDeviceByDeviceID", "err", err)
		ctx.String(http.StatusInternalServerError, "server db error")
		return
	}

	password := dev.Password
	if password == "" {
		password = g.cfg.Password
	}
	// 免鉴权
	if dev.Password == ignorePassword {
		password = ""
	}
	if password != "" {
		hdrs := ctx.Request.GetHeaders("Authorization")
		if len(hdrs) == 0 {
			resp := sip.NewResponseFromRequest("", ctx.Request, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			resp.AppendHeader(&sip.GenericHeader{HeaderName: "WWW-Authenticate", Contents: fmt.Sprintf("Digest nonce=\"%s\", algorithm=MD5, realm=\"%s\",qop=\"auth\"", sip.RandString(32), g.cfg.Domain)})
			_ = ctx.Tx.Respond(resp)
			return
		}
		authenticateHeader := hdrs[0].(*sip.GenericHeader)
		auth := sip.AuthFromValue(authenticateHeader.Contents)
		auth.SetPassword(password)
		auth.SetUsername(dev.DeviceID)
		auth.SetMethod(ctx.Request.Method())
		auth.SetURI(auth.Get("uri"))
		if auth.CalcResponse() != auth.Get("response") {
			ctx.Log.Info("设备注册鉴权失败")
			ctx.String(http.StatusUnauthorized, "wrong password")
			return
		}
	}

	respFn := func() {
		resp := sip.NewResponseFromRequest("", ctx.Request, http.StatusOK, "OK", nil)
		resp.AppendHeader(&sip.GenericHeader{
			HeaderName: "Date",
			Contents:   time.Now().Format("2006-01-02T15:04:05.000"),
		})
		_ = ctx.Tx.Respond(resp)
	}

	expire := ctx.GetHeader("Expires")
	if expire == "0" {
		ctx.Log.Info("设备注销")
		g.logout(ctx.DeviceID, func(b *gb28181.Device) {
			b.IsOnline = false
			b.Address = ctx.Source.String()
		})
		respFn()
		return
	}
	g.login(ctx.DeviceID, func(b *gb28181.Device) {
		b.IsOnline = true
		b.Address = ctx.Source.String()
		b.Trasnport = strings.ToUpper(ctx.Source.Network())
		b.RegisteredAt = orm.Now()
		b.Expires, _ = strconv.Atoi(expire)
	})

	ctx.Log.Info("设备注册成功")

	respFn()

	g.QueryDeviceInfo(ctx)
	g.QueryCatalog(ctx)
}

func (g GB28181API) login(deviceID string, changeFn func(*gb28181.Device)) {
	g.store.Login(deviceID, changeFn)
}

func (g GB28181API) logout(deviceID string, changeFn func(*gb28181.Device)) {
	g.store.Logout(deviceID, changeFn)
}
