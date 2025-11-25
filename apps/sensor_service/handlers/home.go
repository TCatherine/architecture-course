// file: handlers/house.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"smart-home-service/db"
	"smart-home-service/message_broker"
	"smart-home-service/models"

	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HomesExchange = "homes_exchange"
)

// HomeHandler инкапсулирует зависимости для обработчиков домов.
type HomeHandler struct {
	DB        *db.DB
	Publisher *message_broker.Publisher
}

// NewHomeHandler создает новый экземпляр HomeHandler.
func NewHomeHandler(db *db.DB, publisher *message_broker.Publisher) *HomeHandler {
	return &HomeHandler{
		DB:        db,
		Publisher: publisher,
	}
}


// GetHomesHandler обрабатывает GET-запрос для получения всех домов.
func (h *HomeHandler) GetHomesHandler(c *gin.Context) {
	homes, err := h.DB.GetHomes(c.Request.Context())
	if err != nil {
	    log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve homes"})
		return
	}
	c.JSON(http.StatusOK, homes)
}

// GetHomeByIDHandler обрабатывает GET-запрос для получения дома по ID.
func (h *HomeHandler) GetHomeByIDHandler(c *gin.Context) {
	homeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
	    log.Printf("ERROR: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home ID"})
		return
	}

	home, err := h.DB.GetHomeByID(c.Request.Context(), homeID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve home"})
		}
		return
	}

	c.JSON(http.StatusOK, home)
}

// CreateHomeHandler обрабатывает POST-запрос для создания нового дома.
func (h *HomeHandler) CreateHomeHandler(c *gin.Context) {
	var homeCreate models.HomeCreate
	if err := c.ShouldBindJSON(&homeCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newHome, err := h.DB.CreateHome(c.Request.Context(), homeCreate)
	if err != nil {
	    log.Printf("WARN: Failed to creat home: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create home"})
		return
	}

	// Публикуем событие о создании дома
	body, _ := json.Marshal(newHome)
	if err := h.Publisher.Publish(HomesExchange, "home.created", body); err != nil {
		log.Printf("WARN: Failed to publish home.created event: %v", err)
		// Не возвращаем ошибку клиенту, т.к. основная операция (создание) прошла успешно
	}

	c.JSON(http.StatusCreated, newHome)
}

// UpdateHomeHandler обрабатывает PUT-запрос для обновления дома.
func (h *HomeHandler) UpdateHomeHandler(c *gin.Context) {
	homeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home ID"})
		return
	}

	var homeUpdate models.HomeUpdate
	if err := c.ShouldBindJSON(&homeUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedHome, err := h.DB.UpdateHome(c.Request.Context(), homeID, homeUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update home"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedHome)
}

// DeleteHomeHandler обрабатывает DELETE-запрос для удаления дома.
func (h *HomeHandler) DeleteHomeHandler(c *gin.Context) {
	homeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home ID"})
		return
	}

	err = h.DB.DeleteHome(c.Request.Context(), homeID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete home"})
		}
		return
	}

	// Публикуем событие об удалении
	body, _ := json.Marshal(gin.H{"home_id": homeID})
	if err := h.Publisher.Publish(HomesExchange, "home.deleted", body); err != nil {
		log.Printf("WARN: Failed to publish home.deleted event: %v", err)
	}

	c.Status(http.StatusNoContent)
}