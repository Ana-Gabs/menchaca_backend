package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"menchaca-backend/config"
)

func main() {
	// 1. Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	// 2. Conectar a Postgres
	config.InitDB()

	// 3. Conectar a MongoDB
	if err := config.InitMongoDB(); err != nil {
		log.Fatalf("Error inicializando MongoDB: %v", err)
	}
	defer config.CloseMongo()

	// 4. Crear instancia de Fiber
	app := fiber.New()

	// 5. Ruta de prueba
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡Servidor funcionando con Fiber, Postgres y MongoDB! 🎉")
	})

	// 6. Graceful shutdown
	go gracefulShutdown(app)

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor iniciado en http://localhost:%s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func gracefulShutdown(app *fiber.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("🧹 Cerrando servidor y conexiones...")
	config.CloseMongo()
	_ = app.Shutdown()
}
