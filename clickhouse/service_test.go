package clickhouse

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/baowk/dilu-plugin/clickhouse/config"
	"github.com/jmoiron/sqlx"
	"log"
	"testing"
	"time"
)

type OrderLogOperateGrimmTest struct {
	LogTime     time.Time         `db:"logTime" ch:"logTime"`
	TraceId     string            `db:"traceId" ch:"traceId"     `     // 日志跟踪ID
	Channel     string            `db:"channel" ch:"channel"     `     // 日志通道
	Application string            `db:"application" ch:"application" ` // 应用
	Project     string            `db:"project" ch:"project"     `     // 项目
	Source      string            `db:"source" ch:"source"      `      // 来源
	BusinessKey string            `db:"businessKey" ch:"businessKey" ` // 业务key
	Msg         string            `db:"msg" ch:"msg"         `         // 消息内容
	Level       string            `db:"level" ch:"level"       `       // 日志等级
	LevelNum    int32             `db:"levelNum" ch:"levelNum"    `    // 日志等级数字
	Context     map[string]string `db:"context" ch:"context"     `     // 上下文
	SpanId      string            `db:"spanId" ch:"spanId"      `      // 链路ID
	PSpanId     string            `db:"pSpanId" ch:"pSpanId"     `     // 父级链路ID
}

func TestClickhouse(t *testing.T) {
	dsn := "tcp://127.0.0.1:9000/log?username=default&password=test"
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "clickhouse") // 将普通数据库连接转换为sqlx可用的形式

	// 设置上下文
	ctx := context.Background()

	// 原 database/sql 数据库驱动
	selectSQL := "SELECT * FROM order_log_operate_grimm_test"
	rows, err := db.QueryContext(ctx, selectSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sql Open : ", rows)

	var orderLogOperateGrimmTest []OrderLogOperateGrimmTest
	err = sqlxDB.Select(&orderLogOperateGrimmTest, "SELECT * FROM order_log_operate_grimm_test")
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range orderLogOperateGrimmTest {
		fmt.Println(u)
	}

	fmt.Println("sqlx Open : ", len(orderLogOperateGrimmTest))
}

func TestClickhouse2(t *testing.T) {

	testConfig := config.DBCfg{
		DSN:   "clickhouse://default:test@127.0.0.1:9000/log",
		Debug: false,
	}
	Cfg = testConfig
	InitDSN()
	ctx := context.Background()
	list, err := NewService("sys").QueryResMap(ctx, "SELECT * FROM business_log_test")
	if err != nil {
		fmt.Println("SQL ERROR : ", err)
		return
	}

	for _, val := range list {
		if project, ok := val["project"].(*string); ok {
			fmt.Println("输出项目project : ", *project)
		} else {
			fmt.Println("断言错误 : ", ok)
		}
		if operateResult, ok := val["operate_result"].(*int32); ok {
			fmt.Println("输出operate_result : ", *operateResult)
		} else {
			fmt.Println("断言失败 : ", ok)
		}
		if logTime, ok := val["log_time"].(*time.Time); ok {
			fmt.Println("输出log_time : ", *logTime)
		} else {
			fmt.Println("断言失败 : ", ok)
		}
		if con, ok := val["context"].(*map[string]string); ok {
			fmt.Println("输出context : ", *con)
		} else {
			fmt.Println("断言失败 : ", ok)
		}
	}
}

func TestClickhouse3(t *testing.T) {
	testConfig := config.DBCfg{
		Addr:     "127.0.0.1:9000",
		Database: "log",
		Username: "default",
		Password: "test",
		Debug:    false,
	}
	Cfg = testConfig
	Init()
	ctx := context.Background()
	list, err := NewService("sys").QueryResMap(ctx, "SELECT * FROM business_log_test")
	if err != nil {
		fmt.Println("SQL ERROR : ", err)
		return
	}

	for _, val := range list {

		if project, ok := val["project"].(*string); ok {
			fmt.Println("输出项目project : ", *project)
		} else {
			fmt.Println("断言错误 : ", ok)
		}

		if operateResult, ok := val["operate_result"].(*int32); ok {
			fmt.Println("输出operate_result : ", *operateResult)
		} else {
			fmt.Println("断言失败 : ", ok)
		}

		if logTime, ok := val["log_time"].(*time.Time); ok {
			fmt.Println("输出log_time : ", *logTime)
		} else {
			fmt.Println("断言失败 : ", ok)
		}

		if cont, ok := val["context"].(*map[string]string); ok {
			contMap := *cont
			fmt.Println("输出context : ", contMap)
			marshal, err := json.Marshal(contMap)
			if err != nil {
				fmt.Println("json err : ", err)
				return
			}
			fmt.Println(string(marshal))

		} else {
			fmt.Println("断言失败 : ", ok)
		}
	}
}
