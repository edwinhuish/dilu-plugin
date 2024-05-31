package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/baowk/dilu-plugin/clickhouse/utils"
	"log"
)

func NewService(dbname string) *BaseService {
	return &BaseService{
		DbName: dbname,
	}
}

type BaseService struct {
	DbName string
}

// DB DBS操作指定数据库
func (s *BaseService) DB() clickhouse.Conn {

	return Db(s.DbName)
}

// Insert 增加
func (s *BaseService) Insert(ctx context.Context, table string, data any) error {
	var err error
	batch, err := s.DB().PrepareBatch(ctx, fmt.Sprintf("INSERT INTO %s", table))
	err = batch.AppendStruct(data)
	if err != nil {
		return err
	}

	err = batch.Send()
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (s *BaseService) Delete(ctx context.Context, table string, query Query, values []any) error {
	var err error
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE %s", table, pgSql(query))
	err = s.DB().Exec(ctx, deleteStmt, values...)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新
func (s *BaseService) Update(ctx context.Context, table string, column string, setValue any, query Query, values any) error {
	var err error
	updateStmt := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s", table, column, pgSql(query))

	err = s.DB().Exec(ctx, updateStmt, setValue, values)
	if err != nil {
		return err
	}

	return nil
}

// QueryResMap 查询
func (s *BaseService) QueryResMap(ctx context.Context, sql string) ([]map[string]interface{}, error) {
	var err error
	rows, err := s.DB().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	columns := rows.Columns()
	columnCount := len(columns)

	var records []map[string]interface{}
	for rows.Next() {
		rowRow := make(map[string]interface{}, columnCount)
		values := utils.TypeAdaptation(rows, columnCount)
		if err = rows.Scan(values...); err != nil {
			log.Fatal(err)
		} else {

			for i, col := range columns {
				rowRow[col] = values[i]
			}
			records = append(records, rowRow)
		}
	}

	return records, nil
}

// Creates 批量增加
func (s *BaseService) Creates(ctx context.Context, table string, data [][]interface{}) error {
	var err error
	batch, err := s.DB().PrepareBatch(ctx, fmt.Sprintf("INSERT INTO %s", table))
	if err != nil {
		return err
	}

	defer func() {
		_ = batch.Flush()
	}()

	batch = utils.Append(batch, data)
	err = batch.Send()
	if err != nil {
		return err
	}

	return nil
}
