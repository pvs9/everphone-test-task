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

type GiftHandler struct {
	service service.Gift
}

func NewGiftHandler(service service.Gift) *GiftHandler {
	return &GiftHandler{service: service}
}

func (h *GiftHandler) getAll(ctx *gin.Context) {
	gifts, err := h.service.GetAll()

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if gifts == nil {
		gifts = []christmas.Gift{}
	}

	newStatusResponse(ctx, http.StatusOK, gifts)
}

func (h *GiftHandler) getById(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.ParseInt(stringId, 10, 64)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entity found")
		return
	}

	gift, err := h.service.GetById(id)

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	newStatusResponse(ctx, http.StatusOK, gift)
}

func (h *GiftHandler) updateById(ctx *gin.Context) {
	stringId := ctx.Param("id")

	var giftData request.GiftUpdateRequest

	id, err := strconv.ParseInt(stringId, 10, 64)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entity found")
		return
	}

	gift, err := h.service.GetById(id)

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if gift == nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entity found")
	}

	if err := ctx.ShouldBind(&giftData); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	updatedGift, err := h.service.Update(*gift, giftData)

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	newStatusResponse(ctx, http.StatusOK, updatedGift)
}

func (h *GiftHandler) uploadDataset(ctx *gin.Context) {
	var datasetData request.GiftDatasetRequest

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
