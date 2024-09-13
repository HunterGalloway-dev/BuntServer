package database

import (
	"BuntServer/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetAllProjects([]string) []models.Project
	GetAllTags() []string
}

type service struct {
	db *mongo.Client
}

func New() Service {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_ROOT_PASSWORD")

	mongoURL := fmt.Sprintf("mongodb://%s:%s", host, port)

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: dbUsername,
		Password: dbPassword,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return &service{
		db: client,
	}
}

// Health implements Service.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)

	if err != nil {
		log.Fatalf("Databse Offline: %v", err)
	}

	return map[string]string{
		"message": "Database online",
	}
}

func (s *service) getCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("INIT_DB")
	collection := s.db.Database(dbName).Collection(collectionName)

	return collection
}

func (s *service) GetAllProjects(tags []string) []models.Project {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var projects []models.Project
	defer cancel()

	projectCollection := s.getCollection("projects")

	filter := bson.D{{
		Key: "tags", Value: bson.D{{Key: "$all", Value: tags}}}}

	results, err := projectCollection.Find(ctx, filter)

	if err != nil {
		log.Fatal("Failed to get all objects")
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleProject models.Project
		if err = results.Decode(&singleProject); err != nil {
			log.Fatal("Failed decoding project from result")
		}

		if len(tags) > 0 {

		} else {
			projects = append(projects, singleProject)
		}

	}

	return projects
}

func (s *service) GetAllTags() []string {
	var tags []string
	projects := s.GetAllProjects(tags)

	for _, project := range projects {
		tags = append(tags, project.Tags...)
	}

	return tags
}
