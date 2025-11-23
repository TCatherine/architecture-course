package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// TemperatureResponse — структура для формирования JSON-ответа

type TemperatureResponse struct {
	Temperature float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

// getTemperatureByQuery обрабатывает запросы с query-параметрами
// Пример: /temperature?location=Kitchen
func getTemperatureByQuery(c *gin.Context) {
	location := c.Query("location")
	id := ""

	generateAndRespond(c, location, id)
}

// getTemperatureByID обрабатывает запросы с ID в URL
// Пример: /temperature/2
func getTemperatureByID(c *gin.Context) {
	id := c.Param("id")
	location := ""

	generateAndRespond(c, location, id)
}

// generateAndRespond содержит общую логику для обоих обработчиков
func generateAndRespond(c *gin.Context, location string, sensorID string) {
	// Если location не указан, определяем его по sensorID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// Если sensorID не указан (может случиться только в getTemperatureByQuery), определяем его по location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	// --- ГЕНЕРАЦИЯ И ОТВЕТ ---
	rand.Seed(time.Now().UnixNano())

	minTemp := 18.0
	maxTemp := 28.0
	randomTemp := minTemp + rand.Float64()*(maxTemp-minTemp)
	randomTemp = math.Round(randomTemp*10) / 10

	response := TemperatureResponse{
		Temperature: randomTemp,
		Unit:        "°C",
		Timestamp:   time.Now().UTC(), // Используем UTC - это лучшая практика для API
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature", // Просто пример типа сенсора
		Description: fmt.Sprintf("Temperature reading from virtual sensor %s located in the %s.", sensorID, location),
	}

	c.JSON(http.StatusOK, response)
	fmt.Printf("Request processed: Location=%s, SensorID=%s, Temp=%.1f\n", location, sensorID, randomTemp)
}

func main() {
	// Получаем порт из переменной окружения "PORT"
	port := os.Getenv("PORT")
	// Если переменная не установлена, используем "8080" по умолчанию
	if port == "" {
		port = "8081"
	}

	listenAddr := ":" + port

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/temperature", getTemperatureByQuery)
	router.GET("/temperature/:id", getTemperatureByID)

	fmt.Printf("Server started on port %s\n", port)

	if err := router.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
