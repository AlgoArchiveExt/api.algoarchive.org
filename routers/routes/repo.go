package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func RepositoryRoutes(route *gin.RouterGroup) {
	var ctrl controllers.RepositoryController
	v1 := route.Group("/v1/repository")

	v1.POST("/commit", ctrl.CommitProblemSolution)
}
