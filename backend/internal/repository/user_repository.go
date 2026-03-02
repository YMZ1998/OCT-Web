package repository

import (
	"context"
	"oct-backend/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Coll: db.Collection("users"),
	}
}

func (r *UserRepository) Create(user *model.User) error {
	user.CreatedAt = time.Now().Unix()
	_, err := r.Coll.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.Coll.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	return &user, err
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.Coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	return &user, err
}

func (r *UserRepository) Update(id string, update bson.M) error {
	_, err := r.Coll.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.Coll.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *UserRepository) UpdateLastLoginAt(id primitive.ObjectID, ts int64) error {
	_, err := r.Coll.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"last_login_at": ts}},
	)
	return err
}
