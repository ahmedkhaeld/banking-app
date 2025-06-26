package account

import (
	"github.com/ahmedkhaeld/banking-app/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routerGroup *gin.RouterGroup) {
	service := InitService()
	controller := NewController(service)

	//Admin route
	// routerGroup.GET("", controller.findAll)

	routerGroup.GET(":id", auth.UserMiddleware(), controller.findOne)
	routerGroup.POST("", auth.UserMiddleware(), controller.create)
	routerGroup.GET(":id/balance", auth.UserMiddleware(), controller.getAccountBalance)
	// routerGroup.DELETE(":id", auth.BearerMiddleware(), controller.delete)
	routerGroup.PATCH(":id/balance", auth.UserMiddleware(), controller.updateBalance)
}
