package mongodb

import (
	"context"

	"mongodb/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection mongo驱动程序集合
type Collection struct {
	collection *mongo.Collection
}

// Find 按给定筛选器查找文档
func (c *Collection) Find(filter interface{}) *MongoSession {
	return &MongoSession{filter: filter, collection: c.collection}
}

// Insert 将单个文档插入到集合中
func (c *Collection) Insert(document interface{}) error {
	var err error
	if _, err = c.collection.InsertOne(context.TODO(), document); err != nil {
		return err
	}
	return nil
}

// InsertWithResult 将单个文档插入到集合中，并返回insert one result
func (c *Collection) InsertWithResult(document interface{}) (result *mongo.InsertOneResult, err error) {
	result, err = c.collection.InsertOne(context.TODO(), document)
	return
}

// InsertAll 插入提供的文档
func (c *Collection) InsertAll(documents []interface{}) error {
	var err error
	if _, err = c.collection.InsertMany(context.TODO(), documents); err != nil {
		return err
	}
	return nil
}

// InsertAllWithResult 插入提供的文档并返回insert many result
func (c *Collection) InsertAllWithResult(documents []interface{}) (result *mongo.InsertManyResult, err error) {
	result, err = c.collection.InsertMany(context.TODO(), documents)
	return
}

// Update 更新集合中的单个文档
func (c *Collection) Update(selector interface{}, update interface{}, upsert ...bool) error {
	if selector == nil {
		selector = bson.D{}
	}

	var err error

	opt := options.Update()
	for _, arg := range upsert {
		if arg {
			opt.SetUpsert(arg)
		}
	}

	if _, err = c.collection.UpdateOne(context.TODO(), selector, update, opt); err != nil {
		return err
	}
	return nil
}

// UpdateWithResult 更新集合中的单个文档并返回更新结果
func (c *Collection) UpdateWithResult(selector interface{}, update interface{}, upsert ...bool) (result *mongo.UpdateResult, err error) {
	if selector == nil {
		selector = bson.D{}
	}

	opt := options.Update()
	for _, arg := range upsert {
		if arg {
			opt.SetUpsert(arg)
		}
	}

	result, err = c.collection.UpdateOne(context.TODO(), selector, update, opt)
	return
}

// UpdateID 按id更新集合中的单个文档
func (c *Collection) UpdateID(id interface{}, update interface{}) error {
	return c.Update(bson.M{"_id": id}, update)
}

// UpdateAll 更新集合中的多个文档
func (c *Collection) UpdateAll(selector interface{}, update interface{}, upsert ...bool) (*mongo.UpdateResult, error) {
	if selector == nil {
		selector = bson.D{}
	}

	var err error

	opt := options.Update()
	for _, arg := range upsert {
		if arg {
			opt.SetUpsert(arg)
		}
	}

	var updateResult *mongo.UpdateResult
	if updateResult, err = c.collection.UpdateMany(context.TODO(), selector, update, opt); err != nil {
		return updateResult, err
	}
	return updateResult, nil
}

// Remove 从集合中删除单个文档
func (c *Collection) Remove(selector interface{}) error {
	if selector == nil {
		selector = bson.D{}
	}
	var err error
	if _, err = c.collection.DeleteOne(context.TODO(), selector); err != nil {
		return err
	}
	return nil
}

// RemoveID 按id从集合中删除单个文档
func (c *Collection) RemoveID(id interface{}) error {
	return c.Remove(id)
}

// RemoveAll 从集合中删除多个文档
func (c *Collection) RemoveAll(selector interface{}) error {
	if selector == nil {
		selector = bson.D{}
	}
	var err error

	if _, err = c.collection.DeleteMany(context.TODO(), selector); err != nil {
		return err
	}
	return nil
}

// Count 获取与筛选器匹配的文档数
func (c *Collection) Count(selector interface{}) (int64, error) {
	if selector == nil {
		selector = bson.D{}
	}
	var err error
	var count int64
	count, err = c.collection.CountDocuments(context.TODO(), selector)
	return count, err
}
