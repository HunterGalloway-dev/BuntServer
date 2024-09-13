package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id               primitive.ObjectID `json:"_id" bson:"_id"`
	Title            string             `json:"title" validate:"required"`
	ShortDescription string             `json:"shortDescription" validate:"required"`
	GithubURL        string             `json:"githubURL" validate:"required"`
	YoutubeURL       string             `json:"youtubeURL" validate:"required"`
	ImageURL         string             `json:"imageURL" validate:"required"`
	Description      string             `json:"description" validate:"required"`
	Tags             []string           `json:"tags" validate:"required"`
}

type CodeRunRequest struct {
	Language     string        `json:"language"`
	Stdin        string        `json:"stdin"`
	MainFilePath string        `json:"mainFilePath"`
	Files        []CodeRunFile `json:files`
}

type CodeRunFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type CodeRunOutput struct {
	Output  string
	Err     string
	RunTime time.Duration
}
