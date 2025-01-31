// Code generated by gowebx, DO AVOID EDIT.
package gb28181

import (
	"fmt"

	"github.com/ixugo/goweb/pkg/orm"
)

// Device domain model
type Device struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	DeviceID     string    `gorm:"column:device_id;notNull;uniqueIndex;default:'';comment:20 位国标编号" json:"device_id"`                              // 20 位国标编号
	Name         string    `gorm:"column:name;notNull;default:'';comment:设备名称" json:"name"`                                                        // 设备名称
	Trasnport    string    `gorm:"column:trasnport;notNull;default:'';comment:传输协议(TCP/UDP)" json:"trasnport"`                                     // 传输协议(TCP/UDP)
	StreamMode   string    `gorm:"column:stream_mode;notNull;default:'TCP_PASSIVE';comment:数据传输模式(UDP/TCP_PASSIVE,TCP_ACTIVE)" json:"stream_mode"` // 数据传输模式(UDP/TCP_PASSIVE,TCP_ACTIVE)
	IP           string    `gorm:"column:ip;notNull;default:''" json:"ip"`
	Port         int       `gorm:"column:port;notNull;default:0" json:"port"`
	IsOnline     bool      `gorm:"column:is_online;notNull;default:FALSE" json:"is_online"`
	RegisteredAt orm.Time  `gorm:"column:registered_at;notNull;default:CURRENT_TIMESTAMP;comment:注册时间" json:"registered_at"` // 注册时间
	KeepaliveAt  orm.Time  `gorm:"column:keepalive_at;notNull;default:CURRENT_TIMESTAMP;comment:心跳时间" json:"keepalive_at"`   // 心跳时间
	Keepalives   int       `gorm:"column:keepalives;notNull;default:0;comment:心跳间隔" json:"keepalives"`                       // 心跳间隔
	Expires      int       `gorm:"column:expires;notNull;default:0;comment:注册有效期" json:"expires"`                            // 注册有效期
	Channels     int       `gorm:"column:channels;notNull;default:0;comment:通道数量" json:"channels"`                           // 通道数量
	CreatedAt    orm.Time  `gorm:"column:created_at;notNull;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`       // 创建时间
	UpdatedAt    orm.Time  `gorm:"column:updated_at;notNull;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`       // 更新时间
	Password     string    `gorm:"column:password;notNull;default:'';comment:注册密码" json:"password"`
	Address      string    `gorm:"column:address;notNull;default:'';comment:设备网络地址" json:"address"`
	Ext          DeviceExt `gorm:"column:ext;notNull;default:'{}';type:jsonb;comment:设备属性" json:"ext"` // 设备属性
}

// TableName database table name
func (*Device) TableName() string {
	return "devices"
}

func (d Device) Check() error {
	if len(d.DeviceID) < 18 {
		return fmt.Errorf("国标 ID 长度应大于等于 18 位")
	}
	return nil
}

func (d *Device) init(id, deviceID string) {
	d.ID = id
	d.DeviceID = deviceID
}
