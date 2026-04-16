package repository

import (
	"context"
	"portofolio-api/clean/config"
	"portofolio-api/clean/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SocialMediaRepository interface {
	CheckExistingSocialMedia(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, socialMedia *domain.SocialMedia) (primitive.ObjectID, error)
}

type socialMediaRepository struct {
	collection *mongo.Collection
}

func NewSocialMediaRepository() SocialMediaRepository {
	return &socialMediaRepository{
		collection: config.DB.Collection("social_media"),
	}
}

func (r *socialMediaRepository) CheckExistingSocialMedia(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *socialMediaRepository) Create(ctx context.Context, socialMedia *domain.SocialMedia) (primitive.ObjectID, error) {
	socialMedia.ID = primitive.NewObjectID()
	socialMedia.CreatedAt = time.Now()
	socialMedia.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, socialMedia)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
