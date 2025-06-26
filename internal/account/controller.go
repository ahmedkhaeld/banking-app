package account

import (
	"fmt"
	"net/http"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/common"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

// @Success  200  {array}  model
// @Tags     account
// @param    s       query  string    false  "{'and': [ {'title': { 'cont':'cul' } } ]}"
// @param    fields  query  string    false  "fields to select eg: name,age"
// @param    page    query  int       false  "page of pagination"
// @param    limit   query  int       false  "limit of pagination"
// @param    join    query  string    false  "join relations eg: category, parent"
// @param    filter  query  []string  false  "filters eg: name||eq||ad price||gte||200"
// @param    sort    query  []string  false  "filters eg: created_at,desc title,asc"
// @Router   /api/v1/account [get]
// func (c *Controller) findAll(ctx *gin.Context) {
// 	var api crud.GetAllRequest
// 	if api.Limit == 0 {
// 		api.Limit = 20
// 	}
// 	if err := ctx.ShouldBindQuery(&api); err != nil {
// 		ctx.JSON(400, gin.H{"message": err.Error()})
// 		return
// 	}

// 	var result []model
// 	var totalRows int64
// 	err := c.service.Find(api, &result, &totalRows)
// 	if err != nil {
// 		ctx.JSON(400, gin.H{"message": err.Error()})
// 		return
// 	}

// 	var data interface{}
// 	if api.Page > 0 {
// 		data = map[string]interface{}{
// 			"data":       result,
// 			"total":      totalRows,
// 			"totalPages": int(math.Ceil(float64(totalRows) / float64(api.Limit))),
// 		}
// 	} else {
// 		data = result
// 	}
// 	ctx.JSON(200, data)
// }

// @Success  200  {object}  model
// @Tags     account
// @Security JWT
// @param    id    path  string  true  "uuid of item"
// @Router   /api/v1/account/{id} [get]
func (c *Controller) findOne(ctx *gin.Context) {
	var api crud.GetAllRequest
	var item common.ById
	if err := ctx.ShouldBindQuery(&api); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	api.Filter = append(api.Filter, fmt.Sprintf("id||eq||%s", item.ID))

	var result model

	err := c.service.FindOne(api, &result)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, result)
}

// @Success  201  {object}  CreateAccountResponse
// @Tags     account
// @Security JWT
// @param    request  body  CreateAccountRequest  true  "Account creation payload"
// @Router   /api/v1/account [post]
func (c *Controller) create(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "user_id not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "user_id in context is not a string"})
		return
	}

	resp, err := c.service.CreateAccount(req, userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": resp})
}

// @Success  200  {string}  string  "ok"
// @Tags     account
// @param    id  path  string  true  "uuid of item"
// @Router   /api/v1/account/{id} [delete]
// func (c *Controller) delete(ctx *gin.Context) {
// 	var item common.ById
// 	if err := ctx.ShouldBindUri(&item); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}

// 	id, err := uuid.ParseBytes([]byte(item.ID))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	err = c.service.Delete(&model{ID: id})
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
// }

// @Success  200  {string}  string  "ok"
// @Tags     account
// @param    id  path  string  true  "uuid of item"
// @param    item  body  model   true  "update body"
// @Router   /api/v1/account/{id} [patch]
// func (c *Controller) update(ctx *gin.Context) {
// 	var item model
// 	var byId common.ById
// 	if err := ctx.ShouldBind(&item); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	if err := ctx.ShouldBindUri(&byId); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	id, err := uuid.Parse(byId.ID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	if err := ctx.ShouldBindUri(&byId); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	err = c.service.Update(&model{ID: id}, &item)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, item)
// }

// GetAccountBalance godoc
// @Summary  Get account balance
// @Description  Returns the balance and currency for a specific account
// @Tags     account
// @Security JWT
// @Param    id  path  string  true  "Account ID"
// @Success  200  {object}  AccountBalanceResponse
// @Failure  404  {object}  map[string]string
// @Router   /api/v1/account/{id}/balance [get]
func (c *Controller) getAccountBalance(ctx *gin.Context) {
	accountID := ctx.Param("id")
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "user_id not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "user_id in context is not a string"})
		return
	}
	resp, err := c.service.GetAccountBalance(accountID, userIDStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "account not found or not owned by user"})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}
