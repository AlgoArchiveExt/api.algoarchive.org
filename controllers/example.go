package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://github.com/akmamun/gin-boilerplate-examples/blob/main/controllers/create_api.go

type ExampleController struct{}

func (ctrl *ExampleController) GetExampleData(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Hello World"})
}
