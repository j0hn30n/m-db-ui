package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	client *mongo.Client
}

func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client, nil
}

func NewService(client *mongo.Client) *Service {
	return &Service{client: client}
}

// 获取所有数据库
func (s *Service) GetDatabases() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.client.ListDatabaseNames(ctx, bson.M{})
}

// 获取数据库信息
func (s *Service) GetDatabase(dbName string) (*DatabaseInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)

	// 获取所有集合
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// 获取统计信息
	stats := &DatabaseStats{}
	if err := db.RunCommand(ctx, bson.M{"dbStats": 1}).Decode(stats); err != nil {
		// 如果无法获取统计信息，使用默认值
		stats = &DatabaseStats{Collections: len(collections)}
	}

	return &DatabaseInfo{
		Name:       dbName,
		Collections: collections,
		Stats:      stats,
	}, nil
}

// 创建数据库
func (s *Service) CreateDatabase(dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建一个简单的集合来创建数据库
	db := s.client.Database(dbName)
	collectionName := "init_collection"

	// 创建一个空集合来初始化数据库
	return db.CreateCollection(ctx, collectionName)
}

// 删除数据库
func (s *Service) DeleteDatabase(dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.client.Database(dbName).Drop(ctx)
}

// 获取数据库的所有集合
func (s *Service) GetCollections(dbName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	return db.ListCollectionNames(ctx, bson.M{})
}

// 创建集合
func (s *Service) CreateCollection(dbName, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	return db.CreateCollection(ctx, collectionName)
}

// 删除集合
func (s *Service) DeleteCollection(dbName, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	return db.Collection(collectionName).Drop(ctx)
}

// 获取集合中的文档
func (s *Service) GetDocuments(dbName, collectionName string, page, limit int64) (*DocumentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	collection := db.Collection(collectionName)

	// 计算跳过的文档数
	skip := (page - 1) * limit

	// 获取总数
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// 获取文档
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"_id": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []map[string]interface{}
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return &DocumentsResponse{
		Documents: documents,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// 创建文档
func (s *Service) CreateDocument(dbName, collectionName string, document map[string]interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	collection := db.Collection(collectionName)

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// 更新文档
func (s *Service) UpdateDocument(dbName, collectionName, id string, document map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	collection := db.Collection(collectionName)

	objectID, err := toObjectID(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": document})
	return err
}

// 删除文档
func (s *Service) DeleteDocument(dbName, collectionName, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	collection := db.Collection(collectionName)

	objectID, err := toObjectID(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// 查询文档
func (s *Service) QueryDocuments(dbName, collectionName string, query map[string]interface{}, page, limit int64) (*DocumentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	collection := db.Collection(collectionName)

	// 计算跳过的文档数
	skip := (page - 1) * limit

	// 获取总数
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	// 获取文档
	cursor, err := collection.Find(ctx, query, options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"_id": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []map[string]interface{}
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return &DocumentsResponse{
		Documents: documents,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// 获取统计信息
func (s *Service) GetStats() (*ServerStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stats ServerStats
	if err := s.client.Database("admin").RunCommand(ctx, bson.M{"serverStatus": 1}).Decode(&stats); err != nil {
		return nil, err
	}

	// 获取数据库列表
	databases, err := s.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	stats.DatabaseCount = len(databases)
	return &stats, nil
}

// GetDocument 获取单个文档
func (s *Service) GetDocument(dbName, collectionName, id string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := s.client.Database(dbName)
	collection := db.Collection(collectionName)

	objectID, err := toObjectID(id)
	if err != nil {
		return nil, err
	}

	var document map[string]interface{}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&document)
	if err != nil {
		return nil, err
	}

	return formatDocument(document), nil
}