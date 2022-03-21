package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/pvs9/everphone-test-task/pkg/request"
	"github.com/pvs9/everphone-test-task/pkg/service"
	"net/http"
	"path/filepath"
	"strconv"
)

type EmployeeHandler struct {
	service service.Employee
}

func NewEmployeeHandler(service service.Employee) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) getAll(ctx *gin.Context) {
	employees, err := h.service.GetAll()

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if employees == nil {
		employees = []christmas.Employee{}
	}

	newStatusResponse(ctx, http.StatusOK, employees)
}

func (h *EmployeeHandler) getById(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.ParseInt(stringId, 10, 64)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entity found")
		return
	}

	employee, err := h.service.GetById(id)

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	newStatusResponse(ctx, http.StatusOK, employee)
}

func (h *EmployeeHandler) updateById(ctx *gin.Context) {
	stringId := ctx.Param("id")

	var employeeData request.EmployeeUpdateRequest

	id, err := strconv.ParseInt(stringId, 10, 64)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entity found")
		return
	}

	employee, err := h.service.GetById(id)

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if employee == nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entity found")
	}

	if err := ctx.ShouldBind(&employeeData); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	updatedEmployee, err := h.service.Update(*employee, employeeData)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	newStatusResponse(ctx, http.StatusOK, updatedEmployee)
}

func (h *EmployeeHandler) uploadDataset(ctx *gin.Context) {
	var datasetData request.EmployeeDatasetRequest

	if err := ctx.ShouldBind(&datasetData); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	extension := filepath.Ext(datasetData.Dataset.Filename)
	newFileName := uuid.New().String() + extension

	if err := ctx.SaveUploadedFile(datasetData.Dataset, "/go/src/everphone-test-task.io/pkg/storage/"+newFileName); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	messageId, err := h.service.UploadDataset(newFileName)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	newStatusResponse(ctx, http.StatusOK, messageId)
}
