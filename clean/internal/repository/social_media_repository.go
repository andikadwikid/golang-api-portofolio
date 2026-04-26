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

type SocialMediaRepository interface {
	CheckExistingSocialMedia(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, socialMedia *domain.SocialMedia) (primitive.ObjectID, error)
	GetAllSocialMedia(ctx context.Context) ([]domain.SocialMedia, error)
	GetSocialMediaById(ctx context.Context, id primitive.ObjectID) (domain.SocialMedia, error)
	UpdateSocialMedia(ctx context.Context, id primitive.ObjectID, socialMedia *domain.SocialMedia) error
	DeleteSocialMedia(ctx context.Context, id primitive.ObjectID) error
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

func (r *socialMediaRepository) GetAllSocialMedia(ctx context.Context) ([]domain.SocialMedia, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var socialMedias []domain.SocialMedia
	if err := cursor.All(ctx, &socialMedias); err != nil {
		return nil, err
	}

	return socialMedias, nil
}

func (r *socialMediaRepository) GetSocialMediaById(ctx context.Context, id primitive.ObjectID) (domain.SocialMedia, error) {
	var socialMedia domain.SocialMedia
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&socialMedia)
	if err != nil {
		return domain.SocialMedia{}, err
	}
	return socialMedia, nil
}

func (r *socialMediaRepository) UpdateSocialMedia(ctx context.Context, id primitive.ObjectID, socialMedia *domain.SocialMedia) error {
	filter := bson.M{"_id": id, "is_deleted": false}
	update := bson.M{"$set": bson.M{"name": socialMedia.Name, "icon": socialMedia.Icon}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *socialMediaRepository) DeleteSocialMedia(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id, "is_deleted": false}
	update := bson.M{"$set": bson.M{"is_deleted": true}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
