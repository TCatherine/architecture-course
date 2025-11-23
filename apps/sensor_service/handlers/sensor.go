package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"smart-home-service/db"
	"smart-home-service/models"
	"smart-home-service/services"

	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	DB              *db.DB
	SmartHomeClient *services.SmartHomeClient
}

func NewSensorHandler(db *db.DB, client *services.SmartHomeClient) *SensorHandler {
	return &SensorHandler{
		DB:              db,
		SmartHomeClient: client,
	}
}

// CreateSensorProxyHandler обрабатывает создание датчика через проксирование в монолит
func (h *SensorHandler) CreateSensorProxyHandler(c *gin.Context) {
	// 1. Получаем Home ID из URL
	homeID, err := strconv.Atoi(c.Param("id")) // В роутере будет :id (home_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home ID"})
		return
	}

	// 2. Валидируем входящий JSON
	var payload models.SensorCreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Отправляем запрос в Smart Home Monolith
	serviceID, err := h.SmartHomeClient.RegisterDevice(payload)
	if err != nil {
		log.Printf("ERROR: Failed to register device in Smart Home: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to register device in upstream service"})
		return
	}

	// 4. Сохраняем связь в локальной БД
	link := models.Sensor{
		HomeID:    homeID,
		ServiceID: serviceID,
	}

	if err := h.DB.CreateSensorLink(c.Request.Context(), link); err != nil {
		// Примечание: Устройство в монолите уже создано.
		// В продакшене тут нужна бы компенсация (удаление из монолита) или очередь повторных попыток.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Device created remotely but failed to link locally"})
		return
	}

	// 5. Возвращаем результат
	// Возвращаем структуру связи, либо можно вернуть original payload + ID
	c.JSON(http.StatusCreated, gin.H{
		"home_id":    link.HomeID,
		"service_id": link.ServiceID,
		"status":     "linked",
	})
}


// GetSensorsHandler получает список датчиков для дома с данными из монолита
func (h *SensorHandler) GetSensorsHandler(c *gin.Context) {
	// 1. Получаем Home ID
	homeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home ID"})
		return
	}

	// 2. Получаем список ID из локальной БД
	sensorIDs, err := h.DB.GetSensorIDsByHomeID(c.Request.Context(), homeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sensor links"})
		return
	}

	// Если датчиков нет, возвращаем пустой список
	if len(sensorIDs) == 0 {
		c.JSON(http.StatusOK, []models.SensorDetail{})
		return
	}

	// 3. Запрашиваем данные из монолита для каждого датчика
	// Можно делать последовательно или параллельно. Для простоты сделаем параллельно через WaitGroup
	var result []models.SensorDetail
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range sensorIDs {
		wg.Add(1)
		go func(serviceID int) {
			defer wg.Done()
			detail, err := h.SmartHomeClient.GetSensorByID(serviceID)
			if err != nil {
				log.Printf("WARN: Failed to fetch sensor %d details: %v", serviceID, err)
				// В случае ошибки просто пропускаем этот датчик или добавляем с пометкой "error"
				return
			}

			mu.Lock()
			result = append(result, *detail)
			mu.Unlock()
		}(id)
	}

	wg.Wait()

	c.JSON(http.StatusOK, result)
}
