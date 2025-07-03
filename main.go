package main

import (
	"fmt"
	"log"
	"net/http"

	config "menchaca-backend/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env:", err)
	}

	config.InitDB()

	err = config.InitMongoDB()
	if err != nil {
		log.Fatalf("Error inicializando MongoDB: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "¡Servidor funcionando con Postgres y MongoDB!")
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
