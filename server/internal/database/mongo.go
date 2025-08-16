// Package database handles all connection between web server and mongodb.
package database

// Reference: https://www.mongodb.com/developer/products/mongodb/build-go-web-application-gin-mongodb-help-ai/
import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define your MongoDB connection string

// Create a global variable to hold our MongoDB connection
var instance *mongo.Client
var validate *validator.Validate

// This function runs before we call our main function and connects to our MongoDB database.
// If it cannot connect, the application stops.
func init() {
	if err := connectToMongoDB(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
	validate = validator.New()
}

// GetInstance returns the instance of mongo client for uses in other packages.
// Technically, the instance isn't singleton per se.
func GetInstance() *mongo.Client {
	return instance
}

// connectToMongoDB connects our web server to... MongoDB.
func connectToMongoDB() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	uri := os.Getenv("MONGODB_URI")
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
