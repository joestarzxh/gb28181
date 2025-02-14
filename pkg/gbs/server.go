package gbs

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/gowvp/gb28181/internal/conf"
	"github.com/gowvp/gb28181/internal/core/gb28181"
	"github.com/gowvp/gb28181/internal/core/sms"
	"github.com/gowvp/gb28181/pkg/gbs/m"
	"github.com/gowvp/gb28181/pkg/gbs/sip"
)

type Server struct {
	*sip.Server
	gb           *GB28181API
	mediaService sms.Core

	fromAddress *sip.Address

	devices *Client
}

func NewServer(cfg *conf.Bootstrap, store gb28181.GB28181, sc sms.Core) (*Server, func()) {
	api := NewGB28181API(cfg, store, sc.NodeManager)

	uri, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%d", cfg.Sip.ID, cfg.Sip.Host, cfg.Sip.Port))
	from := sip.Address{
		DisplayName: sip.String{Str: "gowvp"},
		URI:         &uri,
		Params:      sip.NewParams(),
	}

	svr = sip.NewServer(&from)
	svr.Register(api.handlerRegister)
	msg := svr.Message()
	msg.Handle("Keepalive", api.sipMessageKeepalive)
	msg.Handle("Catalog", api.sipMessageCatalog)
	msg.Handle("DeviceInfo", api.sipMessageDeviceInfo)

	// msg.Handle("RecordInfo", api.handlerMessage)

	c := Server{
		Server:       svr,
		mediaService: sc,
		fromAddress:  &from,
		devices:      NewClient(),
		gb:           api,
	}
	api.svr = &c

	// devices, err := store.FindDevices(context.TODO())
	// if err != nil {
	// 	panic(err)
	// }
	// for _, device := range devices {
	// c.devices.Store(device.DeviceID, newDevice(device.NetworkAddress(), device.DeviceID))
	// }

	go svr.ListenUDPServer(fmt.Sprintf(":%d", cfg.Sip.Port))
	go svr.ListenTCPServer(fmt.Sprintf(":%d", cfg.Sip.Port))

	return &c, c.Close
}

func Start() {
	// 数据库表初始化 启动时自动同步数据结构到数据库
	// db.DBClient.AutoMigrate(new(Devices))
	// db.DBClient.AutoMigrate(new(Channels))
	// db.DBClient.AutoMigrate(new(Streams))
	// db.DBClient.AutoMigrate(new(m.SysInfo))
	// db.DBClient.AutoMigrate(new(Files))

	LoadSYSInfo()

	// svr = sip.NewServer()
	// go svr.ListenUDPServer(config.UDP)
}

// MODDEBUG MODDEBUG
var MODDEBUG = "DEBUG"

// ActiveDevices 记录当前活跃设备，请求播放时设备必须处于活跃状态
type ActiveDevices struct {
	sync.Map
}

// Get Get
func (a *ActiveDevices) Get(key string) (Devices, bool) {
	if v, ok := a.Load(key); ok {
		return v.(Devices), ok
	}
	return Devices{}, false
}

var _activeDevices ActiveDevices

// 系统运行信息
var (
	_sysinfo *m.SysInfo
	config   *m.Config
)

func LoadSYSInfo() {
	config = m.MConfig
	_activeDevices = ActiveDevices{sync.Map{}}

	StreamList = streamsList{&sync.Map{}, &sync.Map{}, 0}
	ssrcLock = &sync.Mutex{}
	_recordList = &sync.Map{}
	RecordList = apiRecordList{items: map[string]*apiRecordItem{}, l: sync.RWMutex{}}

	// init sysinfo
	// _sysinfo = &m.SysInfo{}
	// if err := db.Get(db.DBClient, _sysinfo); err != nil {
	// 	if db.RecordNotFound(err) {
	// 		//  初始不存在
	// 		_sysinfo = m.DefaultInfo()

	// 		if err = db.Create(db.DBClient, _sysinfo); err != nil {
	// 			// logrus.Fatalf("1 init sysinfo err:%v", err)
	// 		}
	// 	} else {
	// 		// logrus.Fatalf("2 init sysinfo err:%v", err)
	// 	}
	// }
	m.MConfig.GB28181 = _sysinfo

	// uri, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s", _sysinfo.LID, _sysinfo.Region))
	_serverDevices = Devices{
		DeviceID: _sysinfo.LID,
		// Region:   _sysinfo.Region,
		addr: &sip.Address{
			DisplayName: sip.String{Str: "sipserver"},
			// URI:         &uri,
			Params: sip.NewParams(),
		},
	}

	// init media
	url, err := url.Parse(config.Media.RTP)
	if err != nil {
		// logrus.Fatalf("media rtp url error,url:%s,err:%v", config.Media.RTP, err)
	}
	ipaddr, err := net.ResolveIPAddr("ip", url.Hostname())
	if err != nil {
		// logrus.Fatalf("media rtp url error,url:%s,err:%v", config.Media.RTP, err)
	}
	_sysinfo.MediaServerRtpIP = ipaddr.IP
	_sysinfo.MediaServerRtpPort, _ = strconv.Atoi(url.Port())
}

// zlm接收到的ssrc为16进制。发起请求的ssrc为10进制
func ssrc2stream(ssrc string) string {
	if ssrc[0:1] == "0" {
		ssrc = ssrc[1:]
	}
	num, _ := strconv.Atoi(ssrc)
	return fmt.Sprintf("%08X", num)
}

func sipResponse(tx *sip.Transaction) (*sip.Response, error) {
	response := tx.GetResponse()
	if response == nil {
		return nil, sip.NewError(nil, "response timeout", "tx key:", tx.Key())
	}
	if response.StatusCode() != http.StatusOK {
		return response, sip.NewError(nil, "device: ", response.StatusCode(), " ", response.Reason())
	}
	return response, nil
}

// QueryCatalog 查询 catalog
func (s *Server) QueryCatalog(deviceID string) error {
	return s.gb.QueryCatalog(deviceID)
}

func (s *Server) Play(in *PlayInput) error {
	return s.gb.Play(in)
}

func (s *Server) StopPlay(in *StopPlayInput) error {
	return s.gb.StopPlay(in)
}
