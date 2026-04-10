package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portofolio-api/config"
)

var DB *mongo.Database

func Connect() {
	// Use config from global variable
	mongoURI := config.Config.MongoURI
	databaseName := config.Config.MongoDB

	// Set up connection options
	clientOptions := options.Client().ApplyURI(mongoURI).
		SetMaxPoolSize(200).
		SetMinPoolSize(20).
		SetMaxConnIdleTime(1 * time.Minute)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB connection error:", err)
	}

	// Ping the database to test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping error:", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")

	// Set the selected database
	DB = client.Database(databaseName)
	InitIndexes()
}

func InitIndexes() {
	SocialMediaIndexes()
	PortofolioIndexes()
	ProjectIndexes()
	UserIndexes()
}

func SocialMediaIndexes() {
	ctx := context.Background()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := DB.Collection(config.CollectionSocialMedia).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal(err)
	}
}

func PortofolioIndexes() {
	ctx := context.Background()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_deleted", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}, {Key: "is_deleted", Value: 1}},
		},
	}

	_, err := DB.Collection(config.CollectionPortofolio).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal(err)
	}
}

func ProjectIndexes() {
	ctx := context.Background()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "portofolio_id", Value: 1}, {Key: "is_deleted", Value: 1}},
		},
	}

	_, err := DB.Collection(config.CollectionProject).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal(err)
	}
}

func UserIndexes() {
	ctx := context.Background()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := DB.Collection(config.CollectionUsers).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal(err)
	}
}
