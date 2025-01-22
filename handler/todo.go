package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wei840222/go-restful-sample/storage"
)

type TodoHandler struct {
	storage storage.TodoStorage
}

type CreateTodoReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	var req CreateTodoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo storage.Todo
	todo.Title = req.Title
	todo.Description = req.Description

	if err := h.storage.Create(c, &todo); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &GetTodoRes{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	})
}

type ListTodoRes []*GetTodoRes

func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.storage.List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res ListTodoRes
	res = make(ListTodoRes, 0, len(todos))
	for _, todo := range todos {
		res = append(res, &GetTodoRes{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

type GetTodoRes struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *TodoHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.storage.Get(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &GetTodoRes{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	})
}

type UpdateTodoReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateTodoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo storage.Todo
	todo.ID = uint(id)
	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed

	if err := h.storage.Update(c, uint(id), &todo); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &GetTodoRes{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.Delete(c, uint(id)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func RegisterTodoHandler(e *gin.Engine, s storage.TodoStorage) error {
	h := &TodoHandler{
		storage: s,
	}

	todo := e.Group("/todos")
	{
		todo.GET("/", h.List)
		todo.POST("/", h.Create)
		todo.GET("/:id", h.Get)
		todo.PATCH("/:id", h.Update)
		todo.DELETE("/:id", h.Delete)
	}

	return nil
}
