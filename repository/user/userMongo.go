package user

import (
	"belajarGo2/service/user"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	col *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		col: db.Collection("users"),
	}
}

func (r *MongoRepository) Create(user user.User) (err error) {
	_, err = r.col.InsertOne(context.Background(), user)
	return
}

func (r *MongoRepository) GetByEmail(email string) (user user.User, err error) {
	err = r.col.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			err = nil
			return
		}
	}
	return
}

func (r *MongoRepository) UpdateEmailVerification(user user.User) (err error) {
	_, err = r.col.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": user})
	return
}
