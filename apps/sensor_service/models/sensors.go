package models

import (
	"time"
)

// Sensor (пока как заглушка для будущей работы)
type Sensor struct {
	ServiceID int `json:"service_id"`
	HomeID    int `json:"home_id"`
}

type SensorDetail struct {
	ID          int       `json:"id"`          // ID в монолите (service_id)
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

// SensorCreatePayload - то, что присылает клиент
type SensorCreatePayload struct {
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required"` // e.g. "TEMPERATURE_SENSOR"
	Location     string `json:"location" binding:"required"`
	Address      string `json:"address"`       // Игнорируем при сохранении (нет полей в БД)
	SerialNumber int64  `json:"serial_number"` // Игнорируем при сохранении (нет полей в БД)
	State        string `json:"state"`
}

// MonolithSensorCreate - то, что мы отправляем в монолит smart_home
type MonolithSensorCreate struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Unit     string `json:"unit,omitempty"` // Можно добавить дефолтные значения
}

// MonolithSensorResponse - ответ от монолита
type MonolithSensorResponse struct {
	ID int `json:"id"`
}