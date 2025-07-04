package utils

import (
	"context"
	"os"
	"runtime"
	"strings"
	"time"

	"menchaca-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var LogCollection *mongo.Collection

func InitLogger(mongoClient *mongo.Client) {
	LogCollection = mongoClient.Database("nombre_de_tu_db").Collection("logs")
}

func LogAction(c *gin.Context, email, action, logLevel string) {
	start := time.Now()

	c.Next()

	duration := time.Since(start).Seconds() * 1000 // ms

	statusCode := c.Writer.Status()
	dynamicLevel := "info"
	if statusCode >= 400 {
		dynamicLevel = "error"
	}

	entry := models.LogEntry{
		Email:        email,
		Action:       action,
		LogLevel:     ifEmpty(logLevel, dynamicLevel),
		Timestamp:    time.Now(),
		IP:           getIP(c),
		UserAgent:    c.Request.UserAgent(),
		Referer:      c.Request.Referer(),
		Origin:       c.Request.Header.Get("Origin"),
		Method:       c.Request.Method,
		URL:          c.Request.RequestURI,
		Status:       statusCode,
		ResponseTime: duration,
		Protocol:     c.Request.Proto,
		Hostname:     getHostname(),
		Environment:  getEnv("GO_ENV", "development"),
		GoVersion:    runtime.Version(),
		PID:          os.Getpid(),
	}

	go func() {
		_, err := LogCollection.InsertOne(context.Background(), entry)
		if err != nil {
			// Aquí puedes usar logrus para registrar el error si deseas
			// logs.Logger.Error("Error al guardar log: ", err)
		}
	}()
}

// Funciones auxiliares
func ifEmpty(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

func getIP(c *gin.Context) string {
	ip := c.ClientIP()
	// Opción adicional si estás detrás de un proxy
	if forwarded := c.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		ip = strings.TrimSpace(ips[0])
	}
	return ip
}

func getHostname() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
