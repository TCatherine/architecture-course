// file: main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smart-home-service/db"
	"smart-home-service/message_broker"
	"smart-home-service/services"
	"smart-home-service/handlers"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// --- Инициализация БД ---
	dbURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/smart_home")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()
	log.Println("Connected to database successfully")

	// --- Инициализация брокера сообщений (Publisher) ---
	amqpURL := getEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/")
	publisher, err := message_broker.NewPublisher(amqpURL)
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v", err)
	}
	defer publisher.Close()
	log.Println("Connected to RabbitMQ successfully")

    smartHomeURL := getEnv("SMART_HOME_URL", "http://localhost:8080") // URL монолита
    shClient := services.NewSmartHomeClient(smartHomeURL)
	// --- Инициализация роутера ---
	// Передаем в роутер и БД, и паблишер
	router := handlers.SetupRouter(database, publisher, shClient)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

    port := getEnv("PORT", "8080")
	// Если переменная не установлена, используем "8080" по умолчанию
	if port == "" {
		port = "8081"
	}

	listenAddr := ":" + port

	// --- Настройка и запуск сервера ---
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

// getEnv получает переменную окружения или возвращает значение по умолчанию.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
