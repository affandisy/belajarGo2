package inventory

import (
	"belajarGo2/service/inventory"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createInventoryIndex(col *mongo.Collection) error {
	ctx := context.TODO()

	// Check existing indexes
	cur, err := col.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	indexExists := false
	for cur.Next(ctx) {
		var idx bson.M
		if err := cur.Decode(&idx); err != nil {
			return err
		}

		if name, ok := idx["name"].(string); ok && name == "code_1" {
			indexExists = true
			break
		}
	}

	if indexExists {
		return nil // already exists
	}

	// Create unique index on "code"
	model := mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = col.Indexes().CreateOne(ctx, model)
	return err
}

type MongoRepository struct {
	col *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	col := db.Collection("inventories")

	if err := createInventoryIndex(col); err != nil {
		fmt.Println("Error ensuring unique index:", err)
	}

	return &MongoRepository{
		col: col,
	}
}

func (r *MongoRepository) Create(inv inventory.Inventory) (err error) {
	_, err = r.col.InsertOne(context.Background(), inv)
	return err
}

func (r *MongoRepository) ReadAll(page int, limit int) (invs []inventory.Inventory, err error) {
	cursor, err := r.col.Find(context.Background(), bson.M{}, nil)
	if err != nil {
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var inv inventory.Inventory
		if err = cursor.Decode(&inv); err != nil {
			return
		}
		invs = append(invs, inv)
	}
	return
}

func (r *MongoRepository) ReadByCode(code string) (inv inventory.Inventory, err error) {
	err = r.col.FindOne(context.Background(), bson.M{"code": code}).Decode(&inv)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			err = nil
			return
		}
	}
	return
}

func (r *MongoRepository) Update(inv inventory.Inventory) (err error) {
	_, err = r.col.UpdateOne(context.Background(), bson.M{"code": inv.Code}, bson.M{"$set": inv})
	return
}

func (r *MongoRepository) Delete(code string) (err error) {
	_, err = r.col.DeleteOne(context.Background(), bson.M{"code": code})
	return
}
