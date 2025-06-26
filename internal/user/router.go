package user

import (
	"github.com/ahmedkhaeld/banking-app/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routerGroup *gin.RouterGroup) {
	service := InitService()
	controller := NewController(service)

	routerGroup.GET(":id", controller.findOne)
	routerGroup.POST("", controller.create)
	routerGroup.PATCH(":id", auth.UserMiddleware(), controller.update)
	routerGroup.POST("/login", controller.login)
}
