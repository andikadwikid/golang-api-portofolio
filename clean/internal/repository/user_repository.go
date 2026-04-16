package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/clean/config"
	"portofolio-api/clean/internal/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (primitive.ObjectID, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{
		collection: config.DB.Collection("users"),
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err // Kembalikan error mentah, biarkan Service yang interpretasi
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updateData}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return result, err
}
