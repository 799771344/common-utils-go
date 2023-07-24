package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Client struct {
		client *mongo.Client
		dbName string
	}
	QueryParams struct {
		Find       func(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
		Insert     func(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		InsertMany func(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
		Delete     func(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		DeleteMany func(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
		Update     func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
		UpdateMany func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	}
)

func NewClient(host, username, password, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: username,
		Password: password,
	}
	clientOptions := options.Client().ApplyURI(host).SetAuth(credential)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("failed to connect mongo")
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.New("failed to connect mongo")
	}
	return &Client{
		client: mongoClient,
		dbName: dbName,
	}, nil
}

func (client *Client) GetCollection(collectionName string) *QueryParams {
	collection := client.client.Database(client.dbName).Collection(collectionName)

	return &QueryParams{
		Find:       collection.Find,
		Insert:     collection.InsertOne,
		InsertMany: collection.InsertMany,
		Delete:     collection.DeleteOne,
		DeleteMany: collection.DeleteMany,
		Update:     collection.UpdateOne,
		UpdateMany: collection.UpdateMany,
	}
}

func (client *Client) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client.client.Disconnect(ctx)
}

func main() {
	// 创建一个客户端
	client, err := NewClient("mongodb://127.0.0.1:27017", "username", "password", "testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 获取mongo集合
	query := client.GetCollection("user")

	// 查询
	opts := options.Find()
	opts.SetSort(bson.M{"name": 1})
	cursor, err := query.Find(context.Background(), bson.M{"age": bson.M{"$gte": 18}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 插入
	_, err = query.Insert(context.Background(), bson.M{"name": "张三", "age": 18})
	if err != nil {
		log.Fatal(err)
	}

	// 删除
	_, err = query.Delete(context.Background(), bson.M{"name": "张三"})
	if err != nil {
		log.Fatal(err)
	}

	// 更新
	_, err = query.Update(context.Background(), bson.M{"name": "李四"}, bson.M{"$set": bson.M{"age": 19}})
	if err != nil {
		log.Fatal(err)
	}

}
