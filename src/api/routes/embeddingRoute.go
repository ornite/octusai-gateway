package routes

import (
	"context"
	"log"
	"net/http"

	pb "gateway/proto/embedding"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Router struct holds the client for the MicroService to make GRPC calls.
type EmbeddingRouter struct {
	client pb.EmbeddingServiceClient // GRPC client for the Micro Service
}

// NewRouter creates and returns a new Router instance with a connected MicroService client.
func NewEmbeddingRouter(conn *grpc.ClientConn) *EmbeddingRouter {
	return &EmbeddingRouter{
		client: pb.NewEmbeddingServiceClient(conn), // Initialize the MicroService client
	}
}

// RegisterRoutes sets up the API endpoints (routes) for the micro service.
func (gr *EmbeddingRouter) RegisterRoutes(router *gin.Engine) {
	grp := router.Group("/embedding") // Group routes under '/' prefix

	grp.GET("/text", gr.textEmbeddingHandler) // Route
}

func (ar *EmbeddingRouter) textEmbeddingHandler(c *gin.Context) {
	log.Println("Get")
    var request pb.TextRequest
	// Attempt to bind the incoming JSON payload to the Request struct.
	if err := c.ShouldBindJSON(&request); err != nil {
		// If binding fails, return a 400 Bad Request with the error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "BadRequest",
			"message": err.Error(),
		}})
		return
	}

    // Call the GRPC service to get response, passing the structured request.
    response, err := ar.client.GetTextEmbedding(context.Background(), &request)
    if err != nil {
        handleGRPCError(c, err)
        return
    }

    // On successful GRPC call, return the service response with a 200 OK status.
    c.JSON(http.StatusOK, response)
}