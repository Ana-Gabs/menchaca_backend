package utils

import (
	"context"
	"os"
	"runtime"
	"strings"
	"time"

	"menchaca-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// LogCollection es la colección de MongoDB donde se guardan los logs
var LogCollection *mongo.Collection

// InitLogger inicializa la colección de logs de MongoDB
func InitLogger(client *mongo.Client) {
	LogCollection = client.Database(os.Getenv("MONGO_DB")).Collection("logs")
}

// LogAction registra una acción en MongoDB; debe usarse como middleware de Fiber
func LogAction(action string, logLevel string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Ejecuta siguiente handler
		err := c.Next()

		duration := time.Since(start).Seconds() * 1000 // ms

		statusCode := c.Response().StatusCode()
		dynamicLevel := logLevel
		if dynamicLevel == "" {
			dynamicLevel = "info"
		}
		if statusCode >= 400 {
			dynamicLevel = "error"
		}

		entry := models.LogEntry{
			Email:        c.Locals("userEmail").(string),
			Action:       action,
			LogLevel:     dynamicLevel,
			Timestamp:    time.Now(),
			IP:           getIP(c),
			UserAgent:    c.Get("User-Agent"),
			Referer:      c.Get("Referer"),
			Origin:       c.Get("Origin"),
			Method:       c.Method(),
			URL:          string(c.Request().URI().Path()),
			Status:       statusCode,
			ResponseTime: duration,
			Protocol:     c.Request().Protocol(),
			Hostname:     getHostname(),
			Environment:  getEnv("GO_ENV", "development"),
			GoVersion:    runtime.Version(),
			PID:          os.Getpid(),
		}

		// Inserta de forma asíncrona
		go func() {
			_, err := LogCollection.InsertOne(context.Background(), entry)
			if err != nil {
				// opcional: manejar error de logging
			}
		}()

		return err
	}
}

// getIP obtiene la IP del cliente
func getIP(c *fiber.Ctx) string {
	ip := c.IP()
	if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		ip = strings.TrimSpace(ips[0])
	}
	return ip
}

// getHostname retorna el hostname del sistema
func getHostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

// getEnv obtiene variable de entorno con fallback
func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
