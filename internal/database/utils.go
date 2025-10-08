package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// toObjectID 将字符串转换为ObjectID
func toObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// formatDocument 格式化文档，处理ObjectID
func formatDocument(doc map[string]interface{}) map[string]interface{} {
	formatted := make(map[string]interface{})
	for k, v := range doc {
		if oid, ok := v.(primitive.ObjectID); ok {
			formatted[k] = oid.Hex()
		} else {
			formatted[k] = v
		}
	}
	return formatted
}

// parseQuery 解析查询字符串
func parseQuery(query string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(query), true, &result)
	return result, err
}