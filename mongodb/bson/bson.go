package bson

import (
	"encoding/hex"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Zero 零值
type Zero interface {
	IsZero() bool
}

// D 有序的键值对组成的 BSON 文档
type D = primitive.D

// E BSON 文档中的单个键值对
type E = primitive.E

// M map[string]interface{} 的别名
type M = primitive.M

// A []interface{} 的别名
type A = primitive.A

// ObjectID MongoDB 中的唯一标识符
type ObjectID = primitive.ObjectID

// NewObjectID 创建一个新的 MongoDB ObjectID
var NewObjectID = primitive.NewObjectID

// ObjectIDHex 将十六进制字符串转换为 MongoDB 的 ObjectID
var ObjectIDHex = primitive.ObjectIDFromHex

// Marshal 将 Go 对象序列化为 BSON 格式的数据
var Marshal = bson.Marshal

// IsObjectIDHex 检查该字符串是否是合法的 MongoDB ObjectID 的十六进制表示
func IsObjectIDHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
