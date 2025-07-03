package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Faltan variables de entorno para la base de datos")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", user, password, host, port, dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Error al establecer conexión con la base de datos en SUPABASE: %v", err)
	}

	DB = conn
	fmt.Println("Conexión exitosa con la base de datos en SUPABASE")
}
