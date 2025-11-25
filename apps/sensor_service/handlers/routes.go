// file: handlers/house.go
package handlers

import (
	"smart-home-service/db"
	"smart-home-service/message_broker"
	"smart-home-service/services"

	"github.com/gin-gonic/gin"
)


func SetupRouter(db *db.DB, publisher *message_broker.Publisher, shClient *services.SmartHomeClient) *gin.Engine {
	r := gin.Default()

	// Создаем экземпляр обработчиков
	homeHandler := NewHomeHandler(db, publisher)
	sensorHandler := NewSensorHandler(db, shClient)

	// Группируем роуты для API v1
	apiV1 := r.Group("/api/v1")
	{
		// Группа роутов для домов
		homes := apiV1.Group("/homes")
		{
			homes.GET("", homeHandler.GetHomesHandler)
		}
	}
    {
		// Группа роутов для домов
		home := apiV1.Group("/home")
		{
			home.POST("", homeHandler.CreateHomeHandler)
			home.GET("/:id", homeHandler.GetHomeByIDHandler)
			home.PUT("/:id", homeHandler.UpdateHomeHandler)
			home.DELETE("/:id", homeHandler.DeleteHomeHandler)
			home.POST("/:id/sensor", sensorHandler.CreateSensorProxyHandler)
			home.GET("/:id/sensors", sensorHandler.GetSensorsHandler)
		}
	}

	return r
}
