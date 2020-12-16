package handler

import (
	"github.com/wei840222/go-restful-sample/internal/store"

	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUserHandler register /api/users api handlers
func RegisterUserHandler(engine *gin.Engine, userStore store.UserStore) {
	h := &userHandler{userStore}

	user := engine.Group("/api/users")
	{
		user.POST("/", h.create)
	}
}

type userHandler struct {
	userStore store.UserStore
}

type userRes struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type userCreateReq struct {
	Name string `json:"name" binding:"required"`
}

// @Summary createUser
// @tags user
// @Description create user by accessToken
// @Accept  application/json
// @Produce application/json
// @Param Body body userCreateReq true "user data"
// @Success 201 "success"
// @Failure 400 "invalid request body format"
// @Router /api/users [post]
func (h userHandler) create(c *gin.Context) {
	var req userCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newUser, err := h.userStore.Create(c.Request.Context(), &store.User{
		Name: req.Name,
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, userRes{
		ID:   newUser.ID,
		Name: newUser.Name,
	})
}
