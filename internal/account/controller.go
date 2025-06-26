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

	resp, err := c.service.createAccount(req, userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": resp})
}

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
	resp, err := c.service.getAccountBalance(accountID, userIDStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "account not found or not owned by user"})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// UpdateBalance godoc
// @Summary Update account balance
// @Description Updates the balance of a specific account by a given amount
// @Tags account
// @Security JWT
// @Param id path string true "Account ID"
// @Param request body UpdateAccountBalanceRequest true "Balance update payload"
// @Success 200 {object} model
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/account/{id}/balance [patch]
func (c *Controller) updateBalance(ctx *gin.Context) {
	accountID := ctx.Param("id")
	var req UpdateAccountBalanceRequest
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
	// Optionally, check ownership here if needed
	if !c.service.isAccountOwnedByUser(accountID, userIDStr) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "account does not belong to user"})
		return
	}
	// Update the account balance
	if req.Amount == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "amount must be non-zero"})
		return
	}
	if req.Amount < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "amount must be positive"})
		return
	}
	account, err := c.service.updateBalance(ctx, accountID, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": account})
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}
