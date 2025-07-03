package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitMongoDB() error {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("falta la variable de entorno MONGODB_URI")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("error conectando a MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("no se pudo hacer ping a MongoDB: %v", err)
	}

	MongoClient = client
	MongoDB = client.Database("menchaca")

	fmt.Println("Conexión exitosa con MongoDB Atlas")
	return nil
}


func CloseMongo() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Error al cerrar conexión con MongoDB: %v", err)
		}
		fmt.Println("Conexión con MongoDB cerrada")
	}
}
