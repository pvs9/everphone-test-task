package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseBag struct {
	Data    interface{} `json:"data"`
	Message *string     `json:"message"`
}

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	ctx.AbortWithStatusJSON(statusCode, responseBag{nil, &message})
}

func newStatusResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.AbortWithStatusJSON(statusCode, responseBag{data, nil})
}
