package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

// handleGRPCError handles errors from GRPC and sends appropriate HTTP response.
func handleGRPCError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "InternalServerError",
			"message": "An unexpected error occurred",
		}})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
		"code":    st.Code().String(),
		"message": st.Message(),
	}})
}