package responses

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Convenience shorthand for returning a bad request response and message:
/*
Original Use:
c.JSON(http.StatusBadRequest, gin.H{
	"error": fmt.Sprintf("Info message: %s", errorMessage),
})
*/
func GiveErrorResponse(c *gin.Context, infoMessage string, errorMessage string, extraFields *map[string]any) {
	response := gin.H{
		"error": fmt.Sprintf("%s: %s", infoMessage, errorMessage),
	}

	if extraFields != nil {
		for key, value := range *extraFields {
			response[key] = value
		}
	}

	c.JSON(http.StatusBadRequest, response)
}

// Convenience shorthand for returning an OK response and message:
/*
Original Use:
c.JSON(http.StatusOK, gin.H{
	"message": "Message",
})
*/
func GiveOKResponse(c *gin.Context, message string, extraFields *map[string]any) {
	response := gin.H{
		"message": message,
	}

	if extraFields != nil {
		for key, value := range *extraFields {
			response[key] = value
		}
	}

	c.JSON(http.StatusOK, response)
}
