package main

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type UserInfo struct {
	Name   string `bson:"name"`
	Age    uint16 `bson:"age"`
	Weight uint32 `bson:"weight"`
}

func main() {

	ctx := context.Background()

	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://admin:rootroot@10.18.0.2", Database: "class", Coll: "user"})

	defer func() {
		if err = cli.Close(ctx); err != nil {
			panic(err)
		}
	}()

	cli.EnsureIndexes(ctx, []string{}, []string{"age", "name,weight"})

	var userInfo = UserInfo{
		Name:   "xm",
		Age:    7,
		Weight: 40,
	}

	// insert one document
	result, err := cli.InsertOne(ctx, userInfo)
	fmt.Println("result", result)
	// find one document
	one := UserInfo{}
	err = cli.Find(ctx, bson.M{"name": userInfo.Name}).One(&one)

	err = cli.Remove(ctx, bson.M{"age": 7})

	// multiple insert
	var userInfos = []interface{}{
		UserInfo{Name: "a1", Age: 6, Weight: 20},
		UserInfo{Name: "b2", Age: 6, Weight: 25},
		UserInfo{Name: "c3", Age: 6, Weight: 30},
		UserInfo{Name: "d4", Age: 6, Weight: 35},
		UserInfo{Name: "a1", Age: 7, Weight: 40},
		UserInfo{Name: "a1", Age: 8, Weight: 45},
	}
	result1, err := cli.InsertMany(ctx, userInfos)
	fmt.Println("result1", result1)

	// find all 、sort and limit
	batch := []UserInfo{}
	cli.Find(ctx, bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)

	count, err := cli.Find(ctx, bson.M{"age": 6}).Count()
	fmt.Println("count", count)

	// UpdateOne one
	err = cli.UpdateOne(ctx, bson.M{"name": "d4"}, bson.M{"$set": bson.M{"age": 7}})

	// UpdateAll
	result2, err := cli.UpdateAll(ctx, bson.M{"age": 6}, bson.M{"$set": bson.M{"age": 10}})
	fmt.Println("result2", result2)

	err = cli.Find(ctx, bson.M{"age": 10}).Select(bson.M{"age": 1}).One(&one)

	matchStage := bson.D{{"$match", []bson.E{{"weight", bson.D{{"$gt", 30}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$name"}, {"total", bson.D{{"$sum", "$age"}}}}}}
	var showsWithInfo []bson.M
	err = cli.Aggregate(context.Background(), qmgo.Pipeline{matchStage, groupStage}).All(&showsWithInfo)
	fmt.Println("err", err)

}
