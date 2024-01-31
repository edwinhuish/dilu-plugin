# mongodb

## 简单使用案例
**其一**
```go
package main

import (
	"fmt"
	"mongodb"
	"mongodb/bson"
	"strconv"
)

type User struct {
	ID   bson.ObjectID `bson:"_id" json:"_id"`
	Name string        `bson:"name" json:"name"`
}

func main() {
	uri := "mongodb://127.0.0.1:27017"

	mongoSession := mongodb.New(uri)
	mongoSession.SetPoolLimit(10)

	if cleanFunc, err := mongoSession.Connect(); err != nil {
		cleanFunc()
		fmt.Println("客户端连接出错：", err)
		return
	}

	//// 增
	if err := mongoSession.DB("test").Collection("users").Insert(bson.M{"name": "name"}); err != nil {
		fmt.Println("insert error：", err)
	}

	// 增多个
	var docs []interface{}
	for index := 0; index < 10; index++ {
		docs = append(docs, bson.M{"name": strconv.Itoa(index)})
	}

	if err := mongoSession.DB("test").Collection("users").InsertAll(docs); err != nil {
		fmt.Println("insertAll error：", err)
	}

	// 删
	if err := mongoSession.DB("test").Collection("users").Remove(bson.M{"name": "0"}); err != nil {
		fmt.Println("remove error：", err)
	}

	// 删多个
	if err := mongoSession.DB("test").Collection("users").RemoveAll(bson.M{"name": "1"}); err != nil {
		fmt.Println("removeAll error：", err)
	}

	// 改
	if err := mongoSession.DB("test").Collection("users").Update(bson.M{"name": "8"}, bson.M{"$set": bson.M{"name": "name"}}); err != nil {
		fmt.Println("update error：", err)
	}

	// 改多个
	info, err := mongoSession.DB("test").Collection("users").UpdateAll(bson.M{"name": "9"}, bson.M{"$set": bson.M{"name": "name_9"}})
	if err != nil {
		fmt.Println("updateAll error：", err)
	}
	fmt.Println("改多个的返回", info)

	// 查多个
	var resultAll []User
	if err := mongoSession.DB("test").Collection("users").Find(bson.M{}).All(&resultAll); err != nil {
		fmt.Println("find1 error：", err)
	}

	for _, r := range resultAll {
		fmt.Println("查出来了", r.Name)
	}

	var resultAlls []User
	// skip = page * limit
	if err := mongoSession.DB("test").Collection("users").Find(bson.M{"name": "name"}).Skip(2).Limit(2).All(&resultAlls); err != nil {
		fmt.Println("find2 error：", err)
	}

	for _, r := range resultAlls {
		fmt.Println("我又查出来了", r.ID, r.Name)
	}

	// 查
	var resultOne User
	if err := mongoSession.DB("test").Collection("users").Find(bson.M{"name": "name_9"}).One(&resultOne); err != nil {
		fmt.Println("one error：", err)
	}
	fmt.Println("我也查出来了", resultOne, resultOne.ID, resultOne.Name)

	// 统计数量
	count, err := mongoSession.DB("test").Collection("users").Count(bson.M{"name": "name"})
	if err != nil {
		fmt.Println("count error：", err)
	}
	fmt.Println("统计出来了", count)

}


```
**其二**
```go
// cmd/main.go
package main

import (
    "fmt"
    "mongodb/options"
    "mongodb/start"
)

func main() {
    err := start.InitMongo()
    if err != nil {
        fmt.Println(err)
        return
    }
    options.SerSysApi = options.NewSysApiService(start.MongoSession)
    insert()
}

func insert() {
    options.SerSysApi.Insert()
}
```
```go
// start/mongo.go
package start

import (
    "fmt"
    "mongodb"
)

var MongoSession *mongodb.MongoSession

func InitMongo() error {
    uri := "mongodb://127.0.0.1:27017"
    MongoSession = mongodb.New(uri)
    MongoSession.SetPoolLimit(10)
    
    if cleanFunc, err := MongoSession.Connect(); err != nil {
        cleanFunc()
        fmt.Println("客户端连接出错：", err)
        return fmt.Errorf("MongoDB client connection error: %v", err)
    }
    return nil
}

```
```go
// options/mongo.go
package options

import (
    "fmt"
    "mongodb"
    "mongodb/bson"
)

type SysApiService struct {
    *mongodb.MongoSession
}

var SerSysApi = SysApiService{}

func (s *SysApiService) Insert() {
    // options\mongo.go
    //// 增
    if err := s.MongoSession.DB("test").Collection("users").Insert(bson.M{"name": "name"}); err != nil {
        fmt.Println("insert error：", err)
    }
}

```
```go
package options

import "mongodb"

func NewSysApiService(session *mongodb.MongoSession) SysApiService {
    return SysApiService{
        MongoSession: session,
    }
}

```