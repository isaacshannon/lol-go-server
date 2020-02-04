package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		fmt.Println("error: unable to find MONGO_PW in the environment")
		os.Exit(1)
	}
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Connected to MongoDB!")

	users = client.Database("league").Collection("users")
}

func loginUser(email string, password string) (User, error) {
	filter := bson.D{}
	var result User
	err := users.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}