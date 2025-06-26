package transfer

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
// @Tags     transfer
// @Security JWT
// @Summary Get all transfers for an account
// @Description Retrieves all transfers (both incoming and outgoing) for a specific account. The account must belong to the authenticated user. You can use the 'direction' query parameter to filter the results:
//   - direction=all (default): returns both incoming and outgoing transfers for the account.
//   - direction=incoming: returns only transfers where the account is the recipient (deposits).
//   - direction=outgoing: returns only transfers where the account is the sender (withdrawals).
//
// This endpoint supports additional filtering, sorting, and pagination options.
// @param    s       query  string    false  "{'and': [ {'title': { 'cont':'cul' } } ]}"
// @param    fields  query  string    false  "fields to select eg: name,age"
// @param    page    query  int       false  "page of pagination"
// @param    limit   query  int       false  "limit of pagination"
// @param    join    query  string    false  "join relations eg: category, parent"
// @param    filter  query  []string  false  "filters eg: name||eq||ad price||gte||200"
// @param    sort    query  []string  false  "filters eg: created_at,desc title,asc"
// @param    account_id  query  string  true  "ID of the account to filter transfers by"
// @param    direction  query  string  false  "Direction of transfer: incoming, outgoing, or all (default is all)"
// @Router   /api/v1/transfer [get]
func (c *Controller) findAll(ctx *gin.Context) {
	var api crud.GetAllRequest
	if api.Limit == 0 {
		api.Limit = 20
	}
	if err := ctx.ShouldBindQuery(&api); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	accountID := ctx.Query("account_id")
	if accountID == "" {
		ctx.JSON(400, gin.H{"message": "account_id query parameter is required"})
		return
	}
	// validate the account id is belonging to the authenticated user
	authUser, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"message": "unauthorized"})
		return
	}
	authUserID, ok := authUser.(string)
	if !ok {
		ctx.JSON(401, gin.H{"message": "unauthorized"})
		return
	}
	if !c.service.isAccountBelongsToUser(ctx, accountID, authUserID) {
		ctx.JSON(403, gin.H{"message": "forbidden: account does not belong to user"})
		return
	}

	switch ctx.Query("direction") {
	case "incoming":
		api.Filter = append(api.Filter, fmt.Sprintf("to_account_id||eq||%s", accountID))
	case "outgoing":
		api.Filter = append(api.Filter, fmt.Sprintf("from_account_id||eq||%s", accountID))
	default:
		api.Filter = append(api.Filter, fmt.Sprintf("from_account_id||eq||%s|to_account_id||eq||%s", accountID, accountID))
	}

	var result []model
	var totalRows int64
	err := c.service.Find(api, &result, &totalRows)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var data interface{}
	if api.Page > 0 {
		data = map[string]interface{}{
			"data":       result,
			"total":      totalRows,
			"totalPages": int((totalRows + int64(api.Limit) - 1) / int64(api.Limit)),
		}
	} else {
		data = result
	}
	ctx.JSON(200, data)
}

// @Success  200  {object}  model
// @Tags     transfer
// @Security JWT
// @Summary Get a transfer by ID
// @Description Retrieves a single transfer by its UUID
// @param    id    path  string  true  "uuid of item"
// @Router   /api/v1/transfer/{id} [get]
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

// @Summary Execute a money transfer between accounts
// @Description Transfers money from one account to another using a transaction
// @Tags transfer
// @Security JWT
// @Accept json
// @Produce json
// @Param request body CreateTransferRequest true "Transfer payload"
// @Success 201 {object} CreateTransferResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/transfer/execute [post]
func (c *Controller) executeTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	resp, err := c.service.Transfer(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": resp})
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}
