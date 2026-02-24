package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	svc TodoService
}

func NewHandler(svc TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

type createTodoInput struct {
	Title string `json:"title"`
}

type todoResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed string `json:"completed"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	var in createTodoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	uid := c.GetString("userID")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	t, err := h.svc.CreateTodoService(c.Request.Context(), uid, in.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out := todoResponse{ID: t.ID.String(), Title: t.Title, Completed: t.Completed}
	c.JSON(http.StatusCreated, out)
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	uid := c.GetString("userID")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	todos, err := h.svc.GetTodosService(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]todoResponse, 0, len(todos))
	for _, t := range todos {
		out = append(out, todoResponse{ID: t.ID.String(), Title: t.Title, Completed: t.Completed})
	}
	c.JSON(http.StatusOK, out)
}
