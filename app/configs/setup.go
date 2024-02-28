package configs

import (
	"context"
	"log"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("please create an env file")
	}

	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		log.Println("please include your mongodb database url to your env variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("community-api").Collection(collectionName)
	return collection
}

func UseJWT() *jwtauth.JWTAuth {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	authToken := jwtauth.New("HS256", []byte(jwtSecret), nil)

	return authToken
}
