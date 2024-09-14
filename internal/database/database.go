package database

import (
	"BuntServer/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetAllProjects([]string) []models.Project
	GetAllTags() []string
	PostProject(models.Project) primitive.ObjectID
	DeleteProject(string) bool
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
	fmt.Println(tags, len(tags))

	var filter bson.D
	if len(tags) > 0 {
		var test bson.A
		for _, tag := range tags {
			test = append(test, tag)
		}

		filter = bson.D{
			{"tags", bson.D{{"$all", test}}},
		}

	} else {
		filter = bson.D{}
	}

	results, err := projectCollection.Find(ctx, filter)

	if err != nil {
		log.Fatal("Failed to get all objects")
		log.Fatal(err)
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleProject models.Project
		if err = results.Decode(&singleProject); err != nil {
			log.Fatal("Failed decoding project from result")
		}

		projects = append(projects, singleProject)

	}

	return projects
}

func (s *service) GetAllTags() []string {
	var tags []string
	projects := s.GetAllProjects(tags)

	for _, project := range projects {
		for _, tag := range project.Tags {
			if !slices.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

func (s *service) PostProject(project models.Project) primitive.ObjectID {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	projectCollection := s.getCollection("projects")

	if project.Id.IsZero() {
		result, err := projectCollection.InsertOne(ctx, project)

		if err != nil {
			log.Fatal("Failed to create project")
		}

		return result.InsertedID.(primitive.ObjectID)
	} else {
		fmt.Println(project.Id)
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "title", Value: project.Title},
			{Key: "shortDescription", Value: project.ShortDescription},
			{Key: "githubURL", Value: project.GithubURL},
			{Key: "youtubeURL", Value: project.YoutubeURL},
			{Key: "imageURL", Value: project.ImageURL},
			{Key: "description", Value: project.Description},
			{Key: "tags", Value: project.Tags}}}}
		_, err := projectCollection.UpdateByID(ctx, project.Id, update)

		if err != nil {
			log.Fatal(err)
		}

		return project.Id
	}
}

func (s *service) DeleteProject(id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	projectCollection := s.getCollection("projects")

	result, err := projectCollection.DeleteOne(ctx, filter)

	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println(result)

	return true
}
