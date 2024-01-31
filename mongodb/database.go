package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database mongo驱动程序数据库
type Database struct {
	database *mongo.Database
}

// CollectionNames 返回数据库中存在的集合名称
func (d *Database) CollectionNames() (names []string, err error) {
	names, err = d.database.ListCollectionNames(context.TODO(), options.ListCollectionsOptions{})
	return
}

// C 设置集合名称
func (d *Database) C(collection string) *Collection {
	return &Collection{collection: d.database.Collection(collection)}
}

// Collection 设置集合名称
func (d *Database) Collection(collection string) *Collection {
	return &Collection{collection: d.database.Collection(collection)}
}
