package transfer

import (
	"github.com/ahmedkhaeld/banking-app/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routerGroup *gin.RouterGroup) {
	service := InitService()
	controller := NewController(service)

	routerGroup.GET("", auth.UserMiddleware(), controller.findAll)
	routerGroup.GET(":id", auth.UserMiddleware(), controller.findOne)
	routerGroup.POST("", auth.UserMiddleware(), controller.executeTransfer)

}
