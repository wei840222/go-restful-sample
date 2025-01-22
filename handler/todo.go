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

type GetTodoRes struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *TodoHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		return
	}

	todo, err := h.storage.Get(c, id)
	if err != nil {
		if storage.IsNotFound(err) {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorRes{Error: err.Error()})
		} else {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRes{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, &GetTodoRes{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   *todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	})
}

type ListTodoRes []GetTodoRes

func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.storage.List(c)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRes{Error: err.Error()})
		return
	}

	var res ListTodoRes
	res = make(ListTodoRes, 0, len(todos))
	for _, todo := range todos {
		getTodoRes := GetTodoRes{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
		if todo.Completed != nil {
			getTodoRes.Completed = *todo.Completed
		}
		res = append(res, getTodoRes)
	}

	c.JSON(http.StatusOK, res)
}

type CreateTodoReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	var req CreateTodoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		return
	}

	var todo storage.Todo
	todo.Title = req.Title
	todo.Description = req.Description

	if err := h.storage.Create(c, &todo); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRes{Error: err.Error()})
		return
	}

	res := GetTodoRes{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
	if todo.Completed != nil {
		res.Completed = *todo.Completed
	}
	c.JSON(http.StatusCreated, res)
}

type UpdateTodoReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		return
	}

	var req UpdateTodoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		return
	}

	var todo storage.Todo
	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed

	if err := h.storage.Update(c, id, todo); err != nil {
		if storage.IsNotFound(err) {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorRes{Error: err.Error()})
		} else {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRes{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		return
	}

	if err := h.storage.Delete(c, id); err != nil {
		if storage.IsNotFound(err) {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorRes{Error: err.Error()})
		} else {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRes{Error: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
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
