package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pvs9/everphone-test-task/pkg/handler/middleware"
	"github.com/pvs9/everphone-test-task/pkg/service"
)

type Employee interface {
	getAll(ctx *gin.Context)
	getById(ctx *gin.Context)
	updateById(ctx *gin.Context)
	uploadDataset(ctx *gin.Context)
}

type Gift interface {
	getAll(ctx *gin.Context)
	getById(ctx *gin.Context)
	updateById(ctx *gin.Context)
	uploadDataset(ctx *gin.Context)
}

type Handler struct {
	Employee
	Gift
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Employee: NewEmployeeHandler(services.Employee),
		Gift:     NewGiftHandler(services.Gift),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.BasicAuth())

	api := router.Group("/api")
	{
		employee := api.Group("/employee")
		{
			employee.GET("/", h.Employee.getAll)
			employee.GET("/:id", h.Employee.getById)
			employee.PUT("/:id", h.Employee.updateById)
			employee.POST("/dataset", h.Employee.uploadDataset)
		}

		gift := api.Group("/gift")
		{
			gift.GET("/", h.Gift.getAll)
			gift.GET("/:id", h.Gift.getById)
			gift.PUT("/:id", h.Gift.updateById)
			gift.POST("/dataset", h.Gift.uploadDataset)
		}
	}

	return router
}
