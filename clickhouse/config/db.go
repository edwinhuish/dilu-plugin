package config

import "strings"

// DBCfg 数据库连接配置
type DBCfg struct {
	DSN              string `mapstructure:"dsn" json:"dsn" yaml:"dsn"` // 连接参数
	Addr             string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Database         string `mapstructure:"database" json:"database" yaml:"database"`
	Username         string `mapstructure:"username" json:"username" yaml:"username"`
	Password         string `mapstructure:"password" json:"password" yaml:"password"`
	Debug            bool   `mapstructure:"de-bug" json:"de_bug" yaml:"de-bug"` // debug模式
	DialTimeout      int    `mapstructure:"dial-timeout" json:"dial_timeout" yaml:"dial-timeout"`
	MaxOpenConn      int    `mapstructure:"max-open-conn" json:"max_open_conn" yaml:"max-open-conn"`
	MaxIdleConn      int    `mapstructure:"max-idle-conn" json:"max_idle_conn" yaml:"max-idle-conn"`
	ConnMaxLifetime  int    `mapstructure:"conn-max-lifetime" json:"conn_max_lifetime" yaml:"conn-max-lifetime"`
	MaxExecutionTime int    `mapstructure:"max-execution-time" json:"max_execution_time" yaml:"max-execution-time"`
}

// GetAddr 服务地址:端口
func (c *DBCfg) GetAddr() []string {
	return strings.Split(c.Addr, ";")
}

// GetMaxOpenConn 最大链接数量
func (c *DBCfg) GetMaxOpenConn() int {
	if c.MaxOpenConn <= 0 {
		return 5
	}
	return c.MaxIdleConn
}

// GetMaxIdleConn 最大空闲链接数据
func (c *DBCfg) GetMaxIdleConn() int {
	if c.MaxIdleConn <= 0 {
		return 5
	}
	return c.MaxOpenConn
}

// GetConnMaxLifetime 连接的最大生命周期
func (c *DBCfg) GetConnMaxLifetime() int {
	if c.ConnMaxLifetime <= 0 {
		return 15
	}
	return c.ConnMaxLifetime
}

// GetDialTimeout 连接超时时间
func (c *DBCfg) GetDialTimeout() int {
	if c.DialTimeout <= 0 {
		return 30
	}
	return c.DialTimeout
}

// GetMaxExecutionTime 会话设置-最大执行时间
func (c *DBCfg) GetMaxExecutionTime() int {
	if c.MaxExecutionTime <= 0 {
		return 30
	}
	return c.MaxExecutionTime
}
