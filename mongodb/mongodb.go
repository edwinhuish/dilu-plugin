package mongodb

import (
	"context"
	"log"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DefaultTimeout = 60

// MongoSession mongo session
type MongoSession struct {
	client         *mongo.Client
	collection     *mongo.Collection
	maxPoolSize    uint64
	connectTimeout uint
	db             string
	uri            string
	m              sync.RWMutex
	filter         interface{}
	limit          *int64
	project        interface{}
	skip           *int64
	sort           interface{}
	withTimeout    uint
}

// New 创建mongo session
func New(uri string) *MongoSession {
	mongoSession := &MongoSession{
		uri: uri,
	}
	return mongoSession
}

// C Collection 集合名称名称
func (s *MongoSession) C(collection string) *Collection {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.db) == 0 {
		s.db = "test"
	}
	d := &Database{database: s.client.Database(s.db)}
	return &Collection{collection: d.database.Collection(collection)}
}

// Collection 集合名称名称
func (s *MongoSession) Collection(collection string) *Collection {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.db) == 0 {
		s.db = "test"
	}
	d := &Database{database: s.client.Database(s.db)}
	return &Collection{collection: d.database.Collection(collection)}
}

// SetPoolLimit 指定服务器连接池的最大大小
func (s *MongoSession) SetPoolLimit(limit uint64) {
	s.m.Lock()
	defer s.m.Unlock()
	s.maxPoolSize = limit
}

// SetConnectTimeout 连接mongodb服务超时时间
// 秒
func (s *MongoSession) SetConnectTimeout(connectTimeout uint) {
	s.m.Lock()
	defer s.m.Unlock()
	s.connectTimeout = connectTimeout
}

// SetWithTimeOut 操作超时时间
// 秒
func (s *MongoSession) SetWithTimeOut(withTimeout uint) {
	s.m.Lock()
	defer s.m.Unlock()
	s.withTimeout = withTimeout
}

// Connect 连接 mongo 客户端
func (s *MongoSession) Connect() (func(), error) {

	var (
		ctx    = context.Background()
		cancel context.CancelFunc
	)
	if timeout := s.withTimeout; timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	} else {
		// 默认20秒
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(DefaultTimeout)*time.Second)
	}
	defer cancel()

	// 设置连接mongo超时时间
	var connectTimeout uint
	if s.connectTimeout <= 0 {
		connectTimeout = 30
	} else {
		connectTimeout = s.connectTimeout
	}

	opts := options.Client().ApplyURI(s.uri).
		SetConnectTimeout(time.Duration(connectTimeout) * time.Second).
		SetMaxPoolSize(s.maxPoolSize)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 断开
	cleanFunc := func() {
		err = client.Disconnect(context.Background())
		if err != nil {
			log.Printf("Mongo disconnect error: %s", err.Error())
		}
	}
	s.client = client
	return cleanFunc, nil
}

// Ping 验证客户端是否可以连接
func (s *MongoSession) Ping() error {
	return s.client.Ping(context.TODO(), readpref.Primary())
}

// Client 返回mongo 客户端
func (s *MongoSession) Client() *mongo.Client {
	return s.client
}

// DB 设置数据库名
func (s *MongoSession) DB(db string) *Database {
	s.m.Lock()
	defer s.m.Unlock()
	return &Database{database: s.client.Database(db)}
}

// Sort 设置排序
func (s *MongoSession) Sort(sort interface{}) *MongoSession {
	s.sort = sort
	return s
}

// Limit 设置限制
func (s *MongoSession) Limit(limit int64) *MongoSession {
	s.limit = &limit
	return s
}

// Skip 设置跳过
func (s *MongoSession) Skip(skip int64) *MongoSession {
	s.skip = &skip
	return s
}

// One 返回一条数据
func (s *MongoSession) One(result interface{}) error {
	var err error
	err = s.collection.FindOne(context.TODO(), s.filter).Decode(result)

	if err != nil {
		return err
	}
	return nil
}

// All 返回多条数据
func (s *MongoSession) All(result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		panic("result argument must be a slice address")
	}
	slicev := resultv.Elem()

	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	if slicev.Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error

	opt := options.Find()

	if s.sort != nil {
		opt.SetSort(s.sort)
	}

	if s.limit != nil {
		opt.SetLimit(*s.limit)
	}

	if s.skip != nil {
		opt.SetSkip(*s.skip)
	}

	cur, err := s.collection.Find(ctx, s.filter, opt)
	defer cur.Close(ctx)
	if err != nil {
		return err
	}
	if err = cur.Err(); err != nil {
		return err
	}
	i := 0
	for cur.Next(ctx) {
		elemp := reflect.New(elemt)
		if err = bson.Unmarshal(cur.Current, elemp.Interface()); err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		i++
	}
	resultv.Elem().Set(slicev.Slice(0, i))
	return nil
}
