package main

import (
	"log"

	confing "menchaca-backend/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("--Error cargando el archivo .env--")
	}

	confing.InitDB()
}
