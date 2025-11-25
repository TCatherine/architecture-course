package models

import (
	"time"
)

// Home представляет умный дом, принадлежащий пользователю
type Home struct {
	HomeID    int       `json:"home_id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Street    string    `json:"street"`
	Num       int       `json:"num"`
	CreatedAt time.Time `json:"created_at"`
}

// HomeCreate используется для создания нового дома.
type HomeCreate struct {
	UserID int       `json:"user_id"`
	Name   string    `json:"name" binding:"required"`
	City   string    `json:"city"`
	Street string    `json:"street"`
	Num    int       `json:"num"`
}

// HomeUpdate используется для частичного обновления дома.
type HomeUpdate struct {
    UserID int    `json:"user_id"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Street string `json:"street"`
	Num    int    `json:"num"`
}

