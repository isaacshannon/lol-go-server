package main

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var users *mongo.Collection
var ipUsers *mongo.Collection

// User stores data about a user
type User struct {
	UserID string
	RequestTimes []int64
	Predictions []string
}
// User stores data about an IP
type UserIP struct {
	IP string
	RequestTimes []time.Time
	RequestUsers []string
}

func SetupUsers(){
	ctx := context.TODO()

	mongoURI, ok := os.LookupEnv("DB_URI")
	if !ok {
		sentry.CaptureException(errors.New("unable to find MONGO_PW in the environment"))
		return
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	log.Println("Connected to MongoDB!")

	users = client.Database("league").Collection("users")
	ipUsers = client.Database("league").Collection("ipusers")
}

func validateUser(userID string) error {
	return nil
}

func validateUserIP(userID, userIP string) error {
	var result UserIP
	filter := bson.M{"ip": userIP}

	update := bson.M{"$push":bson.M{"requestTimes":time.Now(), "requestUsers":userID}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	upsert := true
	err := users.FindOneAndUpdate(ctx, filter, update, &options.FindOneAndUpdateOptions{Upsert:&upsert}).Decode(&result)
	// New ips are expected
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	return nil
}

