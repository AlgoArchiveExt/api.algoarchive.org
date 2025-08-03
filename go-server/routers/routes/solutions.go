package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SolutionsRoutes(route *gin.Engine) {
	var ctrl controllers.SolutionsController

	v1 := route.Group("/v1/solutions")

	v1.POST("/commits", ctrl.CommitProblemSolution)
	v1.GET("/:owner/:repo", ctrl.GetSolutions)
	v1.GET("/:owner/:repo/count", ctrl.GetSolutionsCount)
	v1.GET("/:owner/:repo/all-count-by-difficulty", ctrl.GetSolutionsCountByDifficulty)
}
