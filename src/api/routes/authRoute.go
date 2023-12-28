package routes

import (
	"context"
	"net/http"

	pb "gateway/proto/auth"
	"gateway/src/api/middleware"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Router struct holds the client for the MicroService to make GRPC calls.
type AuthRouter struct {
	client pb.AuthServiceClient // GRPC client for the Micro Service
}

// NewRouter creates and returns a new Router instance with a connected MicroService client.
func NewAuthRouter(conn *grpc.ClientConn) *AuthRouter {
	return &AuthRouter{
		client: pb.NewAuthServiceClient(conn), // Initialize the MicroService client
	}
}

// RegisterRoutes sets up the API endpoints (routes) for the micro service.
func (gr *AuthRouter) RegisterRoutes(router *gin.Engine) {
	grp := router.Group("/auth") // Group routes under '/' prefix

    grp.POST("/register", gr.registerHandler)
    grp.POST("/login", gr.loginHandler)
    grp.GET("/secretkey", middleware.Authorize(), gr.getSecretKey)
}

func (ar *AuthRouter) registerHandler(c *gin.Context) {
    var request pb.RegisterRequest
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
    response, err := ar.client.Register(context.Background(), &request)
    if err != nil {
        handleGRPCError(c, err)
        return
    }

    // On successful GRPC call, return the service response with a 200 OK status.
    c.JSON(http.StatusOK, response)
}

func (ar *AuthRouter) loginHandler(c *gin.Context) {
    var request pb.LoginRequest
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
    response, err := ar.client.Login(context.Background(), &request)
    if err != nil {
        handleGRPCError(c, err)
        return
    }

    // On successful GRPC call, return the service response with a 200 OK status.
    c.JSON(http.StatusOK, response)
}

func (ar *AuthRouter) getSecretKey(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
        // If the userId doesn't exist in the context for some reason, handle the error
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "StatusUnauthorized",
			"message": "User not authorized or user ID missing",
		}})
        return
    }
	
    var request pb.SecretKeyRequest
	request.UserId = userId.(string)

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
    response, err := ar.client.GetSecretKey(context.Background(), &request)
    if err != nil {
        handleGRPCError(c, err)
        return
    }

    // On successful GRPC call, return the service response with a 200 OK status.
    c.JSON(http.StatusOK, response)
}

