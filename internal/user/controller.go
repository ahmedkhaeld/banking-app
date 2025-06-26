package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/common"
	"github.com/ahmedkhaeld/banking-app/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	//JWT Expire
	JWTExpire = 24 * 60 * 60 * time.Second // 24 hours
)

type Controller struct {
	service *Service
}

// @Success  200  {object}  model
// @Tags     user
// @param    id    path  string  true  "uuid of item"
// @Router   /api/v1/user/{id} [get]
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

// @Success  201  {object}  model
// @Tags     user
// @param    {object}  body  CreateUserRequest  true  "item to create"
// @Router   /api/v1/user [post]
func (c *Controller) create(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := c.service.createUser(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

// @Success  200  {string}  string  "ok"
// @Tags     user
// @Security JWT
// @Description Update user information
// @Description Update user information by ID
// @param    id  path  string  true  "uuid of item"
// @param    item  body  model   true  "update body"
// @Router   /api/v1/user/{id} [patch]
func (c *Controller) update(ctx *gin.Context) {
	var item model
	var byId common.ById
	if err := ctx.ShouldBind(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&byId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id, err := uuid.Parse(byId.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&byId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = c.service.Update(&model{ID: id}, &item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param request body LoginUserRequest true "Login credentials"
// @Success 200 {object} LoginUserResponse "JWT token and user info"
// @Failure 400 {object} map[string]string
// @Router /api/v1/user/login [post]
func (c *Controller) login(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	user, err := c.service.loginUser(req.Username, req.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	// Create JWT token
	jwtMaker, err := auth.NewJWTMaker()
	if err != nil {
		ctx.JSON(500, gin.H{"message": "could not create token maker"})
		return
	}
	userIDStr := user.ID.String()
	if userIDStr == "" {
		ctx.JSON(400, gin.H{"message": "user ID is empty"})
		return
	}
	token, _, err := jwtMaker.CreateToken(userIDStr, JWTExpire)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "could not create token"})
		return
	}
	resp := LoginUserResponse{
		AccessToken: token,
	}
	resp.User.ID = user.ID.String()
	resp.User.Username = user.Username
	resp.User.FullName = user.FullName
	resp.User.Email = user.Email
	ctx.JSON(200, resp)
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}
