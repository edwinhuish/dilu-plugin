package utils

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
	"time"
)

func TypeAdaptation(rows driver.Rows, columnCount int) []interface{} {
	values := make([]interface{}, columnCount)
	for key, val := range rows.ColumnTypes() {
		switch val.ScanType().String() {
		case "string":
			var value string
			values[key] = &value
		case "int8":
			var value int
			values[key] = &value
		case "int16":
			var value int16
			values[key] = &value
		case "int32":
			var value int32
			values[key] = &value
		case "int64":
			var value int64
			values[key] = &value
		case "uint8":
			var value uint8
			values[key] = &value
		case "uint16":
			var value uint16
			values[key] = &value
		case "uint32":
			var value uint32
			values[key] = &value
		case "uint64":
			var value uint64
			values[key] = &value
		case "float":
			var value float64
			values[key] = &value
		case "float32":
			var value float32
			values[key] = &value
		case "float64":
			var value float64
			values[key] = &value
		case "bool":
			var value bool
			values[key] = &value
		case "*time.Time":
			var value time.Time
			values[key] = &value
		case "map[int]int":
			var value map[int]int
			values[key] = &value
		case "map[int]string":
			var value map[int]string
			values[key] = &value
		case "map[string]string":
			var value map[string]string
			values[key] = &value
		case "map[string]int":
			var value map[string]int
			values[key] = &value
		}
	}
	return values
}

func Append[T any](batch driver.Batch, data []T) driver.Batch {
	for _, v := range data {
		err := batch.AppendStruct(v)
		if err != nil {
			log.Printf("db:%s not init", err)
		}
	}
	return batch
}
