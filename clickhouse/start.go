package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/baowk/dilu-plugin/clickhouse/config"
	"sync"
)

var (
	Cfg  config.DBCfg                          // 数据库配置
	lock sync.RWMutex                          // 锁
	dbs  = make(map[string]clickhouse.Conn, 0) // 数据库DBS
)

// Init 初始化连接 DSN
func Init() {
	dbInitDB()
}

// InitDSN 初始化连接 DSN
func InitDSN() {
	dbInitDSN()
}
