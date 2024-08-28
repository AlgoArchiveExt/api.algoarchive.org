package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SolutionsRoutes(route *gin.RouterGroup) {
	var ctrl controllers.SolutionsController

	v1 := route.Group("/v1/solutions")

	v1.POST("/commits", ctrl.CommitProblemSolution)
	v1.GET("/:owner/:repo", ctrl.GetSolutions)
}
