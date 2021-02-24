package notifications

import (
	"sergiocarracedo.es/streambot-go/internal/server"
)

func SetupRoutes(server server.Server) {
	route := server.Engine.Group("notifications")
	route.GET("/subscriber", lastNotificationHandler())
}