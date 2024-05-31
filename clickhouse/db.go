package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/baowk/dilu-plugin/clickhouse/consts"
	"log"
	"time"
)

// dbInit 初始化DSN数据库驱动
func dbInitDSN() {
	opt, err := clickhouse.ParseDSN(Cfg.DSN)
	if err != nil {
		panic(err)
	}

	opt.Debug = Cfg.Debug
	db, err := clickhouse.Open(opt)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()
	SetDb(consts.DB_DEF, db)
}

// dbInitDB 初始化数据库驱动
func dbInitDB() {
	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: Cfg.GetAddr(), // 地址和端口号
		Auth: clickhouse.Auth{ // 身份验证
			Database: Cfg.Database, // 数据库
			Username: Cfg.Username, // 账号
			Password: Cfg.Password, // 密码
		},
		Debug: Cfg.Debug, // 调试模式
		Settings: clickhouse.Settings{ // 会话设置
			"max_execution_time": Cfg.GetMaxExecutionTime(), // 最大执行时间
		},
		Compression: &clickhouse.Compression{ // 压缩设置
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:     time.Second * time.Duration(Cfg.GetDialTimeout()),     // 连接超时时间
		MaxOpenConns:    Cfg.GetMaxOpenConn(),                                  // 最大链接数量
		MaxIdleConns:    Cfg.GetMaxIdleConn(),                                  // 最大空闲链接数据
		ConnMaxLifetime: time.Minute * time.Duration(Cfg.GetConnMaxLifetime()), // 连接的最大生命周期
		//ConnOpenStrategy:     clickhouse.ConnOpenInOrder,      // 连接打开策略
		//BlockBufferSize:      10,    // 缓冲区大小
		//MaxCompressionBuffer: 10240, // 压缩缓冲区大小
		//ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
		//	Products: []struct {
		//		Name    string
		//		Version string
		//	}{
		//		{Name: "my-app", Version: "0.1"},
		//	},
		//},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()
	SetDb(consts.DB_DEF, db)
}

// SetDb 设置指定DBS名称对应的db
// @key DBS名称
// @db 数据库DB
func SetDb(key string, db clickhouse.Conn) {
	lock.Lock()
	defer lock.Unlock()
	dbs[key] = db
}

// Db 指定的DBS（sub）下的db
// @name DBS名称
func Db(name string) clickhouse.Conn {
	lock.RLock()
	defer lock.RUnlock()
	if db, ok := dbs[name]; !ok || db == nil {
		log.Printf("db:%s not init", name)
		panic("db not init")
	} else {
		return db
	}
}
