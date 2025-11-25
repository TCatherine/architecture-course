package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"smart-home-service/models"
	"strings"
	"time"
)

type SmartHomeClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewSmartHomeClient(url string) *SmartHomeClient {
	return &SmartHomeClient{
		BaseURL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterDevice отправляет запрос в монолит и возвращает ID созданного устройства
func (c *SmartHomeClient) RegisterDevice(payload models.SensorCreatePayload) (int, error) {
	// 1. Маппинг данных.
	// Монолит ждет lowercase type (например, "temperature"), а мы получаем "TEMPERATURE_SENSOR"
	monolithType := strings.ToLower(payload.Type)
	if strings.Contains(monolithType, "temperature") {
		monolithType = "temperature"
	}

	reqBody := models.MonolithSensorCreate{
		Name:     payload.Name,
		Type:     monolithType,
		Location: payload.Location,
		Unit: "°C",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	// 2. Вызов API Монолита
	url := fmt.Sprintf("%s/api/v1/sensors", c.BaseURL)
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to call smart_home: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("smart_home returned status: %d", resp.StatusCode)
	}

	// 3. Парсинг ответа для получения ID
	var result models.MonolithSensorResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ID, nil
}

// GetSensorByID запрашивает данные о датчике из монолита
func (c *SmartHomeClient) GetSensorByID(serviceID int) (*models.SensorDetail, error) {
	url := fmt.Sprintf("%s/api/v1/sensors/%d", c.BaseURL, serviceID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sensor %d: %w", serviceID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("smart_home returned status %d for sensor %d", resp.StatusCode, serviceID)
	}

	var sensor models.SensorDetail
	if err := json.NewDecoder(resp.Body).Decode(&sensor); err != nil {
		return nil, fmt.Errorf("failed to decode sensor data: %w", err)
	}

	return &sensor, nil
}
