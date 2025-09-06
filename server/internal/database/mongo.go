// Package database handles all connection between web server and mongodb.
package database

// Reference: https://www.mongodb.com/developer/products/mongodb/build-go-web-application-gin-mongodb-help-ai/
import (
	"context"
	"log"
	"os"
	"time"

	"github.com/lnwdevelopers007/job-applier-3000/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define your MongoDB connection string

// Create a global variable to hold our MongoDB connection
var instance *mongo.Client

func init() {
	if err := connectToMongoDB(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
}

// GetConnection returns the instance of mongo client for uses in other packages.
func GetConnection() *mongo.Client {
	return instance
}

// GetDatabase returns the **DATABASE** for use in other packages.
func GetDatabase() *mongo.Database {
	return instance.Database(os.Getenv("DB_NAME"))
}

// connectToMongoDB connects our web server to... MongoDB.
func connectToMongoDB() error {
	if instance != nil {
		return nil
	}

	uri := config.LoadEnv("MONGODB_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	instance = client
	return err
}
