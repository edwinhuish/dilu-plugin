package clickhouse

import (
	"fmt"
)

type Query struct {
	Column    string `json:"column"`
	Condition string `json:"condition"`
}

type QueryTag string

const (
	EQ     QueryTag = "eq"
	LIKE   QueryTag = "like"
	ILIKE  QueryTag = "ilike"
	GT     QueryTag = "gt"
	GTE    QueryTag = "gte"
	LT     QueryTag = "lt"
	LTE    QueryTag = "lte"
	IN     QueryTag = "in"
	ISNULL QueryTag = "isnull"
)

func pgSql(query Query) string {
	qTag := QueryTag(query.Condition)
	switch qTag {
	case EQ:
		return fmt.Sprintf("%s = ?", query.Column)
	case ILIKE:
		return fmt.Sprintf("%s ilike ?", query.Column)
	case LIKE:
		return fmt.Sprintf("%s like ?", query.Column)
	case GT:
		return fmt.Sprintf("%s > ?", query.Column)
	case GTE:
		return fmt.Sprintf("%s >= ?", query.Column)
	case LT:
		return fmt.Sprintf("%s < ?", query.Column)
	case LTE:
		return fmt.Sprintf("%s <= ?", query.Column)
	case IN:
		return fmt.Sprintf("%s IN (?)", query.Column)
	case ISNULL:
		return fmt.Sprintf("%s IS NULL", query.Column)
	default:
		return fmt.Sprintf("%s = ?", query.Column)
	}
}
