package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/localopsco/go-sample/service"
)

type Handler struct {
	svc *service.TaskService
}

func NewHandler(taskSvc *service.TaskService) *Handler {
	return &Handler{
		taskSvc,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsCompleted bool   `json:"is_completed"`
	}
	err := c.ShouldBindWith(&reqBody, binding.JSON)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
		})
		return
	}
	task, err := h.svc.CreateTask(reqBody.Title, reqBody.Description, reqBody.IsCompleted)
	if err != nil {
		log.Printf("error creating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid data",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) GetTask(c *gin.Context) {
	taskIDStr := strings.TrimSpace(c.Param("task_id"))
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid task_id",
		})
		return
	}
	task, err := h.svc.GetTask(taskID)
	if err != nil {
		if err.Error() == service.TaskNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Printf("error getting task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error. Something went wrong.",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) UpdateTask(c *gin.Context) {

	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsCompleted bool   `json:"is_completed"`
	}
	err := c.ShouldBindWith(&reqBody, binding.JSON)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
		})
		return
	}

	taskIDStr := strings.TrimSpace(c.Param("task_id"))
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid task_id",
		})
		return
	}

	task, err := h.svc.UpdateTask(taskID, reqBody.Title, reqBody.Description, reqBody.IsCompleted)
	if err != nil {
		if err.Error() == service.TaskNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Printf("error creating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid data",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	taskIDStr := strings.TrimSpace(c.Param("task_id"))
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid task_id",
		})
		return
	}
	err = h.svc.DeleteTask(taskID)
	if err != nil {
		if err.Error() == service.TaskNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Printf("error getting task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error. Something went wrong.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *Handler) ListTasks(c *gin.Context) {
	tasks, err := h.svc.ListTasks()
	if err != nil {
		log.Printf("error getting task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error. Something went wrong.",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) AddAttachment(c *gin.Context) {
	taskIDStr := strings.TrimSpace(c.Param("task_id"))
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid task_id",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid attachment",
		})
		return
	}

	task, err := h.svc.AddAttachment(taskID, file)
	if err != nil {
		if err.Error() == service.TaskNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Printf("error getting task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error. Something went wrong.",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}
