package main

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var users *mongo.Collection

// User stores data about a user
type User struct {
	UserID string
	Password string
	MapX0  float64
	MapX1  float64
	MapY0  float64
	MapY1  float64
}

func SetupUsers(){
	ctx := context.TODO()

	mongoURI, ok := os.LookupEnv("DB_URI")
	if !ok {
		sentry.CaptureException(errors.New("unable to find MONGO_PW in the environment"))
		os.Exit(1)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		sentry.CaptureException(err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		sentry.CaptureException(err)
		os.Exit(1)
	}

	log.Println("Connected to MongoDB!")

	users = client.Database("league").Collection("users")
}