package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func ExamplesRoutes(route *gin.Engine) {
	var ctrl controllers.ExampleController
	v1 := route.Group("/v1/examples")
	v1.GET("test/", ctrl.GetExampleData)
}
