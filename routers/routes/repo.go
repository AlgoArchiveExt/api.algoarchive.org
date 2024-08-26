package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func RepositoryRoutes(route *gin.Engine) {
	var ctrl controllers.RepositoryController
	v1 := route.Group("/v1/repo")

	v1.POST("/commit", ctrl.CommitProblemSolution)
}
